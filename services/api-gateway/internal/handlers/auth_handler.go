package handlers

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	authpb "github.com/XRS0/blog/services/api-gateway/proto/auth"
)

type AuthHandler struct {
	authClient authpb.AuthServiceClient
	logger     *slog.Logger
}

func NewAuthHandler(authClient authpb.AuthServiceClient, logger *slog.Logger) *AuthHandler {
	return &AuthHandler{
		authClient: authClient,
		logger:     logger,
	}
}

type registerRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Username string `json:"username" binding:"required,min=3"`
	Password string `json:"password" binding:"required,min=6"`
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.authClient.Register(context.Background(), &authpb.RegisterRequest{
		Email:    req.Email,
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		h.logger.Error("register failed", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	if resp.Error != "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": resp.Error})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"user": gin.H{
			"id":       resp.User.Id,
			"email":    resp.User.Email,
			"username": resp.User.Username,
		},
		"token": resp.Token,
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.authClient.Login(context.Background(), &authpb.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		h.logger.Error("login failed", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	if resp.Error != "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": resp.Error})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":       resp.User.Id,
			"email":    resp.User.Email,
			"username": resp.User.Username,
		},
		"token": resp.Token,
	})
}

func (h *AuthHandler) GetMe(c *gin.Context) {
	userID, _ := c.Get("user_id")
	if userID.(uint64) == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "not authenticated"})
		return
	}

	resp, err := h.authClient.GetUserByID(context.Background(), &authpb.GetUserByIDRequest{
		Id: userID.(uint64),
	})
	if err != nil || resp.Error != "" {
		h.logger.Error("get user failed", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":       resp.User.Id,
		"email":    resp.User.Email,
		"username": resp.User.Username,
	})
}
