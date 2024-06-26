package internal_test

import (
	"strings"
	"sync"

	"github.com/ariel17/jobberwocky/internal/core/domain"
)

// MockFilter is a generic implementation for mocks that need to simulate
// matching behavior from its job source with a given pattern.
type MockFilter struct {
	Error           error
	Jobs            []domain.Job
	filterWasCalled bool
	mutex           sync.Mutex
}

func (m *MockFilter) Filter(pattern *domain.Pattern) ([]domain.Job, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.filterWasCalled = true
	if m.Error != nil {
		return nil, m.Error
	}
	if pattern == nil {
		return m.Jobs, m.Error
	}
	results := make([]domain.Job, 0)
	for _, job := range m.Jobs {
		if Matches(*pattern, job) {
			results = append(results, job)
		}
	}
	return results, m.Error
}

func (m *MockFilter) FilterWasCalled() bool {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	return m.filterWasCalled
}

func (m *MockFilter) SetFilterCalled() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.filterWasCalled = true
}

// Matches implements a basic logic to filter jobs that matches the given
// pattern. Values are inclusive.
func Matches(pattern domain.Pattern, job domain.Job) bool {
	return (pattern.Company == "" || pattern.Company == job.Company) &&
		(pattern.Text == "" || strings.Contains(strings.ToLower(job.Title+job.Description), strings.ToLower(pattern.Text))) &&
		(pattern.Location == "" || pattern.Location == job.Location) &&
		(pattern.Salary == 0 || (
			(job.SalaryMin > 0 && pattern.Salary >= job.SalaryMin && pattern.Salary <= job.SalaryMax) ||
				(job.SalaryMin == 0 && pattern.Salary == job.SalaryMax))) &&
		(pattern.Type == "" || pattern.Type == job.Type) &&
		(pattern.IsRemoteFriendly == nil ||
			(pattern.IsRemoteFriendly != nil && job.IsRemoteFriendly != nil && *pattern.IsRemoteFriendly == *job.IsRemoteFriendly)) &&
		domain.AllKeywordsContained(pattern.Keywords, job.Keywords)
}

// MockSave provides common fields/methods to be used on a Save() method and
// check its usage.
type MockSave struct {
	saveWasCalled bool
	mutex         sync.Mutex
}

func (m *MockSave) SaveWasCalled() bool {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	return m.saveWasCalled
}

func (m *MockSave) SetSaveCalled() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.saveWasCalled = true
}