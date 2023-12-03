package controllers

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"keyboard-analysis/internal/services"
	"keyboard-analysis/internal/transport/dto/auth"
	"keyboard-analysis/internal/transport/dto/secret"
)

type UserController struct {
	service *services.UserService
}

func NewUserController(service *services.UserService) *UserController {
	return &UserController{service: service}
}

func (con *UserController) Register(ctx *fiber.Ctx) error {
	var dto auth.UserCredentials
	if err := json.Unmarshal(ctx.Request().Body(), &dto); err != nil {
		return ctx.JSON(fiber.Map{
			"success": false,
			"error":   ErrInvalidJson.Error(),
		})
	}

	_, err := con.service.Register(dto)
	if err != nil {
		return ctx.JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return ctx.JSON(fiber.Map{
		"success": true,
	})
}

func (con *UserController) Login(ctx *fiber.Ctx) error {
	var dto auth.UserCredentials
	if err := json.Unmarshal(ctx.Request().Body(), &dto); err != nil {
		return ctx.JSON(fiber.Map{
			"success": false,
			"error":   ErrInvalidJson.Error(),
		})
	}

	_, err := con.service.Login(dto)
	if err != nil {
		return ctx.JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return ctx.JSON(fiber.Map{
		"success": true,
	})
}

func (con *UserController) UserHasSecret(ctx *fiber.Ctx) error {
	var dto auth.UserCredentials
	if err := json.Unmarshal(ctx.Request().Body(), &dto); err != nil {
		return ctx.JSON(fiber.Map{
			"success": false,
			"error":   ErrInvalidJson.Error(),
		})
	}

	usr, err := con.service.Login(dto)
	if err != nil {
		return ctx.JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return ctx.JSON(fiber.Map{
		"success": true,
		"data":    usr.HasSecretNote(),
	})
}

func (con *UserController) GetSecret(ctx *fiber.Ctx) error {
	dto, err := secret.FromContext(ctx)
	if err != nil {
		return ctx.JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	usr, err := con.service.RetrieveByToken(dto.Token)
	if err != nil {
		return ctx.JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return ctx.JSON(fiber.Map{
		"success": true,
		"data":    usr.SecretNote,
	})
}

func (con *UserController) SetSecret(ctx *fiber.Ctx) error {
	dto, err := secret.FromContext(ctx)
	if err != nil {
		return ctx.JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	if !dto.HasValue() {
		return ctx.JSON(fiber.Map{
			"success": false,
			"error":   "Вы не указали value",
		})
	}

	usr, err := con.service.RetrieveByToken(dto.Token)
	if err != nil {
		return ctx.JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	if err := con.service.SetUserSecret(usr, *dto.Value); err != nil {
		return ctx.JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return ctx.JSON(fiber.Map{
		"success": true,
		"data":    usr.SecretNote,
	})
}
