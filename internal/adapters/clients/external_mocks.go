package clients

import (
	"io"
	"net/http"
	"strings"

	"github.com/ariel17/jobberwocky/internal/core/domain"
	"github.com/ariel17/jobberwocky/internal/internal_test"
)

type MockExternalJobClient struct {
	internal_test.MockFilter
}

func (m *MockExternalJobClient) Name() string {
	return "mock"
}

func (m *MockExternalJobClient) PatternIsSearchable(pattern *domain.Pattern) bool {
	return true
}

type MockHTTPClient struct {
	StatusCode int
	Body       string
	Error      error
}

func (m *MockHTTPClient) Do(_ *http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: m.StatusCode,
		Body:       io.NopCloser(strings.NewReader(m.Body)),
	}, m.Error
}