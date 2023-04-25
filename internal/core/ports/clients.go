package ports

type ExternalJobClient interface {
	JobFilter
}

type EmailProviderClient interface {
	Send(from, to, subject, body string) error
}