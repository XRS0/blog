package server

import (
	"context"
	"log/slog"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/XRS0/blog/services/auth-service/internal/service"
	pb "github.com/XRS0/blog/services/auth-service/proto"
)

type AuthServer struct {
	pb.UnimplementedAuthServiceServer
	authService *service.AuthService
	logger      *slog.Logger
}

func NewAuthServer(authService *service.AuthService, logger *slog.Logger) *AuthServer {
	return &AuthServer{
		authService: authService,
		logger:      logger,
	}
}

func (s *AuthServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	user, token, err := s.authService.Register(req.Email, req.Username, req.Password)
	if err != nil {
		s.logger.Error("registration failed", "email", req.Email, "error", err)
		return &pb.RegisterResponse{Error: err.Error()}, nil
	}

	return &pb.RegisterResponse{
		User: &pb.User{
			Id:        user.ID,
			Email:     user.Email,
			Username:  user.Username,
			CreatedAt: timestamppb.New(user.CreatedAt),
			UpdatedAt: timestamppb.New(user.UpdatedAt),
		},
		Token: token,
	}, nil
}

func (s *AuthServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	user, token, err := s.authService.Login(req.Email, req.Password)
	if err != nil {
		s.logger.Error("login failed", "email", req.Email, "error", err)
		return &pb.LoginResponse{Error: err.Error()}, nil
	}

	return &pb.LoginResponse{
		User: &pb.User{
			Id:        user.ID,
			Email:     user.Email,
			Username:  user.Username,
			CreatedAt: timestamppb.New(user.CreatedAt),
			UpdatedAt: timestamppb.New(user.UpdatedAt),
		},
		Token: token,
	}, nil
}

func (s *AuthServer) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	userID, err := s.authService.ValidateToken(req.Token)
	if err != nil {
		return &pb.ValidateTokenResponse{
			Valid: false,
			Error: err.Error(),
		}, nil
	}

	return &pb.ValidateTokenResponse{
		Valid:  true,
		UserId: userID,
	}, nil
}

func (s *AuthServer) GetUserByID(ctx context.Context, req *pb.GetUserByIDRequest) (*pb.GetUserByIDResponse, error) {
	user, err := s.authService.GetUserByID(req.Id)
	if err != nil {
		return &pb.GetUserByIDResponse{Error: err.Error()}, nil
	}

	return &pb.GetUserByIDResponse{
		User: &pb.User{
			Id:        user.ID,
			Email:     user.Email,
			Username:  user.Username,
			CreatedAt: timestamppb.New(user.CreatedAt),
			UpdatedAt: timestamppb.New(user.UpdatedAt),
		},
	}, nil
}

func (s *AuthServer) GetUserByEmail(ctx context.Context, req *pb.GetUserByEmailRequest) (*pb.GetUserByEmailResponse, error) {
	user, err := s.authService.GetUserByEmail(req.Email)
	if err != nil {
		return &pb.GetUserByEmailResponse{Error: err.Error()}, nil
	}

	return &pb.GetUserByEmailResponse{
		User: &pb.User{
			Id:        user.ID,
			Email:     user.Email,
			Username:  user.Username,
			CreatedAt: timestamppb.New(user.CreatedAt),
			UpdatedAt: timestamppb.New(user.UpdatedAt),
		},
	}, nil
}
