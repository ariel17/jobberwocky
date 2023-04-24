package ports

import "github.com/ariel17/jobberwocky/internal/core/domain"

type JobService interface {
	Match() []domain.Job
	Create(domain.Job) error
}

type NotificationService interface {
	Send(domain.Job) error
}