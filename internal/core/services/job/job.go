package job

import (
	"log"

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
	if err := j.repository.Save(job); err != nil {
		return err
	}
	log.Printf("New job created: %v", job)
	j.notifications.Enqueue(job)
	return nil
}