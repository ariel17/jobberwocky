package job

import (
	"fmt"
	"strings"

	"gorm.io/gorm"

	"github.com/ariel17/jobberwocky/internal/adapters/repositories/keyword"
	"github.com/ariel17/jobberwocky/internal/core/domain"
	"github.com/ariel17/jobberwocky/internal/core/ports"
)

type jobRepository struct {
	db *gorm.DB
}

func NewJobRepository(db *gorm.DB) ports.JobRepository {
	return &jobRepository{
		db: db,
	}
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
		tx = j.db.Preload("Keywords").Find(&modelJobs)
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
	if tx.Error != nil {
		return tx.Error
	}
	needsReplacement, newKeywords, err := keyword.ReuseExistingKeywords(j.db, jm.Keywords)
	if err != nil {
		return err
	}
	if needsReplacement {
		return j.db.Model(&jm).Association("Keywords").Replace(newKeywords)
	}
	return nil
}

func (j *jobRepository) SyncSchemas() error {
	if err := j.db.AutoMigrate(&Job{}); err != nil {
		return err
	}
	if err := j.db.AutoMigrate(&keyword.Keyword{}); err != nil {
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
	if len(pattern.Keywords) > 0 {
		sub := "SELECT jk.job_id AS job_id, COUNT(*) FROM jobs_keywords jk INNER JOIN keywords k ON (jk.keyword_id=k.id) WHERE k.value IN @keywords GROUP BY jk.job_id HAVING COUNT(*) = @count"
		query = append(query, fmt.Sprintf("id IN (SELECT job_id FROM (%s))", sub))
		parameters["count"] = len(pattern.Keywords)
		parameters["keywords"] = pattern.Keywords
	}
	return strings.Join(query, " AND "), parameters
}