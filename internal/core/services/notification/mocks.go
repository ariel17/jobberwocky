package notification

import "github.com/ariel17/jobberwocky/internal/core/domain"

type MockNotificationService struct {
	enqueueWasCalled bool
}

func (m *MockNotificationService) Enqueue(_ domain.Job) {
	m.enqueueWasCalled = true
}

func (m *MockNotificationService) EnqueueWasCalled() bool {
	return m.enqueueWasCalled
}

func (m *MockNotificationService) StartWorkers() {}

func (m *MockNotificationService) StopWorkers() {}

func (m *MockNotificationService) Process() {}