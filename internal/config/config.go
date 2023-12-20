package config

import (
	"github.com/gofiber/fiber/v2"
	"keyboard-analysis/internal/utils/flow"
	"time"
)

type Config struct {
	Fiber               fiber.Config
	AppUrl              string
	DbUser              string
	DbPassword          string
	DbHost              string
	DbPort              string
	DbName              string
	PasswordsCount      uint
	PasswordMinLength   uint
	PasswordMaxLength   uint
	TokenLength         int
	TokenLifetime       time.Duration
	MinSamples          int // Сколько раз пользователь должен ввести каждый пароль чтобы начать работать с секретом
	SamplesCompareCount int // Сколько последних паролей надо использовать при проверке сходства
	AnalyserProps       flow.AnalyserProps
}

func New() Config {
	return Config{}
}

func WithDefaults() Config {
	return Config{
		Fiber:               fiber.Config{},
		AppUrl:              "localhost:8000",
		DbUser:              "root",
		DbPassword:          "",
		DbHost:              "localhost",
		DbPort:              "3306",
		DbName:              "bmil_lab2",
		PasswordsCount:      10,
		PasswordMinLength:   10,
		PasswordMaxLength:   64,
		TokenLength:         200,
		TokenLifetime:       time.Hour,
		MinSamples:          5,
		SamplesCompareCount: 50,
		AnalyserProps: flow.AnalyserProps{
			MaxDeviation:             10,
			MinSuccessfulComparisons: 80,
			MaxErrorsInTask:          20,
		},
	}
}
