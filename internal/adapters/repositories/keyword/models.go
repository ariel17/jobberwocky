package keyword

type Keyword struct {
	ID    int64  `gorm:"primaryKey"`
	Value string `gorm:"size:10;not null;uniqueIndex"`
}

func StringKeywordsToModel(keywords []string) []Keyword {
	km := make([]Keyword, 0)
	for _, k := range keywords {
		km = append(km, Keyword{Value: k})
	}
	return km
}

func ModelKeywordsToString(keywords []Keyword) []string {
	newKeywords := make([]string, 0)
	for _, k := range keywords {
		newKeywords = append(newKeywords, k.Value)
	}
	return newKeywords
}