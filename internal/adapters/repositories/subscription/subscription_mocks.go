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
	results := make([]domain.Subscription, 0)
	for _, subscription := range m.Subscriptions {
		if mocks.Matches(subscription.Pattern, job) {
			results = append(results, subscription)
		}
	}
	return results, m.Error
}

func (m *MockSubscriptionRepository) Save(subscription domain.Subscription) error {
	m.SetSaveCalled()
	m.Subscriptions = append(m.Subscriptions, subscription)
	return m.Error
}

func (m *MockSubscriptionRepository) SyncSchemas() error {
	return nil
}