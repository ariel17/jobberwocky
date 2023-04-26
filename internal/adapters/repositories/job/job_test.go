package job

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/ariel17/jobberwocky/internal/core/domain"
	helpers "github.com/ariel17/jobberwocky/internal/internal_test"

	"gorm.io/driver/sqlite"
)

const (
	dbName = "test.db"
)

func TestJobRepository_Save(t *testing.T) {
	testCases := []struct {
		name string
		job  domain.Job
	}{
		{"success", domain.Job{"Title", "Description", "Company", "Argentina", 60, 80, domain.FullTime, helpers.BoolPointer(true), []string{"k1", "k2", "k3"}, ""}},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_ = os.Remove(dbName)
			db, err := gorm.Open(sqlite.Open(dbName), &gorm.Config{
				Logger: logger.Default.LogMode(logger.Info),
			})
			assert.Nil(t, err)

			repository, err := NewJobRepository(db)
			assert.Nil(t, err)

			err = repository.SyncSchemas()
			assert.Nil(t, err)

			err = repository.Save(tc.job)
			assert.Nil(t, err)

			var job Job
			db.First(&job, 1)
			assert.Equal(t, job.Title, tc.job.Title)
			assert.Equal(t, job.Description, tc.job.Description)
			assert.Equal(t, job.Company, tc.job.Company)
			assert.Equal(t, job.Location, tc.job.Location)
			assert.Equal(t, job.SalaryMin, tc.job.SalaryMin)
			assert.Equal(t, job.SalaryMax, tc.job.SalaryMax)
			assert.Equal(t, job.Type, tc.job.Type)
			assert.Equal(t, *job.IsRemoteFriendly, *tc.job.IsRemoteFriendly)
			for _, k1 := range job.Keywords {
				found := false
				for _, k2 := range tc.job.Keywords {
					if k1.Value == k2 {
						found = true
						break
					}
				}
				assert.True(t, found)
			}
		})
	}
}

func TestJobRepository_Filter(t *testing.T) {
	jobs := []domain.Job{
		{"Title", "Description", "Company1", "Argentina", 60, 80, domain.FullTime, helpers.BoolPointer(true), []string{"java", "python", "golang"}, ""},
		{"Another", "Value", "Company2", "USA", 0, 90, domain.Contractor, helpers.BoolPointer(false), []string{"java", "python", "kotlin"}, ""},
		{"X", "", "SpaceX", "USA", 900, 1000, domain.Contractor, helpers.BoolPointer(true), []string{"java", "python", "php"}, ""},
	}
	testCases := []struct {
		name     string
		pattern  *domain.Pattern
		expected int
	}{
		{"all jobs without filter", nil, 3},
		{"all jobs with empty filter", &domain.Pattern{}, 3},
		{"filter by text", &domain.Pattern{Text: "title"}, 1},
		{"filter by company", &domain.Pattern{Company: "Company1"}, 1},
		{"filter by location", &domain.Pattern{Location: "Argentina"}, 1},
		{"filter by fixed salary", &domain.Pattern{Salary: 90}, 1},
		{"filter by ranged salary", &domain.Pattern{Salary: 70}, 1},
		{"filter by type", &domain.Pattern{Type: domain.Contractor}, 2},
		{"filter by is remote friendly", &domain.Pattern{IsRemoteFriendly: helpers.BoolPointer(true)}, 2},
		{"filter by is remote friendly and company", &domain.Pattern{Company: "SpaceX", IsRemoteFriendly: helpers.BoolPointer(true)}, 1},
		{"not matching", &domain.Pattern{Location: "Argentina", Type: domain.Contractor}, 0},
		{"filter by single keyword", &domain.Pattern{Keywords: []string{"python"}}, 3},
		{"filter by multiple keywords", &domain.Pattern{Keywords: []string{"python", "golang"}}, 1},
		{"not matching by multiple keywords", &domain.Pattern{Keywords: []string{"python", "golang", "xxx"}}, 0},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_ = os.Remove(dbName)
			db, err := gorm.Open(sqlite.Open(dbName), &gorm.Config{
				Logger: logger.Default.LogMode(logger.Info),
			})
			assert.Nil(t, err)

			repository, err := NewJobRepository(db)
			assert.Nil(t, err)

			err = repository.SyncSchemas()
			assert.Nil(t, err)

			for _, job := range jobs {
				err = repository.Save(job)
				assert.Nil(t, err)
			}

			jobs, err := repository.Filter(tc.pattern)
			assert.Nil(t, err)
			assert.NotNil(t, jobs)
			assert.Equal(t, tc.expected, len(jobs))
		})
	}
}