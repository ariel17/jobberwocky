package job

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ariel17/jobberwocky/internal/adapters/repositories"
	"github.com/ariel17/jobberwocky/internal/core/domain"
	"github.com/ariel17/jobberwocky/internal/core/services/notification"
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
				MockFilter: repositories.MockFilter{Error: tc.err},
			}
			notification := notification.MockNotificationService{}
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
		pattern *domain.Pattern
		err     error
		matches int
	}{
		{"all without filter", nil, nil, 3},
		{"all with empty filter", &domain.Pattern{}, nil, 3},
		{"Pattern by title text", &domain.Pattern{Text: "technical"}, nil, 1},
		{"Pattern by description text", &domain.Pattern{Text: "you"}, nil, 1},
		{"Pattern by salary in range", &domain.Pattern{Salary: 7000}, nil, 1},
		{"Pattern by ranged/fixed salary", &domain.Pattern{Salary: 8000}, nil, 3},
		{"Pattern by company", &domain.Pattern{Company: "IBM"}, nil, 1},
		{"Pattern by location", &domain.Pattern{Location: "USA"}, nil, 1},
		{"Pattern by type", &domain.Pattern{Type: domain.Contractor}, nil, 1},
		{"Pattern by remote friendly", &domain.Pattern{IsRemoteFriendly: boolPointer(true)}, nil, 2},
		{"Pattern by single keyword", &domain.Pattern{Keywords: []string{"sql"}}, nil, 1},
		{"Pattern by multiple keywords that does not match", &domain.Pattern{Keywords: []string{"sql", "java"}}, nil, 0},
		{"Pattern by multiple keywords that matches", &domain.Pattern{Keywords: []string{"golang", "java"}}, nil, 1},
		{"Pattern by keywords and text", &domain.Pattern{Text: "technical", Keywords: []string{"sql", "java"}}, nil, 0},
		{"failed by repository error", nil, errors.New("mocked error"), 0},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repository := repositories.MockJobRepository{
				MockFilter: repositories.MockFilter{
					Error: tc.err,
					Jobs: []domain.Job{
						{"Looking for a Technical Leader", "Very long description.", "Ariel Labs", "Argentina", 6000, 8000, domain.FullTime, true, []string{"golang", "java", "python", "mysql"}},
						{"Sr Java developer", "We need you.", "IBM", "USA", 0, 8000, domain.PartTime, false, []string{"java"}},
						{"Junior Python developer", "We need more.", "Globant", "", 0, 8000, domain.Contractor, true, []string{"sql"}},
					},
				},
			}
			service := NewJobService(&repository, nil)
			jobs, err := service.Filter(tc.pattern)
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