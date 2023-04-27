package configs

import "os"

const (
	httpAddressKey     = "HTTP_ADDRESS"
	defaultHTTPAddress = ":8080"
)

var (
	httpAddress = ""
)

func GetHTTPAddress() string {
	return httpAddress
}

func init() {
	if httpAddress = os.Getenv(httpAddressKey); httpAddress == "" {
		httpAddress = defaultHTTPAddress
	}
}