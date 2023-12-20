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
var ErrAlreadyHasPasswords = errors.New("passwords are already set")
var ErrNoPasswordsProvided = errors.New("no passwords provided")

func NewUser(login, password string) (*User, error) {
	pwHash, err := HashPassword(password)
	if err != nil {
		return nil, ErrUnableToHashPassword
	}

	return &User{
		Login:    login,
		Password: pwHash,
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

func (user *User) SetPasswords(passwords []*Password) error {
	if user.HasPasswords() {
		return ErrAlreadyHasPasswords
	}
	if len(passwords) == 0 {
		return ErrNoPasswordsProvided
	}
	user.Passwords = passwords
	return nil
}

func (user *User) HasPasswords() bool {
	return len(user.Passwords) != 0
}
