package subscription

import (
	"github.com/ariel17/jobberwocky/internal/adapters/repositories/keyword"
	"github.com/ariel17/jobberwocky/internal/core/domain"
)

type Subscription struct {
	ID               int64  `gorm:"primaryKey"`
	Email            string `gorm:"size:50;not null;uniqueIndex"`
	Text             string `gorm:"size:20"`
	Company          string `gorm:"size:50"`
	Location         string `gorm:"size:20"`
	Salary           int
	Type             string `gorm:"size:10"`
	IsRemoteFriendly *bool
	Keywords         []keyword.Keyword `gorm:"many2many:subscriptions_keywords;"`
}

func subscriptionDomainToModel(s domain.Subscription) Subscription {
	return Subscription{
		Email:            s.Email,
		Text:             s.Text,
		Company:          s.Company,
		Location:         s.Location,
		Salary:           s.Salary,
		Type:             s.Type,
		IsRemoteFriendly: s.IsRemoteFriendly,
		Keywords:         keyword.StringKeywordsToModel(s.Keywords),
	}
}

func subscriptionModelToDomain(s Subscription) domain.Subscription {
	return domain.Subscription{
		Pattern: domain.Pattern{
			Text:             s.Text,
			Company:          s.Company,
			Location:         s.Location,
			Salary:           s.Salary,
			Type:             s.Type,
			IsRemoteFriendly: s.IsRemoteFriendly,
			Keywords:         keyword.ModelKeywordsToString(s.Keywords),
		},
		Email: s.Email,
	}
}