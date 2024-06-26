package job

import (
	"github.com/ariel17/jobberwocky/internal/adapters/repositories/keyword"
	"github.com/ariel17/jobberwocky/internal/core/domain"
)

type Job struct {
	ID               int64  `gorm:"primaryKey"`
	Title            string `gorm:"size:50;not null"`
	Description      string `gorm:"size:255"`
	Company          string `gorm:"size:50"`
	Location         string `gorm:"size:50;not null"`
	SalaryMin        int
	SalaryMax        int
	Type             string `gorm:"size:10"`
	IsRemoteFriendly *bool
	Keywords         []keyword.Keyword `gorm:"many2many:jobs_keywords;"`
}

func jobDomainToModel(job domain.Job) Job {
	return Job{
		Title:            job.Title,
		Description:      job.Description,
		Company:          job.Company,
		Location:         job.Location,
		SalaryMin:        job.SalaryMin,
		SalaryMax:        job.SalaryMax,
		Type:             job.Type,
		IsRemoteFriendly: job.IsRemoteFriendly,
		Keywords:         keyword.StringKeywordsToModel(job.Keywords),
	}
}

func jobModelToDomain(job Job) (domain.Job, error) {
	keywords := keyword.ModelKeywordsToString(job.Keywords)
	j := domain.Job{Title: job.Title, Description: job.Description, Company: job.Company, Location: job.Location, SalaryMin: job.SalaryMin, SalaryMax: job.SalaryMax, Type: job.Type, IsRemoteFriendly: job.IsRemoteFriendly, Keywords: keywords}
	return j, j.Validate(true)
}