package auth

import (
	"errors"
	"strings"
)

type UserCredentials struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

var ErrShortLogin = errors.New("login is too short")
var ErrShortPassword = errors.New("password is too short")
var ErrSameLoginPassword = errors.New("login and password are same")

func (cred UserCredentials) Validate() error {
	cred.Login = strings.Trim(cred.Login, " ")
	cred.Password = strings.Trim(cred.Password, " ")

	if len(cred.Login) < 3 {
		return ErrShortLogin
	}
	if len(cred.Password) < 8 {
		return ErrShortPassword
	}
	if cred.Login == cred.Password {
		return ErrSameLoginPassword
	}
	return nil
}
