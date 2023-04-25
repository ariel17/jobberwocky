package clients

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ariel17/jobberwocky/internal/core/domain"
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
		{"all jobs without filter", nil, http.StatusOK, `[["Jr Java Developer",24000,"Argentina",["Java","OOP"]],["SSr Java Developer",34000,"Argentina",["Java","OOP","Design Patterns"]],["Sr Java Developer",44000,"Argentina",["Java","OOP","Design Patterns"]]]`, nil, 3, true},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			httpClient := MockHTTPClient{
				StatusCode: tc.status,
				Body:       tc.body,
				Error:      tc.err,
			}
			client := NewJobberwockyExteralJobClient(&httpClient)
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