package process

import (
	"encoding/json"
	"keyboard-analysis/internal/models"
	"keyboard-analysis/internal/transport/dto/auth"
)

type KeyboardFlowResponse struct {
	ID     uint                    `json:"id"`
	Phrase string                  `json:"phrase"`
	Flow   []*models.KeyboardEvent `json:"flow"`
}

type KeyboardFlowResults struct {
	Auth  auth.UserCredentials `json:"auth"`
	Flows []KeyboardFlowResult `json:"flows"`
}

type KeyboardFlowResult struct {
	Flow   []*models.KeyboardEvent `json:"flow"`
	Phrase string                  `json:"phrase"`
}

func FromBytes(data []byte) (*KeyboardFlowResults, error) {
	var dto *KeyboardFlowResults
	if err := json.Unmarshal(data, &dto); err != nil {
		return nil, err
	}
	return dto, nil
}

func FromModel(flow models.KeyboardFlow) KeyboardFlowResponse {
	return KeyboardFlowResponse{
		ID:     flow.ID,
		Phrase: flow.Password.Password,
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
