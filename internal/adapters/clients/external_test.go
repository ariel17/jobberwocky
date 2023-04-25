package clients

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ariel17/jobberwocky/internal/core/domain"
)

func TestExternalJobClient_Filter(t *testing.T) {
	testCases := []struct {
		name    string
		pattern domain.Pattern
		success bool
	}{
		{},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			client := NewJobberwockyExteralJobClient()
			jobs, err := client.Filter(tc.pattern)
			assert.Equal(t, tc.success, jobs != nil)
			assert.Equal(t, tc.success, err == nil)
		})
	}
}