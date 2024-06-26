package configs

import (
	"fmt"
	"os"
)

const (
	emailFromKey     = "EMAIL_FROM"
	defaultEmailFrom = "jobs@example.com"

	emailSubjectKey     = "EMAIL_SUBJECT"
	defaultEmailSubject = "A new job alert has arrived!"

	emailTemplateKey = "EMAIL_TEMPLATE"
	DefaultTemplate  = "./resources/body.tmpl"
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
	if emailFrom = os.Getenv(emailFromKey); emailFrom == "" {
		emailFrom = defaultEmailFrom
	}

	if emailSubject = os.Getenv(emailSubjectKey); emailSubject == "" {
		emailSubject = defaultEmailSubject
	}

	template := os.Getenv(emailTemplateKey)
	if template == "" {
		template = DefaultTemplate
	}
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	emailTemplate = fmt.Sprintf("%s/%s", wd, template)
}