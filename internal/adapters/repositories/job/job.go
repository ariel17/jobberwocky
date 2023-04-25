package job

import (
	"gorm.io/gorm"

	"github.com/ariel17/jobberwocky/internal/core/domain"
	"github.com/ariel17/jobberwocky/internal/core/ports"
)

type jobRepository struct {
	db *gorm.DB
}

func NewJobRepository(dialector gorm.Dialector) (ports.JobRepository, error) {
	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return &jobRepository{
		db: db,
	}, nil
}

func (j *jobRepository) Filter(pattern *domain.Pattern) ([]domain.Job, error) {
	return nil, nil
}

func (j *jobRepository) Save(job domain.Job) error {
	jm := jobDomainToModels(job)
	tx := j.db.Create(&jm)
	return tx.Error
}

func (j *jobRepository) SyncSchemas() error {
	if err := j.db.AutoMigrate(&Job{}); err != nil {
		return err
	}
	if err := j.db.AutoMigrate(&Keyword{}); err != nil {
		return err
	}
	return nil
}