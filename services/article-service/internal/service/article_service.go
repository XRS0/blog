package service

import (
	"context"
	"fmt"
	"log/slog"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/XRS0/blog/services/article-service/internal/repository"
	authpb "github.com/XRS0/blog/services/auth-service/proto"
	"github.com/XRS0/blog/shared/rabbitmq"
)

type ArticleService struct {
	repo       *repository.ArticleRepository
	authClient authpb.AuthServiceClient
	mq         *rabbitmq.Client
	logger     *slog.Logger
}

func NewArticleService(
	repo *repository.ArticleRepository,
	authServiceURL string,
	mq *rabbitmq.Client,
	logger *slog.Logger,
) (*ArticleService, error) {
	// Connect to auth service
	conn, err := grpc.Dial(authServiceURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to auth service: %w", err)
	}

	authClient := authpb.NewAuthServiceClient(conn)

	return &ArticleService{
		repo:       repo,
		authClient: authClient,
		mq:         mq,
		logger:     logger,
	}, nil
}

func (s *ArticleService) Create(ctx context.Context, userID uint64, title, content string, visibility repository.Visibility) (*repository.Article, error) {
	article, err := s.repo.Create(userID, title, content, visibility)
	if err != nil {
		return nil, err
	}

	// Publish event
	event := rabbitmq.Event{
		Type: rabbitmq.EventArticleCreated,
		Data: map[string]interface{}{
			"article_id": article.ID,
			"user_id":    article.UserID,
			"visibility": string(article.Visibility),
		},
	}
	if err := s.mq.Publish(ctx, "articles", "article.created", event); err != nil {
		s.logger.Error("failed to publish article created event", "article_id", article.ID, "error", err)
	}

	s.logger.Info("article created", "article_id", article.ID, "user_id", userID)
	return article, nil
}

func (s *ArticleService) GetByID(ctx context.Context, id, viewerID uint64, accessToken string) (*repository.Article, string, error) {
	article, err := s.repo.GetByID(id)
	if err != nil {
		return nil, "", err
	}

	// Check access
	hasAccess, err := s.repo.CheckAccess(id, viewerID, accessToken)
	if err != nil {
		return nil, "", err
	}
	if !hasAccess {
		return nil, "", fmt.Errorf("access denied")
	}

	// Get author username
	userResp, err := s.authClient.GetUserByID(ctx, &authpb.GetUserByIDRequest{Id: article.UserID})
	if err != nil || userResp.Error != "" {
		s.logger.Error("failed to get author", "user_id", article.UserID, "error", err)
		return article, "", nil
	}

	// Publish view event
	event := rabbitmq.Event{
		Type: rabbitmq.EventArticleViewed,
		Data: map[string]interface{}{
			"article_id": article.ID,
			"user_id":    viewerID,
		},
	}
	if err := s.mq.Publish(ctx, "articles", "article.viewed", event); err != nil {
		s.logger.Error("failed to publish view event", "article_id", article.ID, "error", err)
	}

	return article, userResp.User.Username, nil
}

func (s *ArticleService) Update(ctx context.Context, id, userID uint64, title, content string, visibility repository.Visibility) (*repository.Article, error) {
	return s.repo.Update(id, userID, title, content, visibility)
}

func (s *ArticleService) Delete(ctx context.Context, id, userID uint64) error {
	return s.repo.Delete(id, userID)
}

func (s *ArticleService) List(ctx context.Context, viewerID uint64, limit, offset int) ([]*repository.Article, []string, error) {
	articles, err := s.repo.List(limit, offset)
	if err != nil {
		return nil, nil, err
	}

	// Get author usernames
	usernames := make([]string, len(articles))
	for i, article := range articles {
		userResp, err := s.authClient.GetUserByID(ctx, &authpb.GetUserByIDRequest{Id: article.UserID})
		if err == nil && userResp.Error == "" {
			usernames[i] = userResp.User.Username
		}
	}

	return articles, usernames, nil
}

func (s *ArticleService) GetByUser(ctx context.Context, userID, viewerID uint64) ([]*repository.Article, error) {
	return s.repo.GetByUser(userID, viewerID)
}

func (s *ArticleService) CheckAccess(ctx context.Context, articleID, viewerID uint64, accessToken string) (bool, error) {
	return s.repo.CheckAccess(articleID, viewerID, accessToken)
}
