package clients

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/ariel17/jobberwocky/internal/configs"
	"github.com/ariel17/jobberwocky/internal/core/domain"
	"github.com/ariel17/jobberwocky/internal/core/ports"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type jobberwockyJobs [][]interface{}

func NewJobberwockyExternalJobClient(client HTTPClient) ports.ExternalJobClient {
	return &jobberwockyExternalJobClient{
		client: client,
	}

}

type jobberwockyExternalJobClient struct {
	client HTTPClient
}

func (j *jobberwockyExternalJobClient) Filter(pattern *domain.Pattern) ([]domain.Job, error) {
	qs := patternToQueryString(pattern)
	req, err := http.NewRequest("GET", configs.GetJobberwockyURL()+"/jobs"+qs, nil)
	if err != nil {
		return nil, err
	}
	response, err := j.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("api error: %d", response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	var rawJobs jobberwockyJobs
	if err = json.Unmarshal(body, &rawJobs); err != nil {
		return nil, fmt.Errorf("invalid body: %s", body)
	}

	jobs := []domain.Job{}
	for _, rawJob := range rawJobs {
		title := rawJob[0].(string)
		location := rawJob[2].(string)
		salaryMax := int(rawJob[1].(float64))
		keywords := []string{}
		for _, k := range rawJob[3].([]interface{}) {
			keywords = append(keywords, k.(string))
		}
		j, err := domain.NewJob(title, "", "", location, 0, salaryMax, "", nil, j.Name(), keywords...)
		if err != nil {
			return nil, err
		}
		jobs = append(jobs, j)
	}
	return jobs, nil
}

func (j *jobberwockyExternalJobClient) Name() string {
	return "JobberwockyExternal"
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