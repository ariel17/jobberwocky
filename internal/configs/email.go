package configs

import (
	"os"
)

const (
	emailFromKey     = "EMAIL_FROM"
	defaultEmailFrom = "jobs@example.com"

	emailSubjectKey     = "EMAIL_SUBJECT"
	defaultEmailSubject = "A new job alert has arrived"

	emailTemplateKey = "EMAIL_TEMPLATE"
	DefaultTemplate  = "body.tmpl"
)

var (
	emailFrom     = ""
	emailSubject  = ""
	emailTemplate = ""
)

func GetEmailFrom() string {
	return emailFrom
}

func GetEmailSubject() string {
	return emailSubject
}

func GetEmailTemplate() string {
	return emailTemplate
}

func init() {
	emailFrom = os.Getenv(emailFromKey)
	if emailFrom == "" {
		emailFrom = defaultEmailFrom
	}

	emailSubject = os.Getenv(emailSubjectKey)
	if emailSubject == "" {
		emailSubject = defaultEmailSubject
	}

	emailTemplate = os.Getenv(emailTemplateKey)
	if emailTemplate == "" {
		emailTemplate = DefaultTemplate
	}
}