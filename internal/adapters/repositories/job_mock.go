package repositories

import (
	"strings"

	"github.com/ariel17/jobberwocky/internal/core/domain"
)

// MockFilter is a generic implementation for mocks that need to simulate
// matching behavior from its job source with a given pattern.
type MockFilter struct {
	Error           error
	Jobs            []domain.Job
	filterWasCalled bool
}

func (m *MockFilter) Filter(pattern *domain.Pattern) ([]domain.Job, error) {
	m.filterWasCalled = true
	if m.Error != nil {
		return nil, m.Error
	}
	if pattern == nil {
		return m.Jobs, m.Error
	}
	results := []domain.Job{}
	for _, job := range m.Jobs {
		if matches(*pattern, job) {
			results = append(results, job)
		}

	}
	return results, m.Error
}

func (m *MockFilter) FilterWasCalled() bool {
	return m.filterWasCalled
}

func allKeywordsContained(patternKeywords, jobKeywords []string) bool {
	for _, pk := range patternKeywords {
		found := false
		for _, jk := range jobKeywords {
			if strings.ToLower(pk) == strings.ToLower(jk) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func matches(pattern domain.Pattern, job domain.Job) bool {
	return (pattern.Company == "" || pattern.Company == job.Company) &&
		(pattern.Text == "" || strings.Contains(strings.ToLower(job.Title+job.Description), strings.ToLower(pattern.Text))) &&
		(pattern.Location == "" || pattern.Location == job.Location) &&
		(pattern.Salary == 0 || (
			(job.SalaryMin > 0 && pattern.Salary >= job.SalaryMin && pattern.Salary <= job.SalaryMax) ||
				(job.SalaryMin == 0 && pattern.Salary == job.SalaryMax))) &&
		(pattern.Type == "" || pattern.Type == job.Type) &&
		(pattern.IsRemoteFriendly == nil || (pattern.IsRemoteFriendly != nil && *pattern.IsRemoteFriendly == job.IsRemoteFriendly)) &&
		allKeywordsContained(pattern.Keywords, job.Keywords)
}

// MockSave provides common fields/methods to be used on a Save() method and
// check its usage.
type MockSave struct {
	saveWasCalled bool
}

func (m *MockSave) SaveWasCalled() bool {
	return m.saveWasCalled
}

// MockJobRepository uses generic mock implementation to comply with a
// repository behavior.
type MockJobRepository struct {
	MockFilter
	MockSave
}

func (m *MockJobRepository) Save(_ domain.Job) error {
	m.saveWasCalled = true
	return m.Error
}