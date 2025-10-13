package client

import (
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	articlepb "github.com/XRS0/blog/services/api-gateway/proto/article"
	authpb "github.com/XRS0/blog/services/api-gateway/proto/auth"
	statspb "github.com/XRS0/blog/services/api-gateway/proto/stats"
)

type ServiceClients struct {
	Auth    authpb.AuthServiceClient
	Article articlepb.ArticleServiceClient
	Stats   statspb.StatsServiceClient
}

func NewServiceClients(authURL, articleURL, statsURL string) (*ServiceClients, error) {
	// Connect to Auth Service
	authConn, err := grpc.Dial(authURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to auth service: %w", err)
	}

	// Connect to Article Service
	articleConn, err := grpc.Dial(articleURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		authConn.Close()
		return nil, fmt.Errorf("failed to connect to article service: %w", err)
	}

	// Connect to Stats Service
	statsConn, err := grpc.Dial(statsURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		authConn.Close()
		articleConn.Close()
		return nil, fmt.Errorf("failed to connect to stats service: %w", err)
	}

	return &ServiceClients{
		Auth:    authpb.NewAuthServiceClient(authConn),
		Article: articlepb.NewArticleServiceClient(articleConn),
		Stats:   statspb.NewStatsServiceClient(statsConn),
	}, nil
}
