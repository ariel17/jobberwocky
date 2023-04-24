package repositories

import (
	"strings"

	"github.com/ariel17/jobberwocky/internal/core/domain"
)

type MockJobRepository struct {
	Jobs  []domain.Job
	Error error
}

func (m *MockJobRepository) Filter(pattern *domain.Filter) ([]domain.Job, error) {
	if m.Error != nil {
		return nil, m.Error
	}
	if pattern == nil {
		return m.Jobs, m.Error
	}
	results := []domain.Job{}
	for _, job := range m.Jobs {
		if (pattern.Text != "" && strings.Contains(job.Title+job.Description, pattern.Text)) ||
			(pattern.Location != "" && pattern.Location == job.Location) ||
			(pattern.Salary > 0 && (
				(job.SalaryMin != nil && pattern.Salary >= *job.SalaryMin && pattern.Salary <= job.SalaryMax) ||
					(job.SalaryMin == nil && pattern.Salary == job.SalaryMax))) ||
			(pattern.Type != "" && pattern.Type == job.Type) ||
			(pattern.IsRemoteFriendly != nil && *pattern.IsRemoteFriendly == job.IsRemoteFriendly) {
			results = append(results, job)
			continue
		}
		for _, patternKeyword := range pattern.Keywords {
			found := false
			for _, jobKeyword := range job.Keywords {
				if patternKeyword == jobKeyword {
					found = true
					break
				}
			}
			if found {
				results = append(results, job)
				break
			}
		}
	}
	return results, m.Error
}

func (m *MockJobRepository) Save(_ domain.Job) error {
	return m.Error
}