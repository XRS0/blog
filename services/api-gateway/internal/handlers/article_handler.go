package handlers

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	articlepb "github.com/XRS0/blog/services/api-gateway/proto/article"
	statspb "github.com/XRS0/blog/services/api-gateway/proto/stats"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Helper function to safely convert protobuf Timestamp to string
func timestampToString(ts *timestamppb.Timestamp) string {
	if ts == nil || !ts.IsValid() {
		return ""
	}
	return ts.AsTime().Format("2006-01-02T15:04:05Z07:00") // RFC3339
}

type ArticleHandler struct {
	articleClient articlepb.ArticleServiceClient
	statsClient   statspb.StatsServiceClient
	logger        *slog.Logger
}

func NewArticleHandler(
	articleClient articlepb.ArticleServiceClient,
	statsClient statspb.StatsServiceClient,
	logger *slog.Logger,
) *ArticleHandler {
	return &ArticleHandler{
		articleClient: articleClient,
		statsClient:   statsClient,
		logger:        logger,
	}
}

type createArticleRequest struct {
	Title      string `json:"title" binding:"required,min=3"`
	Content    string `json:"content" binding:"required,min=10"`
	Visibility string `json:"visibility"` // "public", "private", "link"
}

type updateArticleRequest struct {
	Title      string `json:"title" binding:"required,min=3"`
	Content    string `json:"content" binding:"required,min=10"`
	Visibility string `json:"visibility"`
}

type likeRequest struct {
	Like *bool `json:"like"`
}

func getUserID(c *gin.Context) uint64 {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0
	}
	return userID.(uint64)
}

func visibilityToProto(v string) articlepb.Visibility {
	switch v {
	case "private":
		return articlepb.Visibility_PRIVATE
	case "link":
		return articlepb.Visibility_LINK
	default:
		return articlepb.Visibility_PUBLIC
	}
}

func visibilityFromProto(v articlepb.Visibility) string {
	switch v {
	case articlepb.Visibility_PRIVATE:
		return "private"
	case articlepb.Visibility_LINK:
		return "link"
	default:
		return "public"
	}
}

func (h *ArticleHandler) CreateArticle(c *gin.Context) {
	userID := getUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}

	var req createArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	visibility := visibilityToProto(req.Visibility)

	resp, err := h.articleClient.CreateArticle(context.Background(), &articlepb.CreateArticleRequest{
		UserId:     userID,
		Title:      req.Title,
		Content:    req.Content,
		Visibility: visibility,
	})
	if err != nil {
		h.logger.Error("create article failed", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	if resp.Error != "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": resp.Error})
		return
	}

	article := resp.Article
	result := gin.H{
		"id":         article.Id,
		"title":      article.Title,
		"content":    article.Content,
		"visibility": visibilityFromProto(article.Visibility),
		"created_at": timestampToString(article.CreatedAt),
		"updated_at": timestampToString(article.UpdatedAt),
	}

	// Add access token for link visibility
	if article.Visibility == articlepb.Visibility_LINK && article.AccessToken != "" {
		result["access_token"] = article.AccessToken
		result["access_url"] = "/articles/" + strconv.FormatUint(article.Id, 10) + "?access_token=" + article.AccessToken
	}

	c.JSON(http.StatusCreated, result)
}

func (h *ArticleHandler) GetArticle(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid article id"})
		return
	}

	userID := getUserID(c)
	accessToken := c.Query("access_token")

	resp, err := h.articleClient.GetArticle(context.Background(), &articlepb.GetArticleRequest{
		Id:          id,
		ViewerId:    userID,
		AccessToken: accessToken,
	})
	if err != nil {
		h.logger.Error("get article failed", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	if resp.Error != "" {
		if resp.Error == "access denied" {
			c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		} else if resp.Error == "article not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "article not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": resp.Error})
		}
		return
	}

	// Get stats
	statsResp, err := h.statsClient.GetArticleStats(context.Background(), &statspb.GetArticleStatsRequest{
		ArticleId: id,
	})

	var views, likes uint64
	var viewerLiked bool
	if err == nil && statsResp.Error == "" {
		views = statsResp.Stats.Views
		likes = statsResp.Stats.Likes

		// Check if viewer liked
		if userID > 0 {
			likeResp, err := h.statsClient.GetUserLikeStatus(context.Background(), &statspb.GetUserLikeStatusRequest{
				ArticleId: id,
				UserId:    userID,
			})
			if err == nil && likeResp.Error == "" {
				viewerLiked = likeResp.IsLiked
			}
		}
	}

	article := resp.Article
	c.JSON(http.StatusOK, gin.H{
		"id":          article.Id,
		"title":       article.Title,
		"content":     article.Content,
		"visibility":  visibilityFromProto(article.Visibility),
		"author":      resp.AuthorUsername,
		"views":       views,
		"likes":       likes,
		"viewerLiked": viewerLiked,
		"created_at":  timestampToString(article.CreatedAt),
		"updated_at":  timestampToString(article.UpdatedAt),
	})
}

func (h *ArticleHandler) ListArticles(c *gin.Context) {
	userID := getUserID(c)

	resp, err := h.articleClient.ListArticles(context.Background(), &articlepb.ListArticlesRequest{
		ViewerId: userID,
		Limit:    100,
		Offset:   0,
	})
	if err != nil {
		h.logger.Error("list articles failed", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	if resp.Error != "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": resp.Error})
		return
	}

	// Get stats for all articles
	articleIDs := make([]uint64, len(resp.Articles))
	for i, article := range resp.Articles {
		articleIDs[i] = article.Id
	}

	var statsMap map[uint64]*statspb.ArticleStatsWithLike
	if len(articleIDs) > 0 {
		statsResp, err := h.statsClient.GetArticlesWithStats(context.Background(), &statspb.GetArticlesWithStatsRequest{
			ArticleIds: articleIDs,
			ViewerId:   userID,
		})
		if err == nil && statsResp.Error == "" {
			statsMap = make(map[uint64]*statspb.ArticleStatsWithLike)
			for _, stat := range statsResp.Stats {
				statsMap[stat.ArticleId] = stat
			}
		}
	}

	articles := make([]gin.H, len(resp.Articles))
	for i, article := range resp.Articles {
		var views, likes uint64
		var viewerLiked bool

		if stat, ok := statsMap[article.Id]; ok {
			views = stat.Views
			likes = stat.Likes
			viewerLiked = stat.ViewerLiked
		}

		author := ""
		if i < len(resp.AuthorUsernames) {
			author = resp.AuthorUsernames[i]
		}

		articles[i] = gin.H{
			"id":          article.Id,
			"title":       article.Title,
			"content":     article.Content,
			"visibility":  visibilityFromProto(article.Visibility),
			"author":      author,
			"views":       views,
			"likes":       likes,
			"viewerLiked": viewerLiked,
			"created_at":  timestampToString(article.CreatedAt),
			"updated_at":  timestampToString(article.UpdatedAt),
		}
	}

	c.JSON(http.StatusOK, articles)
}

func (h *ArticleHandler) UpdateArticle(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid article id"})
		return
	}

	userID := getUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}

	var req updateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	visibility := visibilityToProto(req.Visibility)

	resp, err := h.articleClient.UpdateArticle(context.Background(), &articlepb.UpdateArticleRequest{
		Id:         id,
		UserId:     userID,
		Title:      req.Title,
		Content:    req.Content,
		Visibility: visibility,
	})
	if err != nil {
		h.logger.Error("update article failed", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	if resp.Error != "" {
		if resp.Error == "unauthorized" {
			c.JSON(http.StatusForbidden, gin.H{"error": "unauthorized"})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": resp.Error})
		}
		return
	}

	article := resp.Article
	result := gin.H{
		"id":         article.Id,
		"title":      article.Title,
		"content":    article.Content,
		"visibility": visibilityFromProto(article.Visibility),
		"created_at": timestampToString(article.CreatedAt),
		"updated_at": timestampToString(article.UpdatedAt),
	}

	// Add access token for link visibility
	if article.Visibility == articlepb.Visibility_LINK && article.AccessToken != "" {
		result["access_token"] = article.AccessToken
		result["access_url"] = "/articles/" + strconv.FormatUint(article.Id, 10) + "?access_token=" + article.AccessToken
	}

	c.JSON(http.StatusOK, result)
}

func (h *ArticleHandler) DeleteArticle(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid article id"})
		return
	}

	userID := getUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}

	resp, err := h.articleClient.DeleteArticle(context.Background(), &articlepb.DeleteArticleRequest{
		Id:     id,
		UserId: userID,
	})
	if err != nil {
		h.logger.Error("delete article failed", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	if resp.Error != "" {
		if resp.Error == "article not found or unauthorized" {
			c.JSON(http.StatusForbidden, gin.H{"error": "unauthorized"})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": resp.Error})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *ArticleHandler) LikeArticle(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid article id"})
		return
	}

	userID := getUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}

	var req likeRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.Like == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "like field required"})
		return
	}

	if *req.Like {
		// Add like
		resp, err := h.statsClient.RecordLike(context.Background(), &statspb.RecordLikeRequest{
			ArticleId: id,
			UserId:    userID,
		})
		if err != nil || !resp.Success {
			h.logger.Error("record like failed", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}
	} else {
		// Remove like
		resp, err := h.statsClient.RemoveLike(context.Background(), &statspb.RemoveLikeRequest{
			ArticleId: id,
			UserId:    userID,
		})
		if err != nil || !resp.Success {
			h.logger.Error("remove like failed", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}
	}

	// Get updated stats
	statsResp, err := h.statsClient.GetArticleStats(context.Background(), &statspb.GetArticleStatsRequest{
		ArticleId: id,
	})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": true})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"likes":   statsResp.Stats.Likes,
	})
}
