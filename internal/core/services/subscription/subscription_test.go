package subscription

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ariel17/jobberwocky/internal/adapters/repositories/subscription"
	"github.com/ariel17/jobberwocky/internal/core/domain"
	helpers "github.com/ariel17/jobberwocky/internal/internal_test"
)

func TestSubscriptionService_Create(t *testing.T) {
	newSubscription := domain.Subscription{
		Pattern: domain.Pattern{Text: "Title"},
		Email:   "person1@example.com",
	}
	testCases := []struct {
		name string
		err  error
	}{
		{"subscription created", nil},
		{"failed by repository error", errors.New("mocked error")},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repository := subscription.MockSubscriptionRepository{
				MockFilter: helpers.MockFilter{Error: tc.err},
				MockSave:   helpers.MockSave{},
			}
			service := NewSubscriptionService(&repository)
			err := service.Create(newSubscription)
			assert.True(t, repository.SaveWasCalled())
			assert.Equal(t, tc.err, err)
		})
	}
}