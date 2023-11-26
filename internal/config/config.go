package config

import "github.com/gofiber/fiber/v2"

type Config struct {
	Fiber      fiber.Config
	AppUrl     string
	DbUser     string
	DbPassword string
	DbHost     string
	DbPort     string
	DbName     string
}

func New() Config {
	return Config{}
}

func WithDefaults() Config {
	return Config{
		Fiber:      fiber.Config{},
		AppUrl:     "localhost:8000",
		DbUser:     "root",
		DbPassword: "",
		DbHost:     "localhost",
		DbPort:     "3306",
		DbName:     "bmil_lab2",
	}
}
