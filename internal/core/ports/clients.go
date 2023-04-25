package ports

type ExternalJobClient interface {
	JobFilter
	Name() string
}

type EmailProviderClient interface {
	Send(from, to, subject, body string) error
}