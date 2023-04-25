package job

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ariel17/jobberwocky/internal/adapters/clients"
	"github.com/ariel17/jobberwocky/internal/adapters/repositories"
	"github.com/ariel17/jobberwocky/internal/core/domain"
	"github.com/ariel17/jobberwocky/internal/core/services/notification"
	"github.com/ariel17/jobberwocky/internal/internal_test"
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
				MockFilter: internal_test.MockFilter{Error: tc.err},
			}
			ns := notification.MockNotificationService{}
			service := NewJobService(&repository, &ns)
			job := domain.Job{"Title", "Description", "Company", "Argentina", 60, 80, domain.FullTime, internal_test.BoolPointer(true), []string{"k1", "k2", "k3"}, ""}
			err := service.Create(job)
			assert.True(t, repository.SaveWasCalled)
			assert.Equal(t, tc.err, err)
			if err == nil {
				assert.True(t, ns.EnqueueWasCalled())
			}
		})
	}
}

func TestJobService_Match(t *testing.T) {
	testCases := []struct {
		name          string
		pattern       *domain.Pattern
		repositoryErr error
		externalErr   error
		matches       int
	}{
		{"all without filter", nil, nil, nil, 4},
		{"all with empty filter", &domain.Pattern{}, nil, nil, 4},
		{"pattern by title text", &domain.Pattern{Text: "technical"}, nil, nil, 1},
		{"pattern by description text", &domain.Pattern{Text: "you"}, nil, nil, 1},
		{"pattern by salary in range", &domain.Pattern{Salary: 7000}, nil, nil, 1},
		{"pattern by ranged/fixed salary", &domain.Pattern{Salary: 8000}, nil, nil, 4},
		{"pattern by company", &domain.Pattern{Company: "IBM"}, nil, nil, 1},
		{"pattern by location", &domain.Pattern{Location: "USA"}, nil, nil, 1},
		{"pattern by type", &domain.Pattern{Type: domain.Contractor}, nil, nil, 1},
		{"pattern by remote friendly", &domain.Pattern{IsRemoteFriendly: internal_test.BoolPointer(true)}, nil, nil, 2},
		{"pattern by single keyword", &domain.Pattern{Keywords: []string{"sql"}}, nil, nil, 2},
		{"pattern by multiple keywords that does not match", &domain.Pattern{Keywords: []string{"sql", "java"}}, nil, nil, 0},
		{"pattern by multiple keywords that matches", &domain.Pattern{Keywords: []string{"golang", "java"}}, nil, nil, 1},
		{"pattern by keywords and text", &domain.Pattern{Text: "technical", Keywords: []string{"sql", "java"}}, nil, nil, 0},
		{"pattern by text matches external", &domain.Pattern{Text: "external"}, nil, nil, 1},
		{"failed by repository error", nil, errors.New("mocked repository error"), nil, 0},
		{"failed by external error", nil, nil, errors.New("mocked external error"), 0},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repository := repositories.MockJobRepository{
				MockFilter: internal_test.MockFilter{
					Error: tc.repositoryErr,
					Jobs: []domain.Job{
						{"Looking for a Technical Leader", "Very long description.", "Ariel Labs", "Argentina", 6000, 8000, domain.FullTime, internal_test.BoolPointer(true), []string{"golang", "java", "python", "mysql"}, ""},
						{"Sr Java developer", "We need you.", "IBM", "USA", 0, 8000, domain.PartTime, internal_test.BoolPointer(false), []string{"java"}, ""},
						{"Junior Python developer", "We need more.", "Globant", "", 0, 8000, domain.Contractor, internal_test.BoolPointer(true), []string{"sql"}, ""},
					},
				},
			}
			external := clients.MockExternalJobClient{
				MockFilter: internal_test.MockFilter{
					Error: tc.externalErr,
					Jobs: []domain.Job{
						{"External", "", "", "Argentina", 0, 8000, "", nil, []string{"sql"}, "external"},
					},
				},
			}
			service := NewJobService(&repository, nil, &external)
			jobs, err := service.Filter(tc.pattern)
			assert.True(t, repository.FilterWasCalled)
			if tc.repositoryErr != nil {
				assert.Equal(t, tc.repositoryErr, err)
			}
			if tc.externalErr != nil {
				assert.Equal(t, tc.externalErr, err)
			}
			if err == nil {
				assert.NotNil(t, jobs)
				assert.Equal(t, tc.matches, len(jobs))
			}
		})
	}
}