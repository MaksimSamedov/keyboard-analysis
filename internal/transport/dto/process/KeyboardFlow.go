package process

import (
	"encoding/json"
	"keyboard-analysis/internal/models"
)

type KeyboardFlowResponse struct {
	ID     uint                   `json:"id"`
	Phrase string                 `json:"phrase"`
	Flow   []models.KeyboardEvent `json:"flow"`
}

func FromBytes(data []byte) ([]models.KeyboardFlow, error) {
	var dto []models.KeyboardFlow
	if err := json.Unmarshal(data, &dto); err != nil {
		return nil, err
	}
	return dto, nil
}

func FromModel(flow models.KeyboardFlow) KeyboardFlowResponse {
	return KeyboardFlowResponse{
		ID:     flow.ID,
		Phrase: flow.Phrase,
		Flow:   flow.Flow,
	}
}

func FromModels(flows []models.KeyboardFlow) []KeyboardFlowResponse {
	var dtos []KeyboardFlowResponse
	for _, flow := range flows {
		dtos = append(dtos, FromModel(flow))
	}
	return dtos
}
