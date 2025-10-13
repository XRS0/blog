package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Visibility string

const (
	VisibilityPublic  Visibility = "public"
	VisibilityPrivate Visibility = "private"
	VisibilityLink    Visibility = "link"
)

type Article struct {
	bun.BaseModel `bun:"table:articles,alias:a"`

	ID          uint64     `bun:"id,pk,autoincrement"`
	UserID      uint64     `bun:"user_id,notnull"`
	Title       string     `bun:"title,notnull"`
	Content     string     `bun:"content,notnull"`
	Visibility  Visibility `bun:"visibility,notnull,default:'public'"`
	AccessToken string     `bun:"access_token"`
	CreatedAt   time.Time  `bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt   time.Time  `bun:"updated_at,nullzero,notnull,default:current_timestamp"`
}

type ArticleRepository struct {
	db *bun.DB
}

func NewArticleRepository(db *bun.DB) *ArticleRepository {
	return &ArticleRepository{db: db}
}

func (r *ArticleRepository) Create(userID uint64, title, content string, visibility Visibility) (*Article, error) {
	ctx := context.Background()

	article := &Article{
		UserID:     userID,
		Title:      title,
		Content:    content,
		Visibility: visibility,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if visibility == VisibilityLink {
		article.AccessToken = uuid.New().String()
	}

	_, err := r.db.NewInsert().Model(article).Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create article: %w", err)
	}

	return article, nil
}

func (r *ArticleRepository) GetByID(id uint64) (*Article, error) {
	ctx := context.Background()
	article := new(Article)

	err := r.db.NewSelect().
		Model(article).
		Where("id = ?", id).
		Scan(ctx)

	if err != nil {
		return nil, fmt.Errorf("article not found: %w", err)
	}

	return article, nil
}

func (r *ArticleRepository) Update(id, userID uint64, title, content string, visibility Visibility) (*Article, error) {
	ctx := context.Background()

	// Get existing article to check access token
	existing, err := r.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Check ownership
	if existing.UserID != userID {
		return nil, fmt.Errorf("unauthorized")
	}

	// Generate new access token if changing to link visibility
	accessToken := existing.AccessToken
	if visibility == VisibilityLink && existing.Visibility != VisibilityLink {
		accessToken = uuid.New().String()
	} else if visibility != VisibilityLink {
		accessToken = ""
	}

	existing.Title = title
	existing.Content = content
	existing.Visibility = visibility
	existing.AccessToken = accessToken
	existing.UpdatedAt = time.Now()

	_, err = r.db.NewUpdate().
		Model(existing).
		WherePK().
		Exec(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to update article: %w", err)
	}

	return existing, nil
}

func (r *ArticleRepository) Delete(id, userID uint64) error {
	ctx := context.Background()

	result, err := r.db.NewDelete().
		Model((*Article)(nil)).
		Where("id = ? AND user_id = ?", id, userID).
		Exec(ctx)

	if err != nil {
		return fmt.Errorf("failed to delete article: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("article not found or unauthorized")
	}

	return nil
}

func (r *ArticleRepository) List(limit, offset int) ([]*Article, error) {
	ctx := context.Background()
	var articles []*Article

	err := r.db.NewSelect().
		Model(&articles).
		Where("visibility = ?", VisibilityPublic).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Scan(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to list articles: %w", err)
	}

	return articles, nil
}

func (r *ArticleRepository) GetByUser(userID, viewerID uint64) ([]*Article, error) {
	ctx := context.Background()
	var articles []*Article

	query := r.db.NewSelect().
		Model(&articles).
		Where("user_id = ?", userID).
		Order("created_at DESC")

	if userID != viewerID {
		// Other user viewing - show only public
		query = query.Where("visibility = ?", VisibilityPublic)
	}

	err := query.Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get user articles: %w", err)
	}

	return articles, nil
}

func (r *ArticleRepository) CheckAccess(articleID, viewerID uint64, accessToken string) (bool, error) {
	article, err := r.GetByID(articleID)
	if err != nil {
		return false, err
	}

	// Public articles are accessible to everyone
	if article.Visibility == VisibilityPublic {
		return true, nil
	}

	// Private articles only accessible to owner
	if article.Visibility == VisibilityPrivate {
		return article.UserID == viewerID, nil
	}

	// Link articles accessible with token or to owner
	if article.Visibility == VisibilityLink {
		return article.UserID == viewerID || article.AccessToken == accessToken, nil
	}

	return false, nil
}
