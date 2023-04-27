package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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

func TestJob_Validate(t *testing.T) {
	testCases := []struct {
		name    string
		job     Job
		isLocal bool
		isValid bool
	}{
		{"valid for external", Job{Title: "This is the title", SalaryMax: 100, Source: "external"}, false, true},
		{"valid for local", Job{Title: "This is the title", Company: "Company", Location: "Location", SalaryMin: 0, SalaryMax: 100, Type: FullTime}, true, true},
		{"invalid by empty title", Job{Company: "Company", Location: "Location", SalaryMin: 0, SalaryMax: 100, Type: FullTime}, false, false},
		{"invalid by type", Job{Title: "This is the title", Company: "Company", Location: "Location", SalaryMin: 0, SalaryMax: 100, Type: "invalid"}, true, false},
		{"invalid by ranged salary", Job{Title: "This is the title", Company: "Company", Location: "Location", SalaryMin: 500, SalaryMax: 100, Type: FullTime}, true, false},
		{"invalid by fixed salary", Job{Title: "This is the title", Company: "Company", Location: "Location", SalaryMin: 0, SalaryMax: 0, Type: FullTime}, true, false},
		{"invalid by location and is_remote_friendly", Job{Title: "This is the title", Company: "Company", SalaryMin: 0, SalaryMax: 100, IsRemoteFriendly: boolPointer(true), Type: FullTime}, false, false},
		{"invalid by repeated keywords", Job{Title: "This is the title", Company: "Company", SalaryMin: 0, SalaryMax: 100, IsRemoteFriendly: boolPointer(true), Type: FullTime, Keywords: []string{"k1", "k1"}}, true, false},
		{"invalid by source for external", Job{Title: "This is the title", SalaryMax: 100}, false, false},
		{"invalid by source for local", Job{Title: "This is the title", Company: "Company", SalaryMin: 0, SalaryMax: 100, IsRemoteFriendly: boolPointer(true), Type: FullTime, Keywords: []string{"k1", "k2"}, Source: "external"}, true, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.isValid, tc.job.Validate(tc.isLocal) == nil)
		})
	}
}

func TestJob_IsSourceValid(t *testing.T) {
	testCases := []struct {
		name    string
		source  string
		isLocal bool
		isValid bool
	}{
		{"empty source for local", "", true, true},
		{"not empty source for external", "external", false, true},
		{"invalid source for local", "external", true, false},
		{"invalid source for external", "", false, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			job := Job{Title: "A title", SalaryMax: 100, Source: tc.source}
			assert.Equal(t, tc.isValid, job.isSourceValid(tc.isLocal))
		})
	}
}

func TestPattern_IsEmpty(t *testing.T) {
	testCases := []struct {
		name    string
		pattern Pattern
		isEmpty bool
	}{
		{"empty", Pattern{}, true},
		{"not empty", Pattern{Text: "text"}, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.isEmpty, tc.pattern.IsEmpty())
		})
	}
}

func TestPattern_Validate(t *testing.T) {
	testCases := []struct {
		name    string
		pattern Pattern
		isValid bool
	}{
		{"valid by empty", Pattern{}, true},
		{"valid by values", Pattern{Text: "text", Company: "Company", Location: "Location", Salary: 100, Type: FullTime, Keywords: []string{"k1", "k2"}}, true},
		{"invalid by repeated keywords", Pattern{Text: "text", Company: "Company", Location: "Location", Salary: 100, Type: FullTime, Keywords: []string{"k1", "k1"}}, false},
		{"invalid by invalid type", Pattern{Text: "text", Company: "Company", Location: "Location", Salary: 100, Type: "invalid", Keywords: []string{"k1", "k2"}}, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.isValid, tc.pattern.Validate() == nil)
		})
	}
}

func TestSubscription_Validate(t *testing.T) {
	testCases := []struct {
		name         string
		subscription Subscription
		isValid      bool
	}{
		{"valid with empty pattern", Subscription{Pattern: Pattern{}, Email: "person@example.com"}, true},
		{"valid with patterns", Subscription{Pattern: Pattern{Text: "text", Company: "Company", Location: "Location", Salary: 100, Type: FullTime, Keywords: []string{"k1", "k2"}}, Email: "person@example.com"}, true},
		{"invalid by empty email", Subscription{Pattern: Pattern{Text: "text", Company: "Company", Location: "Location", Salary: 100, Type: FullTime, Keywords: []string{"k1", "k2"}}}, false},
		{"invalid by invalid pattern", Subscription{Pattern: Pattern{Text: "text", Company: "Company", Location: "Location", Salary: 100, Type: "invalid", Keywords: []string{"k1", "k2"}}}, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.isValid, tc.subscription.Validate() == nil)
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

func boolPointer(v bool) *bool {
	var newValue = v
	return &newValue
}