package wavix

import (
	"path"

	"github.com/wavix/sdk-go/utils"
)

type TwoFaServiceInterface interface {
	GetServiceVerifications(serviceId string, queryParams GetServiceVerificationsQueryParams) (*[]TwoFaVerificationListItem, *utils.HttpErrorResponse)
	GetServiceVerificationEvents(sessionId string) (*[]TwoFaVerificationEventListItem, *utils.HttpErrorResponse)
	CreateVerification(payload CreateTwoFaVerificationPayload) (*CreateTwoFaVerificationResponse, *utils.HttpErrorResponse)
	ResendVerificationCode(sessionId string, payload ResendTwoFaVerificationCodePayload) (*ResendTwoFaVerificationCodeResponse, *utils.HttpErrorResponse)
	ValidateCode(sessionId string, payload ValidateTwoFaCodePayload) (*ValidateTwoFaCodeResponse, *utils.HttpErrorResponse)
	CancelVerification(sessionId string) (*utils.HttpSuccessBasicResponse, *utils.HttpErrorResponse)
}

type TwoFaService struct {
	httpConfig *utils.HttpConfig
}

type TwoFaChannelType string

const (
	SmsTwoFaChannelType   TwoFaChannelType = "sms"
	VoiceTwoFaChannelType TwoFaChannelType = "voice"
)

type TwoFaVerificationListItem struct {
	CreatedAt          string `json:"created_at"`
	SessionId          string `json:"session_id"`
	PhoneNumber        string `json:"phone_number"`
	DestinationCountry string `json:"destination_country"`
	Status             string `json:"status"`
	Charge             string `json:"charge"`
	ServiceId          string `json:"service_id"`
	ServiceName        string `json:"service_name"`
}

type TwoFaLookup struct {
	NumberType     string `json:"number_type"`
	Country        string `json:"country"`
	CurrentCarrier string `json:"current_carrier"`
}

type TwoFaVerificationEventListItem struct {
	CreatedAt string `json:"created_at"`
	Event     string `json:"event"`
	Status    string `json:"status"`
	Charge    string `json:"charge"`
	Error     string `json:"error"`
}

type GetServiceVerificationsQueryParams struct {
	utils.RequiredDateParams
}

type CreateTwoFaVerificationPayload struct {
	ServiceId string           `validate:"required" json:"service_id,omitempty"`
	To        string           `validate:"required" json:"to,omitempty"`
	Channel   TwoFaChannelType `validate:"required,oneof=sms voice" json:"channel,omitempty"`
}

type ResendTwoFaVerificationCodePayload struct {
	Channel TwoFaChannelType `validate:"required,oneof=sms voice" json:"channel,omitempty"`
}

type ValidateTwoFaCodePayload struct {
	Code string `validate:"required" json:"code,omitempty"`
}

type CreateTwoFaVerificationResponse struct {
	Success     bool        `json:"success"`
	ServiceId   string      `json:"service_id"`
	SessionUrl  string      `json:"session_url"`
	SessionId   string      `json:"session_id"`
	Destination string      `json:"destination"`
	CreatedAt   string      `json:"created_at"`
	Lookup      TwoFaLookup `json:"lookup"`
	Charge      string      `json:"charge"`
}

type ValidateTwoFaCodeResponse struct {
	IsValid bool `json:"is_valid"`
}

type ResendTwoFaVerificationCodeResponse struct {
	Success     bool             `json:"success"`
	Channel     TwoFaChannelType `json:"channel"`
	Destination string           `json:"destination"`
	CreatedAt   string           `json:"created_at"`
}

func (s *TwoFaService) GetServiceVerifications(serviceId string, queryParams GetServiceVerificationsQueryParams) (*[]TwoFaVerificationListItem, *utils.HttpErrorResponse) {
	basePath := path.Join("/v1/two-fa/service", serviceId, "sessions")
	url := utils.BuildUrlWithQueryString(basePath, queryParams)

	return utils.Get[[]TwoFaVerificationListItem](*s.httpConfig, url, []TwoFaVerificationListItem{})
}

func (s *TwoFaService) GetServiceVerificationEvents(sessionId string) (*[]TwoFaVerificationEventListItem, *utils.HttpErrorResponse) {
	url := path.Join("/v1/two-fa/session", sessionId, "events")

	return utils.Get[[]TwoFaVerificationEventListItem](*s.httpConfig, url, []TwoFaVerificationEventListItem{})
}

func (s *TwoFaService) CreateVerification(payload CreateTwoFaVerificationPayload) (*CreateTwoFaVerificationResponse, *utils.HttpErrorResponse) {
	validate := utils.GetValidate()
	err := validate.Struct(payload)

	if err != nil {
		return nil, &utils.HttpErrorResponse{Message: err.Error()}
	}

	return utils.Post[CreateTwoFaVerificationResponse](*s.httpConfig, "/v1/two-fa/verification", payload, CreateTwoFaVerificationResponse{})
}

func (s *TwoFaService) ResendVerificationCode(sessionId string, payload ResendTwoFaVerificationCodePayload) (*ResendTwoFaVerificationCodeResponse, *utils.HttpErrorResponse) {
	validate := utils.GetValidate()
	err := validate.Struct(payload)

	if err != nil {
		return nil, &utils.HttpErrorResponse{Message: err.Error()}
	}

	url := path.Join("/v1/two-fa/verification", sessionId)
	return utils.Post[ResendTwoFaVerificationCodeResponse](*s.httpConfig, url, payload, ResendTwoFaVerificationCodeResponse{})
}

func (s *TwoFaService) ValidateCode(sessionId string, payload ValidateTwoFaCodePayload) (*ValidateTwoFaCodeResponse, *utils.HttpErrorResponse) {
	validate := utils.GetValidate()
	err := validate.Struct(payload)

	if err != nil {
		return nil, &utils.HttpErrorResponse{Message: err.Error()}
	}

	url := path.Join("/v1/two-fa/verification", sessionId, "check")
	return utils.Post[ValidateTwoFaCodeResponse](*s.httpConfig, url, payload, ValidateTwoFaCodeResponse{})
}

func (s *TwoFaService) CancelVerification(sessionId string) (*utils.HttpSuccessBasicResponse, *utils.HttpErrorResponse) {
	url := path.Join("/v1/two-fa/verification", sessionId, "cancel")
	return utils.Patch(*s.httpConfig, url, nil, utils.HttpSuccessBasicResponse{})
}
