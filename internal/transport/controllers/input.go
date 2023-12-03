package controllers

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"keyboard-analysis/internal/services"
	"keyboard-analysis/internal/transport/dto/auth"
	"keyboard-analysis/internal/transport/dto/process"
	"keyboard-analysis/internal/transport/dto/user"
	"strconv"
)

type InputController struct {
	service     *services.KeyboardService
	userService *services.UserService
}

func NewInputController(service *services.KeyboardService, userService *services.UserService) *InputController {
	return &InputController{
		service:     service,
		userService: userService,
	}
}

func (con *InputController) GetPasswords(ctx *fiber.Ctx) error {
	var creds *auth.UserCredentials
	if err := json.Unmarshal(ctx.Request().Body(), &creds); err != nil {
		return ctx.JSON(fiber.Map{
			"success": false,
			"error":   ErrInvalidJson.Error(),
		})
	}

	usr, err := con.userService.RetrieveByCredentials(creds)
	if err != nil {
		return ctx.JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	res := user.PasswordsToDto(usr)
	return ctx.JSON(fiber.Map{
		"success": true,
		"data":    res,
	})
}

func (con *InputController) Process(ctx *fiber.Ctx) error {
	dto, err := process.FromBytes(ctx.Request().Body())
	if err != nil {
		return ctx.JSON(fiber.Map{
			"success": false,
			"error":   ErrInvalidJson.Error(),
		})
	}

	if _, _, err := con.service.ProcessFlow(dto); err != nil {
		return ctx.JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return ctx.JSON(dto)
}

func (con *InputController) History(ctx *fiber.Ctx) error {
	var creds auth.UserCredentials
	if err := json.Unmarshal(ctx.Request().Body(), &creds); err != nil {
		return ctx.JSON(fiber.Map{
			"success": false,
			"error":   ErrInvalidJson,
		})
	}

	usr, err := con.userService.RetrieveByCredentials(&creds)
	if err != nil {
		return ctx.JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	id, err := strconv.ParseUint(ctx.Params("id", ""), 10, 32)
	if err != nil {
		id = 0
	}
	if id != 0 {
		history, err := con.service.SingleHistory(usr, uint(id))
		if err != nil || history == nil {
			return ctx.JSON(fiber.Map{
				"success": false,
				"error":   "Not found",
			})
		}
		return ctx.JSON(process.FromModel(*history))
	} else {
		history, err := con.service.History(usr)
		if err != nil {
			return ctx.JSON(fiber.Map{
				"success": false,
				"error":   "Error while getting history",
			})
		}
		return ctx.JSON(process.FromModels(history))
	}
}

func (con *InputController) GetToken(ctx *fiber.Ctx) error {
	dto, err := process.FromBytes(ctx.Request().Body())
	if err != nil {
		return ctx.JSON(fiber.Map{
			"success": false,
			"error":   ErrInvalidJson.Error(),
		})
	}

	token, err := con.service.GetToken(dto)
	if err != nil {
		return ctx.JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return ctx.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"token":      token.Token,
			"expiration": token.Expiration.Unix(),
		},
	})
}
