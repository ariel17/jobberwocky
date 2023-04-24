package services

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ariel17/jobberwocky/internal/adapters/repositories"
	"github.com/ariel17/jobberwocky/internal/core/domain"
)

func TestJobService_Create(t *testing.T) {
	testCases := []struct {
		name string
		err  error
	}{
		{"job created", nil},
		{"failed by repository error", errors.New("mocked error")},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repository := repositories.MockJobRepository{
				Error: tc.err,
			}
			notification := MockNotificationService{}
			service := NewJobService(&repository, &notification)
			job := domain.Job{
				Title:            "Looking for a Technical Leader",
				Description:      "This is the longest part where we describe all the details about the job and required skills.",
				Company:          "Ariel Labs",
				Location:         "Argentina",
				SalaryMin:        6000,
				SalaryMax:        8000,
				Type:             "Full-Time",
				IsRemoteFriendly: true,
				Keywords:         []string{"golang", "java", "python", "mysql"},
			}
			err := service.Create(job)
			assert.True(t, repository.SaveWasCalled())
			assert.Equal(t, tc.err, err)
			if err == nil {
				assert.True(t, notification.EnqueueWasCalled())
			}
		})
	}
}

func TestJobService_Match(t *testing.T) {
	testCases := []struct {
		name    string
		pattern *domain.Filter
		err     error
		matches int
	}{
		{"all jobs without filter", nil, nil, 3},
		{"failed by repository error", nil, errors.New("mocked error"), 0},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repository := repositories.MockJobRepository{
				Error: tc.err,
				Jobs: []domain.Job{
					{"Looking for a Technical Leader", "Very long description.", "Ariel Labs", "Argentina", 6000, 8000, domain.FullTime, true, []string{"golang", "java", "python", "mysql"}},
					{"Sr Java developer", "We need you.", "IBM", "USA", 0, 8000, domain.PartTime, false, []string{"java"}},
					{"Junior Python developer", "We need more.", "Globant", "", 0, 8000, domain.Contractor, true, []string{"sql"}},
				},
			}
			service := NewJobService(&repository, nil)
			jobs, err := service.Match(tc.pattern)
			assert.True(t, repository.FilterWasCalled())
			assert.Equal(t, tc.err, err)
			if err == nil {
				assert.NotNil(t, jobs)
				assert.Equal(t, tc.matches, len(jobs))
			}
		})
	}
}