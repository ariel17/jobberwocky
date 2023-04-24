package services

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ariel17/jobberwocky/internal/adapters/repositories"
	"github.com/ariel17/jobberwocky/internal/core/domain"
)

func TestJobService_Create(t *testing.T) {
	testCases := []struct {
		name    string
		err     error
		success bool
	}{
		{"ok", nil, true},
		{"failed", errors.New("mocked error"), false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repository := repositories.MockRepository{
				Error: tc.err,
			}
			service := NewJobService(&repository)

			job := domain.Job{
				Title:            "Looking for a Technical Leader",
				Description:      "This is the longest part where we describe all the details about the job and required skills.",
				Company:          "Ariel Labs",
				Location:         "Argentina",
				SalaryMin:        6000,
				SalaryMax:        8000,
				Type:             "Full-Time",
				IsRemoteFriendly: true,
				Keywords:         []string{"golang", "java", "python", "mysql"},
			}
			err := service.Create(job)
			assert.Equal(t, tc.success, err == nil)
		})
	}
}