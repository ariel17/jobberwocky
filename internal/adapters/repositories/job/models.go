package job

import (
	"github.com/ariel17/jobberwocky/internal/adapters/repositories"
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
	Keywords         []repositories.Keyword `gorm:"many2many:jobs_keywords;"`
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
		Keywords:         repositories.StringKeywordsToModel(job.Keywords),
	}
}

func jobModelToDomain(job Job) (domain.Job, error) {
	keywords := []string{}
	for _, k := range job.Keywords {
		keywords = append(keywords, k.Value)
	}
	return domain.NewJob(job.Title, job.Description, job.Company, job.Location, job.SalaryMin, job.SalaryMax, job.Type, job.IsRemoteFriendly, "", keywords...)
}