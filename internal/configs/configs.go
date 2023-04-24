package configs

import (
	"fmt"
	"os"
	"strconv"
)

const (
	EmailFromKey    = "EMAIL_FROM"
	EmailSubjectKey = "EMAIL_SUBJECT"

	NotificationWorkersKey = "NOTIFICATION_WORKERS"
)

var (
	emailFrom           = ""
	emailSubject        = ""
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

func init() {
	emailFrom = os.Getenv(EmailFromKey)
	emailSubject = os.Getenv(EmailSubjectKey)

	workers := os.Getenv(NotificationWorkersKey)
	var err error
	notificationWorkers, err = strconv.Atoi(workers)
	if err != nil {
		panic(fmt.Errorf("notification workers cannot be parsed as int: %v", err))
	}
}