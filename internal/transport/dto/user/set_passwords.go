package user

import "keyboard-analysis/internal/transport/dto/auth"

type SetPasswordsDto struct {
	Auth      auth.UserCredentials `json:"auth"`
	Passwords PasswordsDto         `json:"passwords"`
}
