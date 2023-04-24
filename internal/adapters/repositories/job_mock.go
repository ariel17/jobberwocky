package repositories

import (
	"strings"

	"github.com/ariel17/jobberwocky/internal/core/domain"
)

type MockJobRepository struct {
	Jobs            []domain.Job
	Error           error
	saveWasCalled   bool
	filterWasCalled bool
}

func (m *MockJobRepository) Filter(pattern *domain.Filter) ([]domain.Job, error) {
	m.filterWasCalled = true
	if m.Error != nil {
		return nil, m.Error
	}
	if pattern == nil {
		return m.Jobs, m.Error
	}
	results := []domain.Job{}
	for _, job := range m.Jobs {
		if (pattern.Company == "" || pattern.Company == job.Company) &&
			(pattern.Text == "" || strings.Contains(strings.ToLower(job.Title+job.Description), strings.ToLower(pattern.Text))) &&
			(pattern.Location == "" || pattern.Location == job.Location) &&
			(pattern.Salary == 0 || (
				(job.SalaryMin > 0 && pattern.Salary >= job.SalaryMin && pattern.Salary <= job.SalaryMax) ||
					(job.SalaryMin == 0 && pattern.Salary == job.SalaryMax))) &&
			(pattern.Type == "" || pattern.Type == job.Type) &&
			(pattern.IsRemoteFriendly == nil || (pattern.IsRemoteFriendly != nil && *pattern.IsRemoteFriendly == job.IsRemoteFriendly)) &&
			allKeywordsContained(pattern.Keywords, job.Keywords) {
			results = append(results, job)
			continue
		}

	}
	return results, m.Error
}

func (m *MockJobRepository) Save(_ domain.Job) error {
	m.saveWasCalled = true
	return m.Error
}

func (m *MockJobRepository) SaveWasCalled() bool {
	return m.saveWasCalled
}

func (m *MockJobRepository) FilterWasCalled() bool {
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