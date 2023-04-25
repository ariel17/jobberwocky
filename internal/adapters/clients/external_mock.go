package clients

import "github.com/ariel17/jobberwocky/internal/adapters/repositories"

type MockExternalJobClient struct {
	repositories.MockFilter
}