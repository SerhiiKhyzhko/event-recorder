package handler

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/SerhiiKhyzhko/event-recorder/internal/domain"
	"github.com/gin-gonic/gin"
)

type EventService interface {
	CreateEvent(ctx context.Context, event domain.Event) (domain.Event, error)
	GetEvents(ctx context.Context, userID int64, start, end time.Time, limit, offset int32) ([]domain.Event, error)
}

type EventHandler struct {
	svc EventService
}

func NewEventHandler(svc EventService) *EventHandler {
	return &EventHandler{svc: svc}
}

type CreateEventReq struct {
	UserID   int64           `json:"user_id" binding:"required,gt=0"`
	Action   string          `json:"action" binding:"required"`
	Metadata json.RawMessage `json:"metadata"`
}

func (h *EventHandler) CreateEvent(c *gin.Context) {
	var req CreateEventReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	event := domain.Event{
		UserID:   req.UserID,
		Action:   req.Action,
		Metadata: req.Metadata,
	}

	createdEvent, err := h.svc.CreateEvent(c.Request.Context(), event)
	if err != nil {
		log.Printf("CreateEvent failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create event"})
		return
	}

	c.JSON(http.StatusCreated, createdEvent)
}

type GetEventsReq struct {
	UserID    int64     `form:"user_id" binding:"required,gt=0"`
	StartDate time.Time `form:"start_date" binding:"required" time_format:"2006-01-02T15:04:05Z07:00"`
	EndDate   time.Time `form:"end_date" binding:"required" time_format:"2006-01-02T15:04:05Z07:00"`
	Limit     int32     `form:"limit,default=100"`
	Offset    int32     `form:"offset,default=0"`
}

func (h *EventHandler) GetEvents(c *gin.Context) {
	var req GetEventsReq
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	events, err := h.svc.GetEvents(c.Request.Context(), req.UserID, req.StartDate, req.EndDate, req.Limit, req.Offset)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidDateRange) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		log.Printf("GetEvents failed (user_id=%d): %v", req.UserID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch events"})
		return
	}

	c.JSON(http.StatusOK, events)
}
