package notification

import (
	"errors"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/ariel17/jobberwocky/internal/adapters/clients"
	"github.com/ariel17/jobberwocky/internal/adapters/repositories"
	"github.com/ariel17/jobberwocky/internal/configs"
	"github.com/ariel17/jobberwocky/internal/core/domain"
	"github.com/ariel17/jobberwocky/internal/internal_test"
)

func TestCreateBody(t *testing.T) {
	job := domain.Job{"Title", "Description", "Company", "Argentina", 60, 80, domain.FullTime, internal_test.BoolPointer(true), []string{"k1", "k2", "k3"}, ""}
	testCases := []struct {
		name    string
		file    string
		success bool
	}{
		{"ok", configs.DefaultTemplate, true},
		{"template not exists", "notexists.tmpl", false},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			body, err := createBody(tc.file, job)
			assert.Equal(t, tc.success, err == nil)
			if tc.success {
				values := []string{
					job.Title, job.Description, job.Company, job.Location,
					strconv.Itoa(job.SalaryMin), strconv.Itoa(job.SalaryMax),
					job.Type, fmt.Sprintf("%v", *job.IsRemoteFriendly), fmt.Sprintf("%s", job.Keywords)}
				for _, v := range values {
					assert.Contains(t, body, v)
				}
			}
		})
	}
}

func TestNotificationService_Enqueue(t *testing.T) {
	testCases := []struct {
		name          string
		job           domain.Job
		matches       int
		repositoryErr error
		emailErr      error
	}{
		{"matches and sends notification", domain.Job{"Title", "Description", "Company", "Argentina", 60, 80, domain.FullTime, internal_test.BoolPointer(true), []string{"k1", "k2", "k3"}, ""}, 2, nil, nil},
		{"failed by repository error", domain.Job{}, 0, errors.New("mocked repository error"), nil},
		{"failed by email client error", domain.Job{"Another", "Description", "Company", "USA", 60, 80, domain.FullTime, internal_test.BoolPointer(true), []string{"k1", "k2", "k3"}, ""}, 1, nil, errors.New("mocked email error")},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			emailClient := clients.MockEmailProviderClient{Error: tc.emailErr}
			repository := repositories.MockSubscriptionRepository{
				MockFilter: internal_test.MockFilter{Error: tc.repositoryErr},
				Subscriptions: []domain.Subscription{
					{domain.Pattern{Text: "Title"}, "person1@example.com"},
					{domain.Pattern{Text: "a different thing"}, "person2@example.com"},
					{domain.Pattern{Location: "Argentina"}, "person3@example.com"},
					{domain.Pattern{Location: "USA"}, "person3@example.com"},
				},
			}
			service := NewNotificationService(configs.GetNotificationWorkers(), &repository, &emailClient)
			service.StartWorkers()
			defer service.StopWorkers()

			service.Enqueue(tc.job)
			time.Sleep(10 * time.Millisecond)

			assert.True(t, repository.FilterWasCalled)
			if tc.repositoryErr == nil {
				assert.True(t, emailClient.SendWasCalled())
				assert.Equal(t, tc.matches, emailClient.SendCalls())
			}
		})
	}
}