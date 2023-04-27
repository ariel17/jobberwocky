package clients

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ariel17/jobberwocky/internal/core/domain"
	helpers "github.com/ariel17/jobberwocky/internal/internal_test"
)

func TestExternalJobClient_Filter(t *testing.T) {
	testCases := []struct {
		name         string
		pattern      *domain.Pattern
		status       int
		body         string
		err          error
		expectedJobs int
		success      bool
	}{
		{"all jobs without filter", nil, http.StatusOK, `[["Jr Java Developer",24000,"Argentina",["Java","Python","Kotlin"]],["SSr Java Developer",34000,"Argentina",["Java","OOP","Design Patterns"]],["Sr Java Developer",44000,"Argentina",["Java","OOP","Design Patterns"]]]`, nil, 3, true},
		{"matching all keywords", &domain.Pattern{Keywords: []string{"Java", "Python", "Kotlin"}}, http.StatusOK, `[["Jr Java Developer",24000,"Argentina",["Java","Python","Kotlin"]],["SSr Java Developer",34000,"Argentina",["Java","OOP","Design Patterns"]],["Sr Java Developer",44000,"Argentina",["Java","OOP","Design Patterns"]]]`, nil, 1, true},
		{"not matching", nil, http.StatusOK, `[]`, nil, 0, true},
		{"not matching by not searchable fields", &domain.Pattern{Company: "Company"}, http.StatusOK, `[]`, nil, 0, true},
		{"invalid body", nil, http.StatusOK, `{invalid`, nil, 0, false},
		{"api error", nil, http.StatusInternalServerError, ``, nil, 0, false},
		{"network error", nil, 0, ``, errors.New("mocked network error"), 0, false},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			httpClient := MockHTTPClient{
				StatusCode: tc.status,
				Body:       tc.body,
				Error:      tc.err,
			}
			client := NewJobberwockyExternalJobClient(&httpClient, "")
			jobs, err := client.Filter(tc.pattern)
			if tc.success {
				assert.NotNil(t, jobs)
				assert.Nil(t, err)
				assert.Equal(t, tc.expectedJobs, len(jobs))
				for _, job := range jobs {
					assert.NotEmpty(t, job.Title)
					assert.Empty(t, job.Description)
					assert.NotEmpty(t, job.Location)
					assert.Empty(t, job.SalaryMin)
					assert.NotEmpty(t, job.SalaryMax)
					assert.NotEmpty(t, job.Keywords)
				}
			} else {
				assert.Nil(t, jobs)
				assert.NotNil(t, err)
			}
		})
	}
}

func TestExternalJobClient_PatternIsSearchable(t *testing.T) {
	testCases := []struct {
		name     string
		pattern  *domain.Pattern
		expected bool
	}{
		{"empty filter is searchable v1", &domain.Pattern{}, true},
		{"empty filter is searchable v2", &domain.Pattern{}, true},
		{"filter by type is not searchable", &domain.Pattern{Type: domain.FullTime}, false},
		{"filter by company is not searchable", &domain.Pattern{Company: "SpaceX"}, false},
		{"filter by is_remote_friendly is not searchable", &domain.Pattern{IsRemoteFriendly: helpers.BoolPointer(true)}, false},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			client := NewJobberwockyExternalJobClient(nil, "")
			assert.Equal(t, tc.expected, client.PatternIsSearchable(tc.pattern))
		})
	}
}

func TestPatternToQueryString(t *testing.T) {
	testCases := []struct {
		name     string
		pattern  *domain.Pattern
		expected string
	}{
		{"without pattern", nil, ""},
		{"empty", &domain.Pattern{}, ""},
		{"text", &domain.Pattern{Text: "text"}, "?name=text"},
		{"location", &domain.Pattern{Location: "location"}, "?country=location"},
		{"salary", &domain.Pattern{Salary: 10}, "?salary_min=10&salary_max=10"},
		{"all fields", &domain.Pattern{Text: "text", Location: "location", Salary: 10}, "?name=text&country=location&salary_min=10&salary_max=10"},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			qs := patternToQueryString(tc.pattern)
			assert.Equal(t, tc.expected, qs)
		})
	}
}