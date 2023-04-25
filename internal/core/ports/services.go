package ports

import "github.com/ariel17/jobberwocky/internal/core/domain"

// JobService is responsible for storing new services and retrieve all those
// that match a given pattern, if present.
type JobService interface {
	Filter(pattern *domain.Filter) ([]domain.Job, error)
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