package wavix

import "github.com/wavix/sdk-go/utils"

type VoiceCampaignServiceInterface interface {
	TriggerScenario(payload TriggerScenarioPayload) (*TriggerScenarioResponse, *utils.HttpErrorResponse)
}

type VoiceCampaignService struct {
	httpConfig *utils.HttpConfig
}

type VoiceCampaignPayload struct {
	CallflowId int    `validate:"required" json:"callflow_id,omitempty"`
	CallerId   string `validate:"required" json:"caller_id,omitempty"`
	Contact    string `validate:"required" json:"contact,omitempty"`
}

type VoiceCampaignItem struct {
	Id        int    `json:"id"`
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
	CallerId  string `json:"caller_id"`
	Contact   string `json:"contact"`
}

type TriggerScenarioPayload struct {
	VoiceCampaign VoiceCampaignPayload `json:"voice_campaign"`
}

type TriggerScenarioResponse struct {
	VoiceCampaign VoiceCampaignItem `json:"voice_campaign"`
}

func (s *VoiceCampaignService) TriggerScenario(payload TriggerScenarioPayload) (*TriggerScenarioResponse, *utils.HttpErrorResponse) {
	validate := utils.GetValidate()
	err := validate.Struct(payload)

	if err != nil {
		return nil, &utils.HttpErrorResponse{Message: err.Error()}
	}

	return utils.Post[TriggerScenarioResponse](*s.httpConfig, "/v1/voice_campaigns", payload, TriggerScenarioResponse{})
}
