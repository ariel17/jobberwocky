package keyword

import (
	"gorm.io/gorm"

	"github.com/ariel17/jobberwocky/internal/core/ports"
)

type keywordRepository struct {
	db *gorm.DB
}

func NewKeywordRepository(db *gorm.DB) ports.Repository {
	return &keywordRepository{
		db: db,
	}
}

func (k *keywordRepository) SyncSchemas() error {
	return k.db.AutoMigrate(&Keyword{})
}

func ReuseExistingKeywords(db *gorm.DB, oldKeywords []Keyword) (bool, []Keyword, error) {
	needsReplacement := false
	newKeywords := make([]Keyword, 0)
	for _, k := range oldKeywords {
		if k.ID != 0 {
			newKeywords = append(newKeywords, k)
			continue
		}
		needsReplacement = true
		var existing Keyword
		tx := db.Where("value = ?", k.Value).First(&existing)
		if tx.Error != nil {
			return false, nil, tx.Error
		}
		k.ID = existing.ID
		newKeywords = append(newKeywords, k)
	}
	return needsReplacement, newKeywords, nil
}