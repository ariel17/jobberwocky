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
	Title            string
	Description      string
	Company          string
	Location         string
	SalaryMin        int
	SalaryMax        int
	Type             string
	IsRemoteFriendly *bool
	Keywords         []string
	Source           string
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
	} else if j.Type != "" {
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
		return Job{}, fmt.Errorf("location and remote-friendly values are incorrect: location=%s, remote friendly=%v", location, isRemoteFriendly)
	}
	return job, nil
}

// Pattern contains value patterns to match when searching for matching jobs.
type Pattern struct {
	Text             string
	Company          string
	Location         string
	Salary           int
	Type             string
	IsRemoteFriendly *bool
	Keywords         []string
}

// Subscription contains the contact details to notify a person about a matching
// job.
type Subscription struct {
	Pattern
	Email string
}