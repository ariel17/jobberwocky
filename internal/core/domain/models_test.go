package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJob_IsTypeValid(t *testing.T) {
	testCases := []struct {
		name    string
		jobType string
		isValid bool
	}{
		{"Contractor", Contractor, true},
		{"Full-Time", Contractor, true},
		{"Part-Time", Contractor, true},
		{"other value", "other value", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			job := Job{Type: tc.jobType}
			assert.Equal(t, tc.isValid, job.IsTypeValid())
		})
	}
}