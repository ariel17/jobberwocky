package subscription

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/ariel17/jobberwocky/internal/adapters/repositories/job"
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
		{"success", domain.Subscription{Pattern: domain.Pattern{Text: "Title", Company: "Company", Location: "Argentina", Salary: 60, Type: domain.FullTime, IsRemoteFriendly: helpers.BoolPointer(true), Keywords: []string{"k1", "k2", "k3"}}, Email: "arielgerardorios@gmail.com"}},
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

func TestSubscriptionRepository_Filter(t *testing.T) {
	subscriptions := []domain.Subscription{
		{Pattern: domain.Pattern{}, Email: "user1@example.com"},
		{Pattern: domain.Pattern{Text: "unique"}, Email: "user2@example.com"},
		{Pattern: domain.Pattern{Text: "wow!"}, Email: "user3@example.com"},
		{Pattern: domain.Pattern{Company: "X"}, Email: "user4@example.com"},
		{Pattern: domain.Pattern{Type: domain.FullTime}, Email: "user5@example.com"},
		{Pattern: domain.Pattern{Keywords: []string{"k1"}}, Email: "user6@example.com"},
		{Pattern: domain.Pattern{Keywords: []string{"k5"}}, Email: "user7@example.com"},
	}
	testCases := []struct {
		name     string
		job      domain.Job
		expected int
	}{
		{"matches by empty filter and type and keywords", domain.Job{Title: "This is the title", Description: "The description", Company: "Company1", Location: "Argentina", SalaryMin: 60, SalaryMax: 80, Type: domain.FullTime, IsRemoteFriendly: helpers.BoolPointer(true), Keywords: []string{"k1", "k2", "k3"}}, 3},
		{"matches by empty filter and text in title and keywords", domain.Job{Title: "Unique title", Description: "The description", Company: "Company2", Location: "Argentina", SalaryMin: 60, SalaryMax: 80, Type: domain.FullTime, IsRemoteFriendly: helpers.BoolPointer(true), Keywords: []string{"k1", "k2", "k3"}}, 3},
		{"matches by empty filter and text in description and keywords", domain.Job{Title: "title1", Description: "Wow!", Company: "Company3", Location: "Argentina", SalaryMin: 60, SalaryMax: 80, Type: domain.FullTime, IsRemoteFriendly: helpers.BoolPointer(true), Keywords: []string{"k1", "k2", "k3"}}, 3},
		{"matches by empty filter and company and type and keywords", domain.Job{Title: "title2", Description: "Wow!", Company: "X", Location: "Argentina", SalaryMin: 60, SalaryMax: 80, Type: domain.FullTime, IsRemoteFriendly: helpers.BoolPointer(true), Keywords: []string{"k1", "k2", "k3"}}, 4},
		{"matches by only keywords", domain.Job{Title: "title2", Description: "Wow!", Company: "X", Location: "Argentina", SalaryMin: 60, SalaryMax: 80, Type: domain.FullTime, IsRemoteFriendly: helpers.BoolPointer(true), Keywords: []string{"k5"}}, 4},
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

			jobRepository := job.NewJobRepository(db)
			assert.Nil(t, jobRepository.SyncSchemas())
			assert.Nil(t, jobRepository.Save(tc.job))

			repository := NewSubscriptionRepository(db)
			assert.Nil(t, repository.SyncSchemas())

			for _, s := range subscriptions {
				assert.Nil(t, repository.Save(s))
			}

			matchingSubscriptions, err := repository.Filter(tc.job)
			assert.Nil(t, err)
			assert.NotNil(t, matchingSubscriptions)
			assert.Equal(t, tc.expected, len(matchingSubscriptions))
		})
	}
}