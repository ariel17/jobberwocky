package configs

import "os"

const (
	jobberwockyURLKey = "JOBBERWOCKY_URL"
)

var (
	jobberwockyURL = ""
)

func GetJobberwockyURL() string {
	return jobberwockyURL
}

func init() {
	jobberwockyURL = os.Getenv(jobberwockyURLKey)
}