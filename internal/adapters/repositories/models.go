package repositories

type Keyword struct {
	ID    int64  `gorm:"primaryKey"`
	Value string `gorm:"size:10;uniqueIndex"`
}