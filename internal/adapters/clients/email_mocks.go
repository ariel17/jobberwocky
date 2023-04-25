package clients

import (
	"log"
	"sync"
)

type MockEmailProviderClient struct {
	Error         error
	sendWasCalled bool
	sendCalls     int
	mutex         sync.Mutex
}

func (m *MockEmailProviderClient) Send(from, to, subject, body string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	log.Printf("Sending mock email: from=%s, to=%s, subject=%s, body=%s", from, to, subject, body)
	m.sendWasCalled = true
	m.sendCalls += 1
	return m.Error
}

func (m *MockEmailProviderClient) SendWasCalled() bool {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	return m.sendWasCalled
}

func (m *MockEmailProviderClient) SendCalls() int {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	return m.sendCalls
}