package ports

import "github.com/ariel17/jobberwocky/internal/core/domain"

type JobRepository interface {
	Filter() ([]domain.Job, error)
	Save(job domain.Job) error
}