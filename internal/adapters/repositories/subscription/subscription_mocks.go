package subscription

import (
	"github.com/ariel17/jobberwocky/internal/core/domain"
	mocks "github.com/ariel17/jobberwocky/internal/internal_test"
)

type MockSubscriptionRepository struct {
	mocks.MockFilter
	mocks.MockSave
	Subscriptions []domain.Subscription
}

func (m *MockSubscriptionRepository) Filter(job domain.Job) ([]domain.Subscription, error) {
	m.SetFilterCalled()
	results := []domain.Subscription{}
	for _, subscription := range m.Subscriptions {
		if mocks.Matches(subscription.Pattern, job) {
			results = append(results, subscription)
		}
	}
	return results, m.Error
}

func (m *MockSubscriptionRepository) Save(_ domain.Subscription) error {
	m.SetSaveCalled()
	return m.Error
}

func (m *MockSubscriptionRepository) SyncSchemas() error {
	return nil
}