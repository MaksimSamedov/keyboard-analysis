package models

import (
	"gorm.io/gorm"
	"keyboard-analysis/internal/utils/password"
	"time"
)

type AccessToken struct {
	gorm.Model
	User       User
	UserID     uint      `gorm:"index:user_id"`
	Token      string    `gorm:"column:token;type:char(200);not null;<-:create"`
	Expiration time.Time `gorm:"column:expiration;not null;<-:create"`
}

func (token *AccessToken) IsExpired() bool {
	return time.Since(token.Expiration) > 0
}

func NewToken(usr User, length int, lifetime time.Duration) *AccessToken {
	gen := password.New(password.AllChars)
	return &AccessToken{
		//User:       usr,
		UserID:     usr.ID,
		Token:      gen.String(length),
		Expiration: time.Now().Add(lifetime),
	}
}
