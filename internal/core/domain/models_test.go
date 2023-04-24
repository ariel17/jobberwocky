package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
		{"Without location and not remote friendly", "", true, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			job := Job{Location: tc.location, IsRemoteFriendly: tc.isRemoteFriendly}
			assert.Equal(t, tc.isValid, job.IsLocationAndIsRemoteFriendlyValid())
		})
	}
}