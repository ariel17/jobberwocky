package keyword

type Keyword struct {
	ID    int64  `gorm:"primaryKey"`
	Value string `gorm:"size:10;not null;uniqueIndex"`
}

func StringKeywordsToModel(keywords []string) []Keyword {
	km := []Keyword{}
	for _, k := range keywords {
		km = append(km, Keyword{Value: k})
	}
	return km
}

func ModelKeywordsToString(keywords []Keyword) []string {
	newKeywords := []string{}
	for _, k := range keywords {
		newKeywords = append(newKeywords, k.Value)
	}
	return newKeywords
}