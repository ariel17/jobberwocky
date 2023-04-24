package ports

import "github.com/ariel17/jobberwocky/internal/core/domain"

type JobRepository interface {
	Filter(pattern *domain.Filter) ([]domain.Job, error)
	Save(job domain.Job) error
}

type SubscriptionRepository interface {
	Filter(job domain.Job) ([]domain.Subscription, error)
	Save(subscription domain.Subscription) error
}