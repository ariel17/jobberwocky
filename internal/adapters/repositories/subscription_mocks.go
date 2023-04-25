package repositories

import (
	"github.com/ariel17/jobberwocky/internal/core/domain"
	"github.com/ariel17/jobberwocky/internal/internal_test"
)

type MockSubscriptionRepository struct {
	internal_test.MockFilter
	internal_test.MockSave
	Subscriptions []domain.Subscription
}

func (m *MockSubscriptionRepository) Filter(job domain.Job) ([]domain.Subscription, error) {
	m.FilterWasCalled = true
	results := []domain.Subscription{}
	for _, subscription := range m.Subscriptions {
		if internal_test.Matches(subscription.Pattern, job) {
			results = append(results, subscription)
		}
	}
	return results, m.Error
}

func (m *MockSubscriptionRepository) Save(_ domain.Subscription) error {
	m.SaveWasCalled = true
	return m.Error
}