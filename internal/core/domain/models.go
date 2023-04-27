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

// isSalaryValid checks for correct ranges or fixed values to be correct.
func (j Job) isSalaryValid() bool {
	if j.SalaryMin > 0 {
		return j.SalaryMax > j.SalaryMin
	}
	return j.SalaryMax > 0
}

// isLocationAndIsRemoteFriendlyValid checks combination of both fields ensure
// that no location is remote-friendly.
func (j Job) isLocationAndIsRemoteFriendlyValid() bool {
	if j.Location == "" && j.IsRemoteFriendly != nil {
		return *j.IsRemoteFriendly == false
	}
	return true
}

func (j Job) isSourceValid(isLocal bool) bool {
	if isLocal {
		return j.Source == ""
	}
	return j.Source != ""
}

func (j Job) Validate(isLocal bool) error {
	if j.Title == "" {
		return errors.New("title cannot be empty")
	}
	if !isTypeValid(j.Source, j.Type) {
		return fmt.Errorf("type value is invalid: %s", j.Type)
	}
	if !j.isSalaryValid() {
		return fmt.Errorf("fixed/ranged salary is invalid: min=%d, max=%d", j.SalaryMin, j.SalaryMax)
	}
	if !j.isLocationAndIsRemoteFriendlyValid() {
		return fmt.Errorf("location and remote-friendly values are incorrect: location=%s, is_remote_friendly=%v", j.Location, *j.IsRemoteFriendly)
	}
	if !isKeywordsValid(j.Keywords) {
		return errors.New("keywords needs to be unique")
	}
	if !j.isSourceValid(isLocal) {
		if isLocal {
			return errors.New("source is not allowed")
		}
		return errors.New("source is required")
	}
	return nil
}

// Pattern contains value patterns to match when searching for matching jobs.
type Pattern struct {
	Text             string   `form:"text" json:"text"`
	Company          string   `form:"company" json:"company"`
	Location         string   `form:"location" json:"location"`
	Salary           int      `form:"salary" json:"salary"`
	Type             string   `form:"type" json:"type"`
	IsRemoteFriendly *bool    `form:"is_remote_friendly" json:"is_remote_friendly,omitempty"`
	Keywords         []string `form:"keywords" json:"keywords"`
}

func (p Pattern) IsEmpty() bool {
	return p.Text == "" && p.Company == "" && p.Location == "" &&
		p.Salary == 0 && p.Type == "" && p.IsRemoteFriendly == nil &&
		len(p.Keywords) == 0
}

func (p Pattern) Validate() error {
	if !isKeywordsValid(p.Keywords) {
		return errors.New("keywords needs to be unique")
	}
	if p.Type != "" && !isTypeValid("", p.Type) {
		return fmt.Errorf("type value is invalid: %s", p.Type)
	}
	return nil
}

// Subscription contains the contact details to notify a person about a matching
// job.
type Subscription struct {
	Pattern
	Email string
}

func (s Subscription) Validate() error {
	if err := s.Pattern.Validate(); err != nil {
		return err
	}
	if s.Email == "" {
		return errors.New("email cannot be empty")
	}
	return nil
}

// isKeywordsValid checks for repeated values in list.
func isKeywordsValid(keywords []string) bool {
	for index, k1 := range keywords {
		for _, k2 := range keywords[index+1:] {
			if k1 == k2 {
				return false
			}
		}
	}
	return true
}

// isTypeValid checks that `Type` field only contains specific values.
func isTypeValid(source, jobType string) bool {
	if source == "" {
		for _, t := range []string{Contractor, FullTime, PartTime} {
			if jobType == t {
				return true
			}
		}
		return false
	}
	return true
}