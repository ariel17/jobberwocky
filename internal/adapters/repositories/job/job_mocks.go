package job

import (
	"github.com/ariel17/jobberwocky/internal/core/domain"
	"github.com/ariel17/jobberwocky/internal/internal_test"
)

// MockJobRepository uses generic mock implementation to comply with a
// repository behavior.
type MockJobRepository struct {
	internal_test.MockFilter
	internal_test.MockSave
}

func (m *MockJobRepository) Save(_ domain.Job) error {
	m.SetSaveCalled()
	return m.Error
}

func (m *MockJobRepository) SyncSchemas() error {
	return nil
}