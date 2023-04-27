package job

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ariel17/jobberwocky/internal/adapters/clients"
	jobRepository "github.com/ariel17/jobberwocky/internal/adapters/repositories/job"
	"github.com/ariel17/jobberwocky/internal/core/domain"
	jobService "github.com/ariel17/jobberwocky/internal/core/services/job"
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
		{"all jobs from all sources", "", "all", nil, nil, http.StatusOK, 2},
		{"filter by text", "?text=java", "external", nil, nil, http.StatusOK, 1},
		{"filter by description", "?text=work", "local", nil, nil, http.StatusOK, 1},
		{"filter by fixed salary", "?salary=1000", "external", nil, nil, http.StatusOK, 1},
		{"filter by ranged salary", "?salary=2500", "local", nil, nil, http.StatusOK, 1},
		{"filter by single keyword", "?keywords=java", "all", nil, nil, http.StatusOK, 2},
		{"filter by multiple keywords", "?keywords=java&keywords=python", "external", nil, nil, http.StatusOK, 1},
		{"filter by type", "?type=Contractor", "local", nil, nil, http.StatusOK, 1},
		{"filter by company", "?company=SpaceX", "local", nil, nil, http.StatusOK, 1},
		{"filter by is_remote_friendly", "?is_remote_friendly=true", "local", nil, nil, http.StatusOK, 1},
		{"not matching", "?is_remote_friendly=false&company=Avature", "empty", nil, nil, http.StatusOK, 0},
		{"invalid query string", "?salary=invalid", "error_invalid_qs", nil, nil, http.StatusBadRequest, 0},
		{"failed by repository error", "", "error_job_repository", errors.New("mocked job repository error"), nil, http.StatusInternalServerError, 0},
		{"failed by external error", "", "error_external", nil, errors.New("mocked external error"), http.StatusInternalServerError, 0},
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
			router.GET(searchPath, handler.Search)

			req, _ := http.NewRequest(http.MethodGet, searchPath+tc.qs, nil)
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

func getGoldenFile(t *testing.T, name string) string {
	content, err := os.ReadFile(fmt.Sprintf("goldenfiles/%s.json", name))
	assert.Nil(t, err)
	return string(content)
}