package job

import (
	"net/http"

	"github.com/gin-gonic/gin"

	handlers "github.com/ariel17/jobberwocky/internal/adapters/http"
	"github.com/ariel17/jobberwocky/internal/core/domain"
	"github.com/ariel17/jobberwocky/internal/core/ports"
)

const (
	jobsPath = "/jobs"
)

type jobHTTPHandler struct {
	handlers.Handler
	jobService ports.JobService
}

func NewJobHTTPHandler(s ports.JobService) handlers.Handler {
	return &jobHTTPHandler{jobService: s}
}

// Search handles requests for searching jobs in the local database and in external resources.
// @Summary      Search for published jobs
// @Description  Based on filter parameters it searches in jobs in the local database and in external resources concurrently.
// @Tags         jobs
// @Param        text               query     string      false  "Filters jobs by matching text in title or description (case-insensitive)."
// @Param        company            query     string      false  "Filters jobs by matching company (case-sensitive)."
// @Param        location           query     string      false  "Filters jobs by matching location (case-sensitive)."
// @Param        salary             query     int         false  "Filters jobs by matching salary, fixed or in range."
// @Param        type               query     string      false  "Filters jobs by matching work type (case-sensitive). Values: Full-Time, Contractor, Part-Time."
// @Param        is_remote_friendly query     bool        false  "Filters jobs by remote condition."
// @Param        keywords           query     []string    false  "Filters jobs by keywords (case-sensitive, inclusive)."
// @Accept       json
// @Produce      json
// @Success      200  {object}  []domain.Job
// @Failure      400  {object}  http.ErrorResponse
// @Failure      500  {object}  http.ErrorResponse
// @Router       /jobs [get]
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

// Post is the handler that supports publishing new jobs.
// @Summary      Publish a new job
// @Description  Creates a new job receiving a JSON body with the details. If matching subscriptions exists, it sends notifications by email asynchronously.
// @Tags         jobs
// @Param        job	body		domain.Job	true	"New job details."
// @Accept       json
// @Produce      json
// @Success      201  {object}	domain.Job
// @Failure      400  {object}  http.ErrorResponse
// @Failure      500  {object}  http.ErrorResponse
// @Router       /jobs [post]
func (j *jobHTTPHandler) Post(c *gin.Context) {
	var job domain.Job
	if err := c.BindJSON(&job); err != nil {
		c.JSON(http.StatusBadRequest, handlers.ErrorResponse{
			Error:       err.Error(),
			Description: "could not parse job to publish",
		})
		return
	}

	if err := job.Validate(true); err != nil {
		c.JSON(http.StatusBadRequest, handlers.ErrorResponse{
			Error:       err.Error(),
			Description: "job is not valid",
		})
		return
	}

	if err := j.jobService.Create(job); err != nil {
		c.JSON(http.StatusInternalServerError, handlers.ErrorResponse{
			Error:       err.Error(),
			Description: "error saving new job",
		})
		return
	}
	c.Status(http.StatusCreated)
}

func (j *jobHTTPHandler) ConfigureRoutes(router *gin.Engine) {
	router.GET(jobsPath, j.Search)
	router.POST(jobsPath, j.Post)
}