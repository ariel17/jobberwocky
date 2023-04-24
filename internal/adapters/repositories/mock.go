package repositories

import "github.com/ariel17/jobberwocky/internal/core/domain"

type MockRepository struct {
	Jobs  []domain.Job
	Error error
}

func (m *MockRepository) Filter(pattern domain.Filter) ([]domain.Job, error) {
	return m.Jobs, m.Error
}

func (m *MockRepository) Save(_ domain.Job) error {
	return m.Error
}