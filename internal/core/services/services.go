package services

import (
	"github.com/ariel17/jobberwocky/internal/core/domain"
	"github.com/ariel17/jobberwocky/internal/core/ports"
)

type jobService struct {
	repository    ports.JobRepository
	notifications ports.NotificationService
}

func NewJobService(repository ports.JobRepository, notifications ports.NotificationService) ports.JobService {
	return &jobService{
		repository:    repository,
		notifications: notifications,
	}
}

func (j *jobService) Match(pattern *domain.Filter) ([]domain.Job, error) {
	return j.repository.Filter(pattern)
}

func (j *jobService) Create(job domain.Job) error {
	err := j.repository.Save(job)
	if err == nil {
		j.notifications.Enqueue(job)
	}
	return err
}