package repositories

import "gorm.io/gorm"

func ReuseExistingKeywords(db *gorm.DB, oldKeywords []Keyword) (bool, []Keyword, error) {
	needsReplacement := false
	newKeywords := []Keyword{}
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