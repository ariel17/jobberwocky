package ports

import "github.com/ariel17/jobberwocky/internal/core/domain"

type JobService interface {
	Match() []domain.Job
	Create(domain.Job) error
}

type NotificationService interface {
	Notify(domain.Job) error
}

type JobRepository interface {
	Filter() ([]domain.Job, error)
	Save(job domain.Job) error
}

type ExternalService interface {
	Match() ([]domain.Job, error)
}