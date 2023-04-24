package domain

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
	IsRemoteFriendly bool
	Keywords         []string
}

// IsTypeValid checks that `Type` field only contains specific values.
func (j Job) IsTypeValid() bool {
	for _, t := range []string{Contractor, FullTime, PartTime} {
		if j.Type == t {
			return true
		}
	}
	return false
}

// IsSalaryValid checks for correct ranges or fixed values to be correct.
func (j Job) IsSalaryValid() bool {
	if j.SalaryMin > 0 {
		return j.SalaryMax > j.SalaryMin
	}
	return j.SalaryMax > 0
}

// Filter contains value patterns to match when searching for matching jobs.
type Filter struct {
	Text             string
	Location         string
	Salary           int
	Type             string
	IsRemoteFriendly *bool
	Keywords         []string
}

// Subscription contains the contact details to notify a person about a matching
// job.
type Subscription struct {
	Filter
	Email string
}