package clients

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/ariel17/jobberwocky/internal/core/domain"
	"github.com/ariel17/jobberwocky/internal/core/ports"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

func NewJobberwockyExternalJobClient(client HTTPClient) ports.ExternalJobClient {
	return &jobberwockyExternalJobClient{
		client: client,
	}

}

type jobberwockyExternalJobClient struct {
	client HTTPClient
}

func (j *jobberwockyExternalJobClient) Filter(pattern *domain.Pattern) ([]domain.Job, error) {
	// req, err := http.NewRequest("GET", configs.GetJobberwockyURL()+"/jobs", nil)
	// if err != nil {
	// 	return nil, err
	// }
	return nil, nil
}

func patternToQueryString(pattern *domain.Pattern) string {
	if pattern == nil {
		return ""
	}
	fields := []string{}
	if pattern.Text != "" {
		fields = append(fields, "name="+pattern.Text)
	}
	if pattern.Location != "" {
		fields = append(fields, "country="+pattern.Location)
	}
	if pattern.Salary > 0 {
		fields = append(fields, fmt.Sprintf("salary_min=%d", pattern.Salary), fmt.Sprintf("salary_max=%d", pattern.Salary))
	}
	if len(fields) == 0 {
		return ""
	}
	return fmt.Sprintf("?%s", strings.Join(fields, "&"))
}