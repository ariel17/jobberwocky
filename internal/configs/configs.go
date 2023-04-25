package configs

import (
	"os"
	"strconv"
)

const (
	EmailFromKey     = "EMAIL_FROM"
	DefaultEmailFrom = "jobs@example.com"

	EmailSubjectKey     = "EMAIL_SUBJECT"
	DefaultEmailSubject = "A new job alert has arrived"

	EmailTemplateKey = "EMAIL_TEMPLATE"
	DefaultTemplate  = "body.tmpl"

	NotificationWorkersKey     = "NOTIFICATION_WORKERS"
	DefaultNotificationWorkers = 10
)

var (
	emailFrom           = ""
	emailSubject        = ""
	emailTemplate       = ""
	notificationWorkers = 0
)

func GetEmailFrom() string {
	return emailFrom
}

func GetEmailSubject() string {
	return emailSubject
}

func GetNotificationWorkers() int {
	return notificationWorkers
}

func GetEmailTemplate() string {
	return emailTemplate
}

func init() {
	emailFrom = os.Getenv(EmailFromKey)
	if emailFrom == "" {
		emailFrom = DefaultEmailFrom
	}

	emailSubject = os.Getenv(EmailSubjectKey)
	if emailSubject == "" {
		emailSubject = DefaultEmailSubject
	}

	emailTemplate = os.Getenv(EmailTemplateKey)
	if emailTemplate == "" {
		emailTemplate = DefaultTemplate
	}

	workers := os.Getenv(NotificationWorkersKey)
	var err error
	notificationWorkers, err = strconv.Atoi(workers)
	if err != nil {
		notificationWorkers = DefaultNotificationWorkers
	}
}