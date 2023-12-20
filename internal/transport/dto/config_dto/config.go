package config_dto

import "keyboard-analysis/internal/config"

type ConfigDto struct {
	PasswordsCount    uint `json:"passwords_count"`
	PasswordMinLength uint `json:"password_min_length"`
	PasswordMaxLength uint `json:"password_max_length"`
}

func FromConfig(conf config.Config) ConfigDto {
	return ConfigDto{
		PasswordsCount:    conf.PasswordsCount,
		PasswordMinLength: conf.PasswordMinLength,
		PasswordMaxLength: conf.PasswordMaxLength,
	}
}
