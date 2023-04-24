package domain

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJob_IsTitleValid(t *testing.T) {
	testCases := []struct {
		name    string
		title   string
		isValid bool
	}{
		{"Non-empty title", "This is a title", true},
		{"Empty title", "", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			job := Job{Title: tc.title}
			assert.Equal(t, tc.isValid, job.IsTitleValid())
		})
	}
}

func TestJob_IsTypeValid(t *testing.T) {
	testCases := []struct {
		jobType string
		isValid bool
	}{
		{Contractor, true},
		{FullTime, true},
		{PartTime, true},
		{"other value", false},
	}

	for _, tc := range testCases {
		t.Run(tc.jobType, func(t *testing.T) {
			job := Job{Type: tc.jobType}
			assert.Equal(t, tc.isValid, job.IsTypeValid())
		})
	}
}

func TestJob_IsSalaryValid(t *testing.T) {
	testCases := []struct {
		name      string
		salaryMin int
		salaryMax int
		isValid   bool
	}{
		{"Valid salary range", 100, 200, true},
		{"Fixed salary", 0, 100, true},
		{"Invalid min value", 300, 200, false},
		{"Invalid max value", 100, 0, false},
		{"Invalid max and min value", 0, 0, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			job := Job{SalaryMin: tc.salaryMin, SalaryMax: tc.salaryMax}
			assert.Equal(t, tc.isValid, job.IsSalaryValid())
		})
	}
}

func TestJob_IsLocationAndRemoteFriendlyValid(t *testing.T) {
	testCases := []struct {
		name             string
		location         string
		isRemoteFriendly bool
		isValid          bool
	}{
		{"With location and remote friendly", "Argentina", true, true},
		{"With location and not remote friendly", "Argentina", false, true},
		{"Without location and remote friendly", "", true, true},
		{"Without location and not remote friendly", "", false, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			job := Job{Location: tc.location, IsRemoteFriendly: tc.isRemoteFriendly}
			assert.Equal(t, tc.isValid, job.IsLocationAndIsRemoteFriendlyValid())
		})
	}
}

func TestNewJob(t *testing.T) {
	testCases := []struct {
		name             string
		title            string
		description      string
		company          string
		location         string
		salaryMin        int
		salaryMax        int
		jobType          string
		isRemoteFriendly bool
		keywords         []string
		err              error
	}{
		{"expected values", "title", "description", "company", "location", 10, 20, Contractor, true, []string{"k1", "k2"}, nil},
		{"invalid title", "", "", "", "", 0, 0, "", true, nil, errors.New("title cannot be empty")},
		{"invalid type", "title", "", "", "", 10, 20, "other", true, nil, errors.New("type value is invalid: other")},
		{"invalid salary range", "title", "", "", "", 20, 10, Contractor, true, nil, errors.New("fixed/ranged salary is invalid: min=20, max=10")},
		{"invalid fixed salary", "title", "", "", "", 0, 0, FullTime, true, nil, errors.New("fixed/ranged salary is invalid: min=0, max=0")},
		{"invalid location and no remote friendly", "title", "", "", "", 10, 20, Contractor, false, nil, errors.New("location and remote-friendly values are incorrect: location=, remote friendly=false")},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := NewJob(tc.title, tc.description, tc.company, tc.location, tc.salaryMin, tc.salaryMax, tc.jobType, tc.isRemoteFriendly, tc.keywords...)
			assert.Equal(t, tc.err, err)
		})
	}
}