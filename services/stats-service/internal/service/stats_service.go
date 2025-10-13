package servicepackage service


import (
	"context"
	"log/slog"

	"github.com/XRS0/blog/services/stats-service/internal/repository"
	"github.com/XRS0/blog/shared/rabbitmq"
)

type StatsService struct {
	repo   *repository.StatsRepository
	mq     *rabbitmq.Client
	logger *slog.Logger
}

func NewStatsService(repo *repository.StatsRepository, mq *rabbitmq.Client, logger *slog.Logger) *StatsService {
	return &StatsService{
		repo:   repo,
		mq:     mq,
		logger: logger,
	}
}

func (s *StatsService) StartEventConsumer(ctx context.Context) error {
	// Declare queue for stats events
	if err := s.mq.DeclareQueue("stats-events"); err != nil {
		return err
	}

	// Bind to articles exchange
	if err := s.mq.BindQueue("stats-events", "articles", "article.*"); err != nil {
		return err
	}

	// Start consuming
	return s.mq.Consume("stats-events", s.handleEvent)
}

func (s *StatsService) handleEvent(event rabbitmq.Event) error {
	s.logger.Debug("handling event", "type", event.Type)

	switch event.Type {
	case rabbitmq.EventArticleViewed:
		articleID := uint64(event.Data["article_id"].(float64))
		userID := uint64(0)
		if uid, ok := event.Data["user_id"]; ok {
			userID = uint64(uid.(float64))
		}
		return s.repo.RecordView(articleID, userID)

	case rabbitmq.EventArticleLiked:
		articleID := uint64(event.Data["article_id"].(float64))
		userID := uint64(event.Data["user_id"].(float64))
		return s.repo.RecordLike(articleID, userID)

	case rabbitmq.EventArticleUnliked:
		articleID := uint64(event.Data["article_id"].(float64))
		userID := uint64(event.Data["user_id"].(float64))
		return s.repo.RemoveLike(articleID, userID)
	}

	return nil
}

func (s *StatsService) RecordView(articleID, userID uint64) error {
	return s.repo.RecordView(articleID, userID)
}

func (s *StatsService) RecordLike(articleID, userID uint64) error {
	return s.repo.RecordLike(articleID, userID)
}

func (s *StatsService) RemoveLike(articleID, userID uint64) error {
	return s.repo.RemoveLike(articleID, userID)
}

func (s *StatsService) GetArticleStats(articleID uint64) (views, likes uint64, err error) {
	return s.repo.GetArticleStats(articleID)
}

func (s *StatsService) IsLikedByUser(articleID, userID uint64) (bool, error) {
	return s.repo.IsLikedByUser(articleID, userID)
}

func (s *StatsService) GetMultipleArticleStats(articleIDs []uint64, viewerID uint64) ([]repository.ArticleStats, error) {
	return s.repo.GetMultipleArticleStats(articleIDs, viewerID)
}
