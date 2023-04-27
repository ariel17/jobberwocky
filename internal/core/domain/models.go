package domain

import (
	"errors"
	"fmt"
)

const (
	FullTime   = "Full-Time"
	Contractor = "Contractor"
	PartTime   = "Part-Time"
)

// Job contains the details of a hiring position on a company.
type Job struct {
	Title            string   `json:"title"`
	Description      string   `json:"description"`
	Company          string   `json:"company"`
	Location         string   `json:"location"`
	SalaryMin        int      `json:"salary_min"`
	SalaryMax        int      `json:"salary_max"`
	Type             string   `json:"type"`
	IsRemoteFriendly *bool    `json:"is_remote_friendly,omitempty"`
	Keywords         []string `json:"keywords"`
	Source           string   `json:"source"`
}

func (j Job) IsTitleValid() bool {
	return j.Title != ""
}

// IsTypeValid checks that `Type` field only contains specific values.
func (j Job) IsTypeValid() bool {
	if j.Source == "" {
		for _, t := range []string{Contractor, FullTime, PartTime} {
			if j.Type == t {
				return true
			}
		}
		return false
	}
	return true
}

// IsSalaryValid checks for correct ranges or fixed values to be correct.
func (j Job) IsSalaryValid() bool {
	if j.SalaryMin > 0 {
		return j.SalaryMax > j.SalaryMin
	}
	return j.SalaryMax > 0
}

// IsLocationAndIsRemoteFriendlyValid checks combination of both fields ensure
// that no location is remote-friendly.
func (j Job) IsLocationAndIsRemoteFriendlyValid() bool {
	if j.Location == "" && j.IsRemoteFriendly != nil {
		return *j.IsRemoteFriendly == true
	}
	return true
}

// IsKeywordsValid checks for repeated values in list.
func (j Job) IsKeywordsValid() bool {
	for index, keyword1 := range j.Keywords {
		for _, keyword2 := range j.Keywords[index+1:] {
			if keyword1 == keyword2 {
				return false
			}
		}
	}
	return true
}

// NewJob creates a new job instance and ensures field values are valid.
func NewJob(title, description, company, location string, salaryMin, salaryMax int, jobType string, isRemoteFriendly *bool, source string, keywords ...string) (Job, error) {
	job := Job{
		Title:            title,
		Description:      description,
		Company:          company,
		Location:         location,
		SalaryMin:        salaryMin,
		SalaryMax:        salaryMax,
		Type:             jobType,
		IsRemoteFriendly: isRemoteFriendly,
		Keywords:         keywords,
		Source:           source,
	}
	if !job.IsTitleValid() {
		return Job{}, errors.New("title cannot be empty")
	}
	if !job.IsTypeValid() {
		return Job{}, fmt.Errorf("type value is invalid: %s", jobType)
	}
	if !job.IsSalaryValid() {
		return Job{}, fmt.Errorf("fixed/ranged salary is invalid: min=%d, max=%d", salaryMin, salaryMax)
	}
	if !job.IsLocationAndIsRemoteFriendlyValid() {
		return Job{}, fmt.Errorf("location and remote-friendly values are incorrect: location=%s, remote friendly=%v", location, *isRemoteFriendly)
	}
	if !job.IsKeywordsValid() {
		return Job{}, errors.New("keywords needs to be unique")
	}
	return job, nil
}

// Pattern contains value patterns to match when searching for matching jobs.
type Pattern struct {
	Text             string   `form:"text"`
	Company          string   `form:"company"`
	Location         string   `form:"location"`
	Salary           int      `form:"salary"`
	Type             string   `form:"type"`
	IsRemoteFriendly *bool    `form:"is_remote_friendly"`
	Keywords         []string `form:"keywords"`
}

func (p Pattern) IsEmpty() bool {
	return p.Text == "" && p.Company == "" && p.Location == "" &&
		p.Salary == 0 && p.Type == "" && p.IsRemoteFriendly == nil &&
		len(p.Keywords) == 0
}

// Subscription contains the contact details to notify a person about a matching
// job.
type Subscription struct {
	Pattern
	Email string
}