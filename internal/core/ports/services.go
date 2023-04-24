package ports

import "github.com/ariel17/jobberwocky/internal/core/domain"

// JobService is responsible for storing new services and retrieve all those
// that match a given pattern, if present.
type JobService interface {
	Match() ([]domain.Job, error)
	Create(domain.Job) error
}

// NotificationService sends job alerts when its details matches subscribers'
// patterns.
type NotificationService interface {
	Send(domain.Job) error
}