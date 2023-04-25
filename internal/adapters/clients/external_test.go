package clients

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ariel17/jobberwocky/internal/core/domain"
)

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
		{"all jobs without filter", nil, http.StatusOK, `[["Jr Java Developer",24000,"Argentina",["Java","OOP"]],["SSr Java Developer",34000,"Argentina",["Java","OOP","Design Patterns"]],["Sr Java Developer",44000,"Argentina",["Java","OOP","Design Patterns"]]]`, nil, 3, true},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			httpClient := MockHTTPClient{
				StatusCode: tc.status,
				Body:       tc.body,
				Error:      tc.err,
			}
			client := NewJobberwockyExternalJobClient(&httpClient)
			jobs, err := client.Filter(tc.pattern)
			if tc.status == http.StatusOK {
				assert.NotNil(t, jobs)
				assert.Nil(t, err)
				assert.Equal(t, tc.expectedJobs, len(jobs))
			} else {
				assert.Nil(t, jobs)
				assert.NotNil(t, err)
			}
		})
	}
}