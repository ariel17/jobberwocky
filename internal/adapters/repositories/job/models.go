package job

import (
	"github.com/ariel17/jobberwocky/internal/core/domain"
)

type Job struct {
	ID               int64  `gorm:"primaryKey"`
	Title            string `gorm:"size:50;not null"`
	Description      string `gorm:"size:255"`
	Company          string `gorm:"size:50;not null"`
	Location         string `gorm:"size:50"`
	SalaryMin        int
	SalaryMax        int
	Type             string `gorm:"size:10"`
	IsRemoteFriendly *bool
	Keywords         []Keyword `gorm:"many2many:job_keyword;"`
}

type Keyword struct {
	ID    int64  `gorm:"primaryKey"`
	Value string `gorm:"size:10"`
}

func jobDomainToModels(job domain.Job) Job {
	km := []Keyword{}
	for _, k := range job.Keywords {
		km = append(km, Keyword{Value: k})
	}
	jm := Job{
		Title:            job.Title,
		Description:      job.Description,
		Company:          job.Company,
		Location:         job.Location,
		SalaryMin:        job.SalaryMin,
		SalaryMax:        job.SalaryMax,
		Type:             job.Type,
		IsRemoteFriendly: job.IsRemoteFriendly,
		Keywords:         km,
	}
	return jm
}