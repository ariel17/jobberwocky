package job

import (
	"strings"

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
	var (
		modelJobs []Job
		tx        *gorm.DB
	)
	if pattern != nil && !pattern.IsEmpty() {
		query, parameters := patternToQueryAndParameters(*pattern)
		tx = j.db.Where(query, parameters).Find(&modelJobs)
	} else {
		tx = j.db.Find(&modelJobs)
	}
	if tx.Error != nil {
		return nil, tx.Error
	}

	domainJobs := []domain.Job{}
	for _, job := range modelJobs {
		j, err := jobModelToDomain(job)
		if err != nil {
			return nil, err
		}
		domainJobs = append(domainJobs, j)
	}

	return domainJobs, nil
}

func (j *jobRepository) Save(job domain.Job) error {
	jm := jobDomainToModel(job)
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

func patternToQueryAndParameters(pattern domain.Pattern) (string, map[string]interface{}) {
	query := []string{}
	parameters := make(map[string]interface{})
	if pattern.Text != "" {
		query = append(query, "(title LIKE @text OR description LIKE @text)")
		parameters["text"] = `%` + pattern.Text + `%`
	}
	if pattern.Company != "" {
		query = append(query, "company = @company")
		parameters["company"] = pattern.Company
	}
	if pattern.Location != "" {
		query = append(query, "location = @location")
		parameters["location"] = pattern.Location
	}
	if pattern.Salary > 0 {
		query = append(query, "((salary_min <> 0 AND salary_min <= @salary AND salary_max >= @salary) OR (salary_min = 0 AND salary_max = @salary))")
		parameters["salary"] = pattern.Salary
	}
	if pattern.Type != "" {
		query = append(query, "type = @type")
		parameters["type"] = pattern.Type
	}
	if pattern.IsRemoteFriendly != nil {
		query = append(query, "is_remote_friendly = @is_remote_friendly")
		parameters["is_remote_friendly"] = *pattern.IsRemoteFriendly
	}
	// TODO keywords
	return strings.Join(query, " AND "), parameters
}