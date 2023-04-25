package clients

import "log"

type MockEmailProviderClient struct {
	Error         error
	sendWasCalled bool
	sendCalls     int
}

func (m *MockEmailProviderClient) Send(from, to, subject, body string) error {
	log.Printf("Sending mock email: from=%s, to=%s, subject=%s, body=%s", from, to, subject, body)
	m.sendWasCalled = true
	m.sendCalls += 1
	return m.Error
}

func (m *MockEmailProviderClient) SendWasCalled() bool {
	return m.sendWasCalled
}

func (m *MockEmailProviderClient) SendCalls() int {
	return m.sendCalls
}