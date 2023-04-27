package configs

import (
	"log"
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
	if notificationWorkers, err = strconv.Atoi(workers); err != nil {
		log.Printf("Using default workers: error=%v, workers=%d", err, defaultNotificationWorkers)
		notificationWorkers = defaultNotificationWorkers
	}
}