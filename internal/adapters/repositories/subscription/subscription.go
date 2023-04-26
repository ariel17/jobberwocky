package subscription

import (
	"gorm.io/gorm"

	"github.com/ariel17/jobberwocky/internal/adapters/repositories"
	"github.com/ariel17/jobberwocky/internal/core/domain"
	"github.com/ariel17/jobberwocky/internal/core/ports"
)

type subscriptionRepository struct {
	db *gorm.DB
}

func NewSubscriptionRepository(db *gorm.DB) (ports.SubscriptionRepository, error) {
	return &subscriptionRepository{
		db: db,
	}, nil
}

func (s *subscriptionRepository) Filter(job domain.Job) ([]domain.Subscription, error) {
	return nil, nil
}
func (s *subscriptionRepository) Save(subscription domain.Subscription) error {
	sm := subscriptionDomainToModel(subscription)
	tx := s.db.Create(&sm)
	if tx.Error != nil {
		return tx.Error
	}
	needsReplacement, newKeywords, err := repositories.ReuseExistingKeywords(s.db, sm.Keywords)
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
	if err := s.db.AutoMigrate(&repositories.Keyword{}); err != nil {
		return err
	}
	return nil
}