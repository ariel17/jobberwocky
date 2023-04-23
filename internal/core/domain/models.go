package domain

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

// Subscription contains the contact details to notify a person about a matching
// job.
type Subscription struct {
	Email            string
	Text             string
	Location         string
	Salary           int
	Type             string
	IsRemoteFriendly bool
	Keywords         []string
}