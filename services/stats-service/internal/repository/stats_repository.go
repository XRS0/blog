package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/uptrace/bun"
)

type ArticleView struct {
	bun.BaseModel `bun:"table:article_views,alias:av"`

	ID        uint64    `bun:"id,pk,autoincrement"`
	ArticleID uint64    `bun:"article_id,notnull"`
	UserID    *uint64   `bun:"user_id"`
	CreatedAt time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp"`
}

type ArticleLike struct {
	bun.BaseModel `bun:"table:article_likes,alias:al"`

	ID        uint64    `bun:"id,pk,autoincrement"`
	ArticleID uint64    `bun:"article_id,notnull"`
	UserID    uint64    `bun:"user_id,notnull"`
	CreatedAt time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp"`
}

type StatsRepository struct {
	db *bun.DB
}

func NewStatsRepository(db *bun.DB) *StatsRepository {
	return &StatsRepository{db: db}
}

func (r *StatsRepository) RecordView(articleID, userID uint64) error {
	ctx := context.Background()

	view := &ArticleView{
		ArticleID: articleID,
		CreatedAt: time.Now(),
	}

	if userID > 0 {
		view.UserID = &userID
	}

	_, err := r.db.NewInsert().Model(view).Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to record view: %w", err)
	}
	return nil
}

func (r *StatsRepository) RecordLike(articleID, userID uint64) error {
	ctx := context.Background()

	like := &ArticleLike{
		ArticleID: articleID,
		UserID:    userID,
		CreatedAt: time.Now(),
	}

	_, err := r.db.NewInsert().
		Model(like).
		On("CONFLICT (article_id, user_id) DO NOTHING").
		Exec(ctx)

	if err != nil {
		return fmt.Errorf("failed to record like: %w", err)
	}
	return nil
}

func (r *StatsRepository) RemoveLike(articleID, userID uint64) error {
	ctx := context.Background()

	_, err := r.db.NewDelete().
		Model((*ArticleLike)(nil)).
		Where("article_id = ? AND user_id = ?", articleID, userID).
		Exec(ctx)

	if err != nil {
		return fmt.Errorf("failed to remove like: %w", err)
	}
	return nil
}

func (r *StatsRepository) GetArticleStats(articleID uint64) (views, likes uint64, err error) {
	ctx := context.Background()

	// Get views count
	viewCount, err := r.db.NewSelect().
		Model((*ArticleView)(nil)).
		Where("article_id = ?", articleID).
		Count(ctx)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get views: %w", err)
	}

	// Get likes count
	likeCount, err := r.db.NewSelect().
		Model((*ArticleLike)(nil)).
		Where("article_id = ?", articleID).
		Count(ctx)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get likes: %w", err)
	}

	return uint64(viewCount), uint64(likeCount), nil
}

func (r *StatsRepository) IsLikedByUser(articleID, userID uint64) (bool, error) {
	ctx := context.Background()

	exists, err := r.db.NewSelect().
		Model((*ArticleLike)(nil)).
		Where("article_id = ? AND user_id = ?", articleID, userID).
		Exists(ctx)

	return exists, err
}

type ArticleStats struct {
	ArticleID   uint64
	Views       uint64
	Likes       uint64
	ViewerLiked bool
}

func (r *StatsRepository) GetMultipleArticleStats(articleIDs []uint64, viewerID uint64) ([]ArticleStats, error) {
	if len(articleIDs) == 0 {
		return []ArticleStats{}, nil
	}

	ctx := context.Background()
	stats := make([]ArticleStats, 0, len(articleIDs))

	for _, articleID := range articleIDs {
		// Get views
		viewCount, _ := r.db.NewSelect().
			Model((*ArticleView)(nil)).
			Where("article_id = ?", articleID).
			Count(ctx)

		// Get likes
		likeCount, _ := r.db.NewSelect().
			Model((*ArticleLike)(nil)).
			Where("article_id = ?", articleID).
			Count(ctx)

		// Check if viewer liked
		viewerLiked := false
		if viewerID > 0 {
			viewerLiked, _ = r.db.NewSelect().
				Model((*ArticleLike)(nil)).
				Where("article_id = ? AND user_id = ?", articleID, viewerID).
				Exists(ctx)
		}

		stats = append(stats, ArticleStats{
			ArticleID:   articleID,
			Views:       uint64(viewCount),
			Likes:       uint64(likeCount),
			ViewerLiked: viewerLiked,
		})
	}

	return stats, nil
}
