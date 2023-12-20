package controllers

import (
	"github.com/gofiber/fiber/v2"
	"keyboard-analysis/internal/config"
	"keyboard-analysis/internal/transport/dto/config_dto"
)

type ConfigController struct {
	conf config.Config
}

func NewConfigController(conf config.Config) *ConfigController {
	return &ConfigController{conf: conf}
}

func (con *ConfigController) GetConfig(ctx *fiber.Ctx) error {
	return ctx.JSON(fiber.Map{
		"success": true,
		"config":  config_dto.FromConfig(con.conf),
	})
}
