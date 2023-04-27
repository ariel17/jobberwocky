package ports

import "github.com/ariel17/jobberwocky/internal/core/domain"

type JobFilter interface {
	Filter(pattern *domain.Pattern) ([]domain.Job, error)
}

// JobService is responsible for storing new services and retrieve all those
// that match a given pattern, if present.
type JobService interface {
	JobFilter
	Create(job domain.Job) error
}

// NotificationService sends job alerts when its details matches subscribers'
// patterns.
type NotificationService interface {
	Enqueue(job domain.Job)
	StartWorkers()
	StopWorkers()
	Process()
}

// SubscriptionService takes new subscriptions to be saved for future
// notifications.
type SubscriptionService interface {
	Create(subscription domain.Subscription) error
}