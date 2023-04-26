package subscription

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/ariel17/jobberwocky/internal/adapters/repositories/keyword"
	"github.com/ariel17/jobberwocky/internal/core/domain"
	helpers "github.com/ariel17/jobberwocky/internal/internal_test"
)

const (
	dbName = "subscription.db"
)

func TestSubscriptionRepository_Save(t *testing.T) {
	testCases := []struct {
		name         string
		subscription domain.Subscription
	}{
		{"success", domain.Subscription{Pattern: domain.Pattern{"Title", "Company", "Argentina", 60, domain.FullTime, helpers.BoolPointer(true), []string{"k1", "k2", "k3"}}, Email: "arielgerardorios@gmail.com"}},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_ = os.Remove(dbName)
			db, err := gorm.Open(sqlite.Open(dbName), &gorm.Config{
				Logger: logger.Default.LogMode(logger.Info),
			})
			assert.Nil(t, err)

			keywordRepository := keyword.NewKeywordRepository(db)
			assert.Nil(t, keywordRepository.SyncSchemas())

			repository := NewSubscriptionRepository(db)
			assert.Nil(t, repository.SyncSchemas())

			err = repository.Save(tc.subscription)
			assert.Nil(t, err)

			var subscription Subscription
			db.First(&subscription, 1)
			assert.Equal(t, subscription.Text, tc.subscription.Text)
			assert.Equal(t, subscription.Company, tc.subscription.Company)
			assert.Equal(t, subscription.Location, tc.subscription.Location)
			assert.Equal(t, subscription.Salary, tc.subscription.Salary)
			assert.Equal(t, subscription.Type, tc.subscription.Type)
			assert.Equal(t, *subscription.IsRemoteFriendly, *tc.subscription.IsRemoteFriendly)
			for _, k1 := range subscription.Keywords {
				found := false
				for _, k2 := range tc.subscription.Keywords {
					if k1.Value == k2 {
						found = true
						break
					}
				}
				assert.True(t, found)
			}
		})
	}
}