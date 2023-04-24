package repositories

import "github.com/ariel17/jobberwocky/internal/core/domain"

type MockSubscriptionRepository struct {
	MockRepository
	Subscriptions []domain.Subscription
}

func (m *MockSubscriptionRepository) Filter(job domain.Job) ([]domain.Subscription, error) {
	m.filterWasCalled = true
	results := []domain.Subscription{}
	for _, subscription := range m.Subscriptions {
		if matches(subscription.Filter, job) {
			results = append(results, subscription)
		}
	}
	return results, m.Error
}

func (m *MockSubscriptionRepository) Save(_ domain.Subscription) error {
	m.saveWasCalled = true
	return m.Error
}