package subscription

import (
	"net/http"

	"github.com/gin-gonic/gin"

	handlers "github.com/ariel17/jobberwocky/internal/adapters/http"
	"github.com/ariel17/jobberwocky/internal/core/domain"
	"github.com/ariel17/jobberwocky/internal/core/ports"
)

const (
	subscriptionPath = "/subscriptions"
)

type subscriptionHTTPHandler struct {
	handlers.Handler
	subscriptionService ports.SubscriptionService
}

func NewSubscriptionHTTPHandler(s ports.SubscriptionService) *subscriptionHTTPHandler {
	return &subscriptionHTTPHandler{subscriptionService: s}
}

// Post handles the request to create a new subscription to be notified on new job posts.
// @Summary      Creates a new subscripion
// @Description  Receives a JSON body with the email and filter values to match new job posts and be notified.
// @Tags         jobs
// @Accept       json
// @Produce      json
// @Param        subscription	body		domain.Subscription	true	"New subscription details."
// @Success      201  {object}	domain.Subscription
// @Failure      400  {object}  http.ErrorResponse
// @Failure      500  {object}  http.ErrorResponse
// @Router       /subscriptions [post]
func (s *subscriptionHTTPHandler) Post(c *gin.Context) {
	var subscription domain.Subscription
	if err := c.BindJSON(&subscription); err != nil {
		c.JSON(http.StatusBadRequest, handlers.ErrorResponse{
			Error:       err.Error(),
			Description: "could not parse subscription",
		})
		return
	}

	if err := subscription.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, handlers.ErrorResponse{
			Error:       err.Error(),
			Description: "subscription is invalid",
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

func (s *subscriptionHTTPHandler) ConfigureRoutes(router *gin.Engine) {
	router.POST(subscriptionPath, s.Post)
}