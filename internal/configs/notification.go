package configs

import (
	"os"
	"strconv"
)

const (
	notificationWorkersKey     = "NOTIFICATION_WORKERS"
	defaultNotificationWorkers = 10
)

var (
	notificationWorkers = 0
)

func GetNotificationWorkers() int {
	return notificationWorkers
}

func init() {
	workers := os.Getenv(notificationWorkersKey)
	var err error
	notificationWorkers, err = strconv.Atoi(workers)
	if err != nil {
		notificationWorkers = defaultNotificationWorkers
	}
}