package services

import (
	"github.com/ariel17/jobberwocky/internal/core/domain"
	"github.com/ariel17/jobberwocky/internal/core/ports"
)

type jobService struct {
	repository ports.JobRepository
}

func NewJobService(repository ports.JobRepository) ports.JobService {
	return &jobService{
		repository: repository,
	}
}

func (j *jobService) Match(pattern *domain.Filter) ([]domain.Job, error) {
	return nil, nil
}

func (j *jobService) Create(job domain.Job) error {
	return j.repository.Save(job)
}