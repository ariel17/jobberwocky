package services

import (
	"fmt"
	"log"

	"github.com/ariel17/jobberwocky/internal/configs"
	"github.com/ariel17/jobberwocky/internal/core/domain"
	"github.com/ariel17/jobberwocky/internal/core/ports"
)

type notificationService struct {
	repository  ports.SubscriptionRepository
	emailClient ports.EmailProviderClient
	workers     int
	input       chan domain.Job
	output      chan error
}

func NewNotificationService(workers int, repository ports.SubscriptionRepository, emailClient ports.EmailProviderClient) ports.NotificationService {
	return &notificationService{
		repository:  repository,
		emailClient: emailClient,
		workers:     workers,
		input:       make(chan domain.Job),
		output:      make(chan error),
	}
}

func (n *notificationService) Enqueue(job domain.Job) {
	n.input <- job
}

func (n *notificationService) StartWorkers() {
	log.Printf("starting %d workers ...", n.workers)
	for i := 0; i < n.workers; i++ {
		go n.Process()
	}
	log.Print("workers started")
}

func (n *notificationService) StopWorkers() {
	log.Print("stoping workers ...")
	close(n.input)
	close(n.output)
	log.Print("workers stopped")
}

func (n *notificationService) Process() {
	defer func() {
		if r := recover(); r != nil {
			log.Print("recovered from panic in notification process:", r)
		}
	}()

	for job := range n.input {
		subscriptions, err := n.repository.Filter(job)
		if err != nil {
			log.Printf("failed to retrieve subscribers: %v", err)
			continue
		}
		for _, subscription := range subscriptions {
			subject := fmt.Sprintf("%s: %s", configs.GetEmailSubject(), job.Title)
			// TODO generate body
			body := ""
			err = n.emailClient.Send(configs.GetEmailFrom(), subscription.Email, subject, body)
			if err != nil {
				log.Printf("failed to send email: %v", err)
			}
		}
	}
}