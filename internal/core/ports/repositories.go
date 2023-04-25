package ports

import "github.com/ariel17/jobberwocky/internal/core/domain"

type Repository interface {
	SyncSchemas() error
}

type JobRepository interface {
	Repository
	JobFilter
	Save(job domain.Job) error
}

type SubscriptionRepository interface {
	Repository
	Filter(job domain.Job) ([]domain.Subscription, error)
	Save(subscription domain.Subscription) error
}