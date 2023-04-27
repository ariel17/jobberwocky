package clients

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/ariel17/jobberwocky/internal/core/domain"
	"github.com/ariel17/jobberwocky/internal/core/ports"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type jobberwockyJobs [][]interface{}

type jobberwockyExternalJobClient struct {
	client HTTPClient
	url    string
}

func NewJobberwockyExternalJobClient(client HTTPClient, url string) ports.ExternalJobClient {
	return &jobberwockyExternalJobClient{
		client: client,
		url:    url,
	}
}

func (j *jobberwockyExternalJobClient) Filter(pattern *domain.Pattern) ([]domain.Job, error) {
	url := j.url + "/jobs" + patternToQueryString(pattern)
	log.Printf("Requesting external resource: url=%s", url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	response, err := j.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = response.Body.Close()
	}()

	log.Printf("External resource response: status_code=%d, content_length=%d", response.StatusCode, response.ContentLength)

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

	jobs := make([]domain.Job, 0)
	for _, rawJob := range rawJobs {
		title := rawJob[0].(string)
		location := rawJob[2].(string)
		salaryMax := int(rawJob[1].(float64))
		keywords := make([]string, 0)
		for _, k := range rawJob[3].([]interface{}) {
			keywords = append(keywords, k.(string))
		}
		j := domain.Job{Title: title, Description: "", Company: "", Location: location, SalaryMin: 0, SalaryMax: salaryMax, Type: "", IsRemoteFriendly: nil, Source: j.Name(), Keywords: keywords}
		if err := j.Validate(false); err != nil {
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
	fields := make([]string, 0)
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