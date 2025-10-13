package server

import (
	"context"
	"log/slog"

	"github.com/XRS0/blog/services/stats-service/internal/service"
	pb "github.com/XRS0/blog/services/stats-service/proto"
)

type StatsServer struct {
	pb.UnimplementedStatsServiceServer
	statsService *service.StatsService
	logger       *slog.Logger
}

func NewStatsServer(statsService *service.StatsService, logger *slog.Logger) *StatsServer {
	return &StatsServer{
		statsService: statsService,
		logger:       logger,
	}
}

func (s *StatsServer) RecordView(ctx context.Context, req *pb.RecordViewRequest) (*pb.RecordViewResponse, error) {
	err := s.statsService.RecordView(req.ArticleId, req.UserId)
	if err != nil {
		return &pb.RecordViewResponse{Success: false, Error: err.Error()}, nil
	}
	return &pb.RecordViewResponse{Success: true}, nil
}

func (s *StatsServer) RecordLike(ctx context.Context, req *pb.RecordLikeRequest) (*pb.RecordLikeResponse, error) {
	err := s.statsService.RecordLike(req.ArticleId, req.UserId)
	if err != nil {
		return &pb.RecordLikeResponse{Success: false, Error: err.Error()}, nil
	}
	return &pb.RecordLikeResponse{Success: true}, nil
}

func (s *StatsServer) RemoveLike(ctx context.Context, req *pb.RemoveLikeRequest) (*pb.RemoveLikeResponse, error) {
	err := s.statsService.RemoveLike(req.ArticleId, req.UserId)
	if err != nil {
		return &pb.RemoveLikeResponse{Success: false, Error: err.Error()}, nil
	}
	return &pb.RemoveLikeResponse{Success: true}, nil
}

func (s *StatsServer) GetArticleStats(ctx context.Context, req *pb.GetArticleStatsRequest) (*pb.GetArticleStatsResponse, error) {
	views, likes, err := s.statsService.GetArticleStats(req.ArticleId)
	if err != nil {
		return &pb.GetArticleStatsResponse{Error: err.Error()}, nil
	}

	return &pb.GetArticleStatsResponse{
		Stats: &pb.ArticleStats{
			ArticleId: req.ArticleId,
			Views:     views,
			Likes:     likes,
		},
	}, nil
}

func (s *StatsServer) GetUserLikeStatus(ctx context.Context, req *pb.GetUserLikeStatusRequest) (*pb.GetUserLikeStatusResponse, error) {
	isLiked, err := s.statsService.IsLikedByUser(req.ArticleId, req.UserId)
	if err != nil {
		return &pb.GetUserLikeStatusResponse{Error: err.Error()}, nil
	}

	return &pb.GetUserLikeStatusResponse{IsLiked: isLiked}, nil
}

func (s *StatsServer) GetArticlesWithStats(ctx context.Context, req *pb.GetArticlesWithStatsRequest) (*pb.GetArticlesWithStatsResponse, error) {
	stats, err := s.statsService.GetMultipleArticleStats(req.ArticleIds, req.ViewerId)
	if err != nil {
		return &pb.GetArticlesWithStatsResponse{Error: err.Error()}, nil
	}

	pbStats := make([]*pb.ArticleStatsWithLike, len(stats))
	for i, stat := range stats {
		pbStats[i] = &pb.ArticleStatsWithLike{
			ArticleId:   stat.ArticleID,
			Views:       stat.Views,
			Likes:       stat.Likes,
			ViewerLiked: stat.ViewerLiked,
		}
	}

	return &pb.GetArticlesWithStatsResponse{Stats: pbStats}, nil
}
