package server

import (
	"context"
	"log/slog"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/XRS0/blog/services/article-service/internal/repository"
	"github.com/XRS0/blog/services/article-service/internal/service"
	pb "github.com/XRS0/blog/services/article-service/proto/article"
)

type ArticleServer struct {
	pb.UnimplementedArticleServiceServer
	articleService *service.ArticleService
	logger         *slog.Logger
}

func NewArticleServer(articleService *service.ArticleService, logger *slog.Logger) *ArticleServer {
	return &ArticleServer{
		articleService: articleService,
		logger:         logger,
	}
}

func visibilityToProto(v repository.Visibility) pb.Visibility {
	switch v {
	case repository.VisibilityPublic:
		return pb.Visibility_PUBLIC
	case repository.VisibilityPrivate:
		return pb.Visibility_PRIVATE
	case repository.VisibilityLink:
		return pb.Visibility_LINK
	default:
		return pb.Visibility_PUBLIC
	}
}

func visibilityFromProto(v pb.Visibility) repository.Visibility {
	switch v {
	case pb.Visibility_PRIVATE:
		return repository.VisibilityPrivate
	case pb.Visibility_LINK:
		return repository.VisibilityLink
	default:
		return repository.VisibilityPublic
	}
}

func (s *ArticleServer) CreateArticle(ctx context.Context, req *pb.CreateArticleRequest) (*pb.CreateArticleResponse, error) {
	visibility := visibilityFromProto(req.Visibility)
	article, err := s.articleService.Create(ctx, req.UserId, req.Title, req.Content, visibility)
	if err != nil {
		s.logger.Error("create article failed", "user_id", req.UserId, "error", err)
		return &pb.CreateArticleResponse{Error: err.Error()}, nil
	}

	return &pb.CreateArticleResponse{
		Article: &pb.Article{
			Id:          article.ID,
			UserId:      article.UserID,
			Title:       article.Title,
			Content:     article.Content,
			Visibility:  visibilityToProto(article.Visibility),
			AccessToken: article.AccessToken,
			CreatedAt:   timestamppb.New(article.CreatedAt),
			UpdatedAt:   timestamppb.New(article.UpdatedAt),
		},
	}, nil
}

func (s *ArticleServer) GetArticle(ctx context.Context, req *pb.GetArticleRequest) (*pb.GetArticleResponse, error) {
	article, username, err := s.articleService.GetByID(ctx, req.Id, req.ViewerId, req.AccessToken)
	if err != nil {
		s.logger.Error("get article failed", "article_id", req.Id, "error", err)
		return &pb.GetArticleResponse{Error: err.Error()}, nil
	}

	return &pb.GetArticleResponse{
		Article: &pb.Article{
			Id:          article.ID,
			UserId:      article.UserID,
			Title:       article.Title,
			Content:     article.Content,
			Visibility:  visibilityToProto(article.Visibility),
			AccessToken: article.AccessToken,
			CreatedAt:   timestamppb.New(article.CreatedAt),
			UpdatedAt:   timestamppb.New(article.UpdatedAt),
		},
		AuthorUsername: username,
	}, nil
}

func (s *ArticleServer) UpdateArticle(ctx context.Context, req *pb.UpdateArticleRequest) (*pb.UpdateArticleResponse, error) {
	visibility := visibilityFromProto(req.Visibility)
	article, err := s.articleService.Update(ctx, req.Id, req.UserId, req.Title, req.Content, visibility)
	if err != nil {
		s.logger.Error("update article failed", "article_id", req.Id, "error", err)
		return &pb.UpdateArticleResponse{Error: err.Error()}, nil
	}

	return &pb.UpdateArticleResponse{
		Article: &pb.Article{
			Id:          article.ID,
			UserId:      article.UserID,
			Title:       article.Title,
			Content:     article.Content,
			Visibility:  visibilityToProto(article.Visibility),
			AccessToken: article.AccessToken,
			CreatedAt:   timestamppb.New(article.CreatedAt),
			UpdatedAt:   timestamppb.New(article.UpdatedAt),
		},
	}, nil
}

func (s *ArticleServer) DeleteArticle(ctx context.Context, req *pb.DeleteArticleRequest) (*pb.DeleteArticleResponse, error) {
	err := s.articleService.Delete(ctx, req.Id, req.UserId)
	if err != nil {
		s.logger.Error("delete article failed", "article_id", req.Id, "error", err)
		return &pb.DeleteArticleResponse{Success: false, Error: err.Error()}, nil
	}

	return &pb.DeleteArticleResponse{Success: true}, nil
}

func (s *ArticleServer) ListArticles(ctx context.Context, req *pb.ListArticlesRequest) (*pb.ListArticlesResponse, error) {
	limit := int(req.Limit)
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	offset := int(req.Offset)

	articles, usernames, err := s.articleService.List(ctx, req.ViewerId, limit, offset)
	if err != nil {
		s.logger.Error("list articles failed", "error", err)
		return &pb.ListArticlesResponse{Error: err.Error()}, nil
	}

	pbArticles := make([]*pb.Article, len(articles))
	for i, article := range articles {
		pbArticles[i] = &pb.Article{
			Id:          article.ID,
			UserId:      article.UserID,
			Title:       article.Title,
			Content:     article.Content,
			Visibility:  visibilityToProto(article.Visibility),
			AccessToken: article.AccessToken,
			CreatedAt:   timestamppb.New(article.CreatedAt),
			UpdatedAt:   timestamppb.New(article.UpdatedAt),
		}
	}

	return &pb.ListArticlesResponse{
		Articles:        pbArticles,
		AuthorUsernames: usernames,
		Total:           int32(len(articles)),
	}, nil
}

func (s *ArticleServer) GetArticlesByUser(ctx context.Context, req *pb.GetArticlesByUserRequest) (*pb.GetArticlesByUserResponse, error) {
	articles, err := s.articleService.GetByUser(ctx, req.UserId, req.ViewerId)
	if err != nil {
		s.logger.Error("get articles by user failed", "user_id", req.UserId, "error", err)
		return &pb.GetArticlesByUserResponse{Error: err.Error()}, nil
	}

	pbArticles := make([]*pb.Article, len(articles))
	for i, article := range articles {
		pbArticles[i] = &pb.Article{
			Id:          article.ID,
			UserId:      article.UserID,
			Title:       article.Title,
			Content:     article.Content,
			Visibility:  visibilityToProto(article.Visibility),
			AccessToken: article.AccessToken,
			CreatedAt:   timestamppb.New(article.CreatedAt),
			UpdatedAt:   timestamppb.New(article.UpdatedAt),
		}
	}

	return &pb.GetArticlesByUserResponse{Articles: pbArticles}, nil
}

func (s *ArticleServer) CheckArticleAccess(ctx context.Context, req *pb.CheckArticleAccessRequest) (*pb.CheckArticleAccessResponse, error) {
	hasAccess, err := s.articleService.CheckAccess(ctx, req.ArticleId, req.ViewerId, req.AccessToken)
	if err != nil {
		return &pb.CheckArticleAccessResponse{HasAccess: false, Error: err.Error()}, nil
	}

	return &pb.CheckArticleAccessResponse{HasAccess: hasAccess}, nil
}
