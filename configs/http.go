package configs

import "os"

const (
	httpPortKey     = "PORT"
	defaultHTTPPort = "8080"
)

var (
	httpPort = ""
)

func GetHTTPPort() string {
	return httpPort
}

func init() {
	if httpPort = os.Getenv(httpPortKey); httpPort == "" {
		httpPort = defaultHTTPPort
	}
}