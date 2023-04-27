package job

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/ariel17/jobberwocky/internal/adapters/clients"
	jobRepository "github.com/ariel17/jobberwocky/internal/adapters/repositories/job"
	"github.com/ariel17/jobberwocky/internal/adapters/repositories/subscription"
	"github.com/ariel17/jobberwocky/internal/core/domain"
	jobService "github.com/ariel17/jobberwocky/internal/core/services/job"
	"github.com/ariel17/jobberwocky/internal/core/services/notification"
	helpers "github.com/ariel17/jobberwocky/internal/internal_test"

	"github.com/gin-gonic/gin"
)

func TestJobHTTPHandler_Search(t *testing.T) {
	localJobs := []domain.Job{{Title: "Junior Kotlin developer", Description: "We need you to make the work.", Company: "SpaceX", Location: "Argentina", SalaryMin: 2000, SalaryMax: 3000, Type: domain.Contractor, IsRemoteFriendly: helpers.BoolPointer(true), Keywords: []string{"kotlin", "java"}}}
	externalJobs := []domain.Job{{Title: "Looking for a Sr Java developer", Location: "USA", SalaryMax: 1000, Keywords: []string{"python", "java", "golang"}, Source: "Avature"}}
	testCases := []struct {
		name         string
		qs           string
		goldenPath   string
		jobErr       error
		externalErr  error
		statusCode   int
		matchingJobs int
	}{
		{"all jobs from all sources", "", "search_all", nil, nil, http.StatusOK, 2},
		{"filter by text", "?text=java", "search_external", nil, nil, http.StatusOK, 1},
		{"filter by description", "?text=work", "search_local", nil, nil, http.StatusOK, 1},
		{"filter by fixed salary", "?salary=1000", "search_external", nil, nil, http.StatusOK, 1},
		{"filter by ranged salary", "?salary=2500", "search_local", nil, nil, http.StatusOK, 1},
		{"filter by single keyword", "?keywords=java", "search_all", nil, nil, http.StatusOK, 2},
		{"filter by multiple keywords", "?keywords=java&keywords=python", "search_external", nil, nil, http.StatusOK, 1},
		{"filter by type", "?type=Contractor", "search_local", nil, nil, http.StatusOK, 1},
		{"filter by company", "?company=SpaceX", "search_local", nil, nil, http.StatusOK, 1},
		{"filter by is_remote_friendly", "?is_remote_friendly=true", "search_local", nil, nil, http.StatusOK, 1},
		{"not matching", "?is_remote_friendly=false&company=Avature", "search_empty", nil, nil, http.StatusOK, 0},
		{"invalid query string", "?salary=invalid", "search_error_invalid_qs", nil, nil, http.StatusBadRequest, 0},
		{"failed by repository error", "", "search_error_job_repository", errors.New("mocked job repository error"), nil, http.StatusInternalServerError, 0},
		{"failed by external error", "", "search_error_external", nil, errors.New("mocked external error"), http.StatusInternalServerError, 0},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			jr := jobRepository.MockJobRepository{
				MockFilter: helpers.MockFilter{
					Jobs:  localJobs,
					Error: tc.jobErr,
				},
			}
			e := clients.MockExternalJobClient{
				MockFilter: helpers.MockFilter{
					Jobs:  externalJobs,
					Error: tc.externalErr,
				},
			}
			jobService := jobService.NewJobService(&jr, nil, &e)

			handler := NewJobHTTPHandler(jobService)
			router := gin.Default()
			router.GET(jobsPath, handler.Search)

			req, _ := http.NewRequest(http.MethodGet, jobsPath+tc.qs, nil)
			rr := httptest.NewRecorder()

			router.ServeHTTP(rr, req)
			assert.Equal(t, tc.statusCode, rr.Code)

			expected := getGoldenFile(t, tc.goldenPath)
			assert.Equal(t, expected, rr.Body.String())

			if tc.statusCode == http.StatusOK {
				var jobs []domain.Job
				assert.Nil(t, json.Unmarshal(rr.Body.Bytes(), &jobs))
				assert.Equal(t, tc.matchingJobs, len(jobs))
			}
		})
	}

}

func TestJobHTTPHandler_Post(t *testing.T) {
	subscriptions := []domain.Subscription{
		{domain.Pattern{}, "person1@example.com"},
		{domain.Pattern{Text: "Kotlin"}, "person2@example.com"},
		{domain.Pattern{Type: domain.FullTime}, "person3@example.com"},
	}

	testCases := []struct {
		name                      string
		goldenPathBody            string
		goldenPathResponse        string
		jobRepositoryErr          error
		subscriptionRepositoryErr error
		emailErr                  error
		statusCode                int
		matchingSubscribers       int
	}{
		{"new job published", "post_body_new", "", nil, nil, nil, http.StatusCreated, 2},
		{"failed by job repository", "post_body_new", "post_response_job_repository_error", errors.New("mocked job repository error"), nil, nil, http.StatusInternalServerError, 0},
		{"failed by invalid body", "post_body_invalid", "post_response_invalid", nil, nil, nil, http.StatusBadRequest, 0},
		{"partial success by subscription repository error", "post_body_new", "", nil, errors.New("mocked subscription repository error"), nil, http.StatusCreated, 0},
		{"partial success by email client error", "post_body_new", "", nil, nil, errors.New("mocked email client error"), http.StatusCreated, 2},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			jr := jobRepository.MockJobRepository{
				MockFilter: helpers.MockFilter{
					Error: tc.jobRepositoryErr,
				},
			}

			sr := subscription.MockSubscriptionRepository{
				MockFilter: helpers.MockFilter{
					Error: tc.subscriptionRepositoryErr,
				},
				MockSave:      helpers.MockSave{},
				Subscriptions: subscriptions,
			}

			ec := clients.MockEmailProviderClient{
				Error: tc.emailErr,
			}

			ns := notification.NewNotificationService(10, &sr, &ec, getTemplatePath(t))
			ns.StartWorkers()
			defer ns.StopWorkers()

			jobService := jobService.NewJobService(&jr, ns, nil)

			handler := NewJobHTTPHandler(jobService)
			router := gin.Default()
			router.POST(jobsPath, handler.Post)

			body := getGoldenFile(t, tc.goldenPathBody)
			req, _ := http.NewRequest(http.MethodPost, jobsPath, strings.NewReader(body))
			rr := httptest.NewRecorder()

			router.ServeHTTP(rr, req)
			assert.Equal(t, tc.statusCode, rr.Code)

			time.Sleep(time.Millisecond)

			if tc.statusCode != http.StatusCreated {
				expected := getGoldenFile(t, tc.goldenPathResponse)
				assert.Equal(t, expected, rr.Body.String())
			} else {
				assert.Equal(t, tc.matchingSubscribers, ec.SendCalls())
			}
		})
	}
}

func getGoldenFile(t *testing.T, name string) string {
	content, err := os.ReadFile(fmt.Sprintf("goldenfiles/%s.json", name))
	assert.Nil(t, err)
	return string(content)
}

func getTemplatePath(t *testing.T) string {
	wd, err := os.Getwd()
	assert.Nil(t, err)
	return wd + "/../../../../resources/body.tmpl"
}