package secret

import (
	"encoding/json"
	"errors"
	"github.com/gofiber/fiber/v2"
)

type RWSecretDto struct {
	Token string  `json:"token"`
	Value *string `json:"value,omitempty"`
}

var ErrInvalidJson = errors.New("invalid json")
var ErrEmptyToken = errors.New("token is empty")

func FromContext(ctx *fiber.Ctx) (*RWSecretDto, error) {
	var dto *RWSecretDto
	if err := json.Unmarshal(ctx.Request().Body(), &dto); err != nil {
		return nil, ErrInvalidJson
	}
	if len(dto.Token) == 0 {
		return nil, ErrEmptyToken
	}
	return dto, nil
}

func (dto *RWSecretDto) HasValue() bool {
	return dto.Value != nil
}
