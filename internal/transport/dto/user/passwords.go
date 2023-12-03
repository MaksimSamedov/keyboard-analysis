package user

import "keyboard-analysis/internal/models"

type PasswordsDto []string

func PasswordsToDto(usr *models.User) PasswordsDto {
	var res PasswordsDto
	for _, pw := range usr.Passwords {
		res = append(res, pw.Password)
	}
	return res
}
