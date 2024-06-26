package notification

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"text/template"

	"github.com/ariel17/jobberwocky/configs"
	"github.com/ariel17/jobberwocky/internal/core/domain"
	"github.com/ariel17/jobberwocky/internal/core/ports"
)

type notificationService struct {
	repository   ports.SubscriptionRepository
	emailClient  ports.EmailProviderClient
	templatePath string
	workers      int
	input        chan domain.Job
	output       chan error
}

func NewNotificationService(workers int, repository ports.SubscriptionRepository, emailClient ports.EmailProviderClient, templatePath string) ports.NotificationService {
	return &notificationService{
		repository:   repository,
		emailClient:  emailClient,
		templatePath: templatePath,
		workers:      workers,
		input:        make(chan domain.Job),
		output:       make(chan error),
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
		log.Print("Looking for subscriptions to notify ...")
		subscriptions, err := n.repository.Filter(job)
		if err != nil {
			log.Printf("failed to retrieve subscribers: %v", err)
			continue
		}

		log.Printf("Found %d subscriptions to be notified: %v", len(subscriptions), subscriptions)
		for _, subscription := range subscriptions {
			subject := fmt.Sprintf("%s %s", configs.GetEmailSubject(), job.Title)
			body, err := createBody(n.templatePath, job)
			if err != nil {
				log.Printf("failed to create email body: %v", err)
				continue
			}
			if err = n.emailClient.Send(configs.GetEmailFrom(), subscription.Email, subject, body); err != nil {
				log.Printf("failed to send email: %v", err)
			}
		}

		log.Print("Subscription notifications finished.")
	}
}

func createBody(templateFile string, job domain.Job) (string, error) {
	templateString, err := os.ReadFile(templateFile)
	if err != nil {
		return "", err
	}
	tmpl, err := template.New(templateFile).Parse(string(templateString))
	if err != nil {
		return "", err
	}
	var b bytes.Buffer
	if err := tmpl.Execute(&b, job); err != nil {
		return "", err
	}
	return b.String(), nil
}