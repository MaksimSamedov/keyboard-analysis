package controllers

import (
	"github.com/gofiber/fiber/v2"
	"keyboard-analysis/internal/services"
	"keyboard-analysis/internal/transport/dto/process"
	"strconv"
)

type InputController struct {
	service *services.KeyboardService
}

func NewInputController(service *services.KeyboardService) *InputController {
	return &InputController{
		service: service,
	}
}

func (con *InputController) Process(ctx *fiber.Ctx) error {
	dto, err := process.FromBytes(ctx.Request().Body())
	if err != nil {
		return ctx.JSON(fiber.Map{
			"success": false,
			"error":   "Invalid data",
		})
	}

	if err := con.service.ProcessFlow(dto); err != nil {
		return ctx.JSON(fiber.Map{
			"success": false,
			"error":   "Error while saving",
		})
	}

	return ctx.JSON(dto)
}

func (con *InputController) History(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id", ""), 10, 32)
	if err != nil {
		id = 0
	}
	if id != 0 {
		history, err := con.service.SingleHistory(uint(id))
		if err != nil || history == nil {
			return ctx.JSON(fiber.Map{
				"success": false,
				"error":   "Not found",
			})
		}
		return ctx.JSON(process.FromModel(*history))
	} else {
		history, err := con.service.History()
		if err != nil {
			return ctx.JSON(fiber.Map{
				"success": false,
				"error":   "Error while getting history",
			})
		}
		return ctx.JSON(process.FromModels(history))
	}
}
