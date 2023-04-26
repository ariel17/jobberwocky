package job

import (
	"log"
	"sync"

	"github.com/ariel17/jobberwocky/internal/core/domain"
	"github.com/ariel17/jobberwocky/internal/core/ports"
)

type jobService struct {
	repository    ports.JobRepository
	notifications ports.NotificationService
	externals     []ports.ExternalJobClient
}

func NewJobService(
	repository ports.JobRepository,
	notifications ports.NotificationService,
	externals ...ports.ExternalJobClient) ports.JobService {

	return &jobService{
		repository:    repository,
		notifications: notifications,
		externals:     externals,
	}
}

func (j *jobService) Filter(pattern *domain.Pattern) ([]domain.Job, error) {
	jobsOutput := make(chan []domain.Job, len(j.externals)+1)
	errOutput := make(chan error, len(j.externals)+1)

	var wg sync.WaitGroup
	wg.Add(1)
	go asyncFilter(j.repository, pattern, &wg, jobsOutput, errOutput)
	for _, external := range j.externals {
		wg.Add(1)
		go asyncFilter(external, pattern, &wg, jobsOutput, errOutput)
	}
	wg.Wait()
	close(errOutput)
	close(jobsOutput)

	var lastErr error
	for err := range errOutput {
		if err != nil {
			log.Printf("Failed to obtain jobs from sources: %v", err)
			lastErr = err
		}
	}
	if lastErr != nil {
		return nil, lastErr
	}

	results := make([]domain.Job, 0)
	for jobs := range jobsOutput {
		results = append(results, jobs...)
	}

	return results, nil
}

func (j *jobService) Create(job domain.Job) error {
	if err := j.repository.Save(job); err != nil {
		return err
	}
	log.Printf("New job created: %v", job)
	j.notifications.Enqueue(job)
	return nil
}

func asyncFilter(filter ports.JobFilter, pattern *domain.Pattern, wg *sync.WaitGroup, jobsOutput chan []domain.Job, errOutput chan error) {
	defer wg.Done()
	jobs, err := filter.Filter(pattern)
	jobsOutput <- jobs
	errOutput <- err
	log.Printf("async job filtering result: jobs=%d, err=%v", len(jobs), err)
}