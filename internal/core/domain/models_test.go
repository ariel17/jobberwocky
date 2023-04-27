package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJob_Validate_Title(t *testing.T) {
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
			job := Job{Title: tc.title, SalaryMax: 100, Source: "external"}
			assert.Equal(t, tc.isValid, job.Validate(false) == nil)
		})
	}
}

func TestIsTypeValid(t *testing.T) {
	testCases := []struct {
		name    string
		jobType string
		source  string
		isValid bool
	}{
		{"local source and contractor", Contractor, "", true},
		{"local source and full time", FullTime, "", true},
		{"local source and part time", PartTime, "", true},
		{"local source and invalid value", "other value", "", false},
		{"external source and no type", "", "external", true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.isValid, isTypeValid(tc.source, tc.jobType))
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
			assert.Equal(t, tc.isValid, job.isSalaryValid())
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
		{"Without location and remote friendly", "", true, false},
		{"Without location and not remote friendly", "", false, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			job := Job{Location: tc.location, IsRemoteFriendly: &tc.isRemoteFriendly}
			assert.Equal(t, tc.isValid, job.isLocationAndIsRemoteFriendlyValid())
		})
	}
}

func TestIsKeywordsValid(t *testing.T) {
	testCases := []struct {
		name     string
		keywords []string
		isValid  bool
	}{
		{"unique keywords", []string{"k1", "k2"}, true},
		{"empty keywords", []string{}, true},
		{"repeated keywords", []string{"k1", "k1"}, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.isValid, isKeywordsValid(tc.keywords))
		})
	}
}