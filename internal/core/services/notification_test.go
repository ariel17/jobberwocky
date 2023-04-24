package services

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/ariel17/jobberwocky/internal/adapters/clients"
	"github.com/ariel17/jobberwocky/internal/adapters/repositories"
	"github.com/ariel17/jobberwocky/internal/core/domain"
)

func TestNotificationService_Enqueue(t *testing.T) {
	testCases := []struct {
		name          string
		job           domain.Job
		matches       int
		repositoryErr error
		emailErr      error
	}{
		{"matches and sends notification", domain.Job{"Title", "Description", "Company", "Argentina", 60, 80, domain.FullTime, true, []string{"k1", "k2", "k3"}}, 2, nil, nil},
		{"failed by repository error", domain.Job{}, 0, errors.New("mocked repository error"), nil},
		{"failed by email client error", domain.Job{"Another", "Description", "Company", "USA", 60, 80, domain.FullTime, true, []string{"k1", "k2", "k3"}}, 1, nil, errors.New("mocked email error")},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			emailClient := clients.MockEmailProviderClient{Error: tc.emailErr}
			repository := repositories.MockSubscriptionRepository{
				MockRepository: repositories.MockRepository{Error: tc.repositoryErr},
				Subscriptions: []domain.Subscription{
					{domain.Filter{Text: "Title"}, "person1@example.com"},
					{domain.Filter{Text: "a different thing"}, "person2@example.com"},
					{domain.Filter{Location: "Argentina"}, "person3@example.com"},
					{domain.Filter{Location: "USA"}, "person3@example.com"},
				},
			}
			service := NewNotificationService(10, &repository, &emailClient)
			service.StartWorkers()
			defer service.StopWorkers()

			service.Enqueue(tc.job)
			time.Sleep(10 * time.Millisecond)

			assert.True(t, repository.FilterWasCalled())
			if tc.repositoryErr == nil {
				assert.True(t, emailClient.SendWasCalled())
				assert.Equal(t, tc.matches, emailClient.SendCalls())
			}
		})
	}
}