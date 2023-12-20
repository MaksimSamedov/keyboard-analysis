package models

import (
	"gorm.io/gorm"
	"keyboard-analysis/internal/utils/password"
)

type Password struct {
	gorm.Model
	Password string `gorm:"column:password;type:char(127);not null;<-:create"`
	User     User
	UserID   uint `gorm:"index:user_id"`
}

var pwGen = password.New(password.LettersLowerCase + password.Numbers)

func GenerateForUser(length uint) *Password {
	return &Password{
		Password: pwGen.String(int(length)),
	}
}

func PasswordFromString(pw string) *Password {
	return &Password{Password: pw}
}
