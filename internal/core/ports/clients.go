package ports

import "github.com/ariel17/jobberwocky/internal/core/domain"

type ExternalJobClient interface {
	Match() ([]domain.Job, error)
}