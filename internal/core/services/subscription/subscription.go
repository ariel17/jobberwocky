package subscription

import (
	"log"

	"github.com/ariel17/jobberwocky/internal/core/domain"
	"github.com/ariel17/jobberwocky/internal/core/ports"
)

type subscriptionService struct {
	repository ports.SubscriptionRepository
}

func NewSubscriptionService(repository ports.SubscriptionRepository) ports.SubscriptionService {
	return &subscriptionService{
		repository: repository,
	}
}

func (s *subscriptionService) Create(subscription domain.Subscription) error {
	err := s.repository.Save(subscription)
	if err != nil {
		return err
	}
	log.Printf("New subscription created: %v", subscription)
	return nil
}