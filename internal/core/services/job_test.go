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
			job := domain.Job{"Title", "Description", "Company", "Argentina", 60, 80, domain.FullTime, true, []string{"k1", "k2", "k3"}}
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
		{"all without filter", nil, nil, 3},
		{"Filter by title text", &domain.Filter{Text: "technical"}, nil, 1},
		{"Filter by description text", &domain.Filter{Text: "you"}, nil, 1},
		{"Filter by salary in range", &domain.Filter{Salary: 7000}, nil, 1},
		{"Filter by ranged/fixed salary", &domain.Filter{Salary: 8000}, nil, 3},
		{"Filter by company", &domain.Filter{Company: "IBM"}, nil, 1},
		{"Filter by location", &domain.Filter{Location: "USA"}, nil, 1},
		{"Filter by type", &domain.Filter{Type: domain.Contractor}, nil, 1},
		{"Filter by remote friendly", &domain.Filter{IsRemoteFriendly: boolPointer(true)}, nil, 2},
		{"Filter by single keyword", &domain.Filter{Keywords: []string{"sql"}}, nil, 1},
		{"Filter by multiple keywords that does not match", &domain.Filter{Keywords: []string{"sql", "java"}}, nil, 0},
		{"Filter by multiple keywords that matches", &domain.Filter{Keywords: []string{"golang", "java"}}, nil, 1},
		{"Filter by keywords and text", &domain.Filter{Text: "technical", Keywords: []string{"sql", "java"}}, nil, 0},
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

func boolPointer(v bool) *bool {
	newValue := v
	return &newValue
}