package job

import (
	"net/http"

	"github.com/gin-gonic/gin"

	handlers "github.com/ariel17/jobberwocky/internal/adapters/http"
	"github.com/ariel17/jobberwocky/internal/core/domain"
	"github.com/ariel17/jobberwocky/internal/core/ports"
)

const (
	searchPath = "/search"
)

type jobHTTPHandler struct {
	jobService ports.JobService
}

func NewJobHTTPHandler(jobService ports.JobService) *jobHTTPHandler {
	return &jobHTTPHandler{jobService: jobService}
}

func (j *jobHTTPHandler) Search(c *gin.Context) {
	var pattern domain.Pattern
	if err := c.Bind(&pattern); err != nil {
		c.JSON(http.StatusBadRequest, handlers.ErrorResponse{
			Error:       err.Error(),
			Description: "could not parse filter parameters",
		})
		return
	}

	jobs, err := j.jobService.Filter(&pattern)
	if err != nil {
		c.JSON(http.StatusInternalServerError, handlers.ErrorResponse{
			Error:       err.Error(),
			Description: "error searching for jobs",
		})
		return
	}
	c.JSON(http.StatusOK, jobs)
}