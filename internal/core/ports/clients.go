package ports

import "github.com/ariel17/jobberwocky/internal/core/domain"

type ExternalJobClient interface {
	JobFilter
	Name() string
	PatternIsSearchable(pattern *domain.Pattern) bool
}

type EmailProviderClient interface {
	Send(from, to, subject, body string) error
}