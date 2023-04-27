package subscription

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/ariel17/jobberwocky/internal/adapters/repositories/subscription"
	"github.com/ariel17/jobberwocky/internal/core/domain"
	subscriptionService "github.com/ariel17/jobberwocky/internal/core/services/subscription"
	helpers "github.com/ariel17/jobberwocky/internal/internal_test"
)

func TestSubscriptionHTTPHandler_Post(t *testing.T) {
	testCases := []struct {
		name                      string
		goldenPathBody            string
		goldenPathResponse        string
		existingSubscriptions     []domain.Subscription
		subscriptionRepositoryErr error
		statusCode                int
	}{
		{"new subscription without filters", "post_body_no_filters", "", nil, nil, http.StatusCreated},
		{"new subscription with filters", "post_body_with_filters", "", nil, nil, http.StatusCreated},
		{"failed by empty email", "post_body_invalid_email", "post_response_invalid_email", nil, nil, http.StatusBadRequest},
		{"failed by invalid type", "post_body_invalid_type", "post_response_invalid_type", nil, nil, http.StatusBadRequest},
		{"failed by repeated keywords", "post_body_invalid_keywords", "post_response_invalid_keywords", nil, nil, http.StatusBadRequest},
		{"failed by existing subscription", "post_body_with_filters", "post_response_already_exists_error", []domain.Subscription{{Email: "person1@example.com", Pattern: domain.Pattern{Company: "SpaceX", Type: "Contractor", Salary: 1000, IsRemoteFriendly: helpers.BoolPointer(true)}}}, nil, http.StatusInternalServerError},
		{"failed by repository error", "post_body_with_filters", "post_response_repository_error", nil, errors.New("mocked repository error"), http.StatusInternalServerError},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := subscription.MockSubscriptionRepository{
				MockFilter:    helpers.MockFilter{Error: tc.subscriptionRepositoryErr},
				MockSave:      helpers.MockSave{},
				Subscriptions: tc.existingSubscriptions,
			}
			s := subscriptionService.NewSubscriptionService(&r)

			handler := NewSubscriptionHTTPHandler(s)
			router := gin.Default()
			handler.ConfigureRoutes(router)

			body := helpers.GetGoldenFile(t, tc.goldenPathBody)
			req, _ := http.NewRequest(http.MethodPost, subscriptionPath, strings.NewReader(body))
			rr := httptest.NewRecorder()

			router.ServeHTTP(rr, req)
			assert.Equal(t, tc.statusCode, rr.Code)

			if tc.statusCode != http.StatusCreated {
				expected := helpers.GetGoldenFile(t, tc.goldenPathResponse)
				assert.Equal(t, expected, rr.Body.String())
			} else {
				assert.True(t, r.SaveWasCalled())

				var s1, s2 domain.Subscription
				assert.Nil(t, json.Unmarshal([]byte(body), &s1))
				s2 = r.Subscriptions[len(r.Subscriptions)-1]
				assert.Equal(t, s1.Email, s2.Email)
				assert.Equal(t, s1.Text, s2.Text)
				assert.Equal(t, s1.Company, s2.Company)
				assert.Equal(t, s1.Location, s2.Location)
				assert.Equal(t, s1.Salary, s2.Salary)
				assert.Equal(t, s1.Type, s2.Type)
				if s1.IsRemoteFriendly == nil {
					assert.Nil(t, s2.IsRemoteFriendly)
				} else {
					assert.Equal(t, *s1.IsRemoteFriendly, *s2.IsRemoteFriendly)
				}
				assert.Equal(t, len(s1.Keywords), len(s2.Keywords))
			}
		})
	}
}