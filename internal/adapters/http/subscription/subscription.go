package subscription

import (
	"net/http"

	"github.com/gin-gonic/gin"

	handlers "github.com/ariel17/jobberwocky/internal/adapters/http"
	"github.com/ariel17/jobberwocky/internal/core/domain"
	"github.com/ariel17/jobberwocky/internal/core/ports"
)

const (
	subscriptionPath = "/subscription"
)

type subscriptionHTTPHandler struct {
	subscriptionService ports.SubscriptionService
}

func NewSubscriptionHTTPHandler(s ports.SubscriptionService) *subscriptionHTTPHandler {
	return &subscriptionHTTPHandler{subscriptionService: s}
}

func (s *subscriptionHTTPHandler) Post(c *gin.Context) {
	var subscription domain.Subscription
	if err := c.BindJSON(&subscription); err != nil {
		c.JSON(http.StatusBadRequest, handlers.ErrorResponse{
			Error:       err.Error(),
			Description: "could not parse subscription",
		})
		return
	}

	if err := s.subscriptionService.Create(subscription); err != nil {
		c.JSON(http.StatusInternalServerError, handlers.ErrorResponse{
			Error:       err.Error(),
			Description: "error saving new subscription",
		})
		return
	}
	c.Status(http.StatusCreated)
}