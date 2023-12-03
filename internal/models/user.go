package models

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Login      string `gorm:"column:login;type:char(31);not null;<-:create"`
	Password   string `gorm:"column:password;type:char(63);not null;<-:create"`
	Passwords  []*Password
	SecretNote *string `gorm:"column:secret_note;type:text"`
}

var ErrUnableToHashPassword = errors.New("unable to get password hash")

func NewUser(login, password string, passwords []*Password) (*User, error) {
	pwHash, err := HashPassword(password)
	if err != nil {
		return nil, ErrUnableToHashPassword
	}

	return &User{
		Login:     login,
		Password:  pwHash,
		Passwords: passwords,
	}, nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (user *User) Attempt(password string) bool {
	return CheckPasswordHash(password, user.Password)
}

func (user *User) HasSecretNote() bool {
	return user.SecretNote != nil
}

func (user *User) SetSecretNote(note string) {
	user.SecretNote = &note
}
