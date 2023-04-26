package subscription

import (
	"fmt"
	"strings"

	"gorm.io/gorm"

	"github.com/ariel17/jobberwocky/internal/adapters/repositories/keyword"
	"github.com/ariel17/jobberwocky/internal/core/domain"
	"github.com/ariel17/jobberwocky/internal/core/ports"
)

type subscriptionRepository struct {
	db *gorm.DB
}

func NewSubscriptionRepository(db *gorm.DB) ports.SubscriptionRepository {
	return &subscriptionRepository{
		db: db,
	}
}

func (s *subscriptionRepository) Filter(job domain.Job) ([]domain.Subscription, error) {
	var modelSubscriptions []Subscription
	query, parameters := jobToQueryAndParameters(job)
	tx := s.db.Where(query, parameters).Find(&modelSubscriptions)
	if tx.Error != nil {
		return nil, tx.Error
	}

	domainSubscriptions := []domain.Subscription{}
	for _, subscription := range modelSubscriptions {
		s := subscriptionModelToDomain(subscription)
		domainSubscriptions = append(domainSubscriptions, s)
	}

	return domainSubscriptions, nil
}
func (s *subscriptionRepository) Save(subscription domain.Subscription) error {
	sm := subscriptionDomainToModel(subscription)
	tx := s.db.Create(&sm)
	if tx.Error != nil {
		return tx.Error
	}
	needsReplacement, newKeywords, err := keyword.ReuseExistingKeywords(s.db, sm.Keywords)
	if err != nil {
		return err
	}
	if needsReplacement {
		return s.db.Model(&sm).Association("Keywords").Replace(newKeywords)
	}
	return nil
}

func (s *subscriptionRepository) SyncSchemas() error {
	if err := s.db.AutoMigrate(&Subscription{}); err != nil {
		return err
	}
	if err := s.db.AutoMigrate(&keyword.Keyword{}); err != nil {
		return err
	}
	return nil
}

func jobToQueryAndParameters(job domain.Job) (string, map[string]interface{}) {
	query := []string{}
	parameters := make(map[string]interface{})

	query = append(query, "(text = '' OR INSTR(@title, text) > 0 OR INSTR(@description, text) > 0)")
	parameters["title"] = job.Title
	parameters["description"] = job.Description

	if job.Company != "" {
		query = append(query, "(company = '' OR company = @company)")
		parameters["company"] = job.Company
	}

	query = append(query, "(location = '' OR location = @location)")
	parameters["location"] = job.Location

	if job.SalaryMin > 0 {
		query = append(query, "(salary = 0 OR salary >= @salary_min)")
		parameters["salary_min"] = job.SalaryMin
	}

	query = append(query, "(salary = 0 OR salary <= @salary_max)")
	parameters["salary_max"] = job.SalaryMax

	query = append(query, "(type = '' OR type = @type)")
	parameters["type"] = job.Type

	if job.IsRemoteFriendly != nil {
		query = append(query, "(is_remote_friendly IS NULL OR is_remote_friendly = @is_remote_friendly)")
		parameters["is_remote_friendly"] = *job.IsRemoteFriendly
	}

	if len(job.Keywords) > 0 {
		sub1 := `SELECT COUNT(*) FROM subscriptions_keywords sk WHERE sk.subscription_id = id`
		sub2 := `SELECT sk.subscription_id AS subscription_id, COUNT(*)
			FROM subscriptions_keywords sk
			INNER JOIN jobs_keywords jk ON (sk.keyword_id=jk.keyword_id)
			GROUP BY sk.subscription_id
			HAVING COUNT(*) <= @count`
		query = append(query, fmt.Sprintf("((%s) = 0 OR id IN (SELECT subscription_id FROM (%s)))", sub1, sub2))
		parameters["count"] = len(job.Keywords)
	}
	return strings.Join(query, " AND "), parameters
}