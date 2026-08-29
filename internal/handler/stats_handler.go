package handler

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/SerhiiKhyzhko/event-recorder/internal/domain"
	"github.com/gin-gonic/gin"
)

type StatsService interface {
	GetStats(ctx context.Context, userID int64, start, end time.Time) ([]domain.UserStat, error)
}

type StatsHandler struct {
	svc StatsService
}

func NewStatsHandler(svc StatsService) *StatsHandler {
	return &StatsHandler{svc: svc}
}

type GetStatsReq struct {
	UserID    int64     `form:"user_id" binding:"required,gt=0"`
	StartDate time.Time `form:"start_date" binding:"required" time_format:"2006-01-02T15:04:05Z07:00"`
	EndDate   time.Time `form:"end_date" binding:"required" time_format:"2006-01-02T15:04:05Z07:00"`
}

func (h *StatsHandler) GetStats(c *gin.Context) {
	var req GetStatsReq
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	stats, err := h.svc.GetStats(c.Request.Context(), req.UserID, req.StartDate, req.EndDate)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidDateRange) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		log.Printf("GetStats failed (user_id=%d): %v", req.UserID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch stats"})
		return
	}

	c.JSON(http.StatusOK, stats)
}