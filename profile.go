package wavix

import (
	"github.com/wavix/sdk-go/utils"
)

type ProfileServiceInterface interface {
	GetCustomerInfo() (*GetCustomerInfoResponse, *utils.HttpErrorResponse)
	UpdateCustomerInfo(payload UpdateCustomerInfoPayload) (*UpdateCustomerInfoResponse, *utils.HttpErrorResponse)
	GetAccountSettings() (*GetAccountSettingsResponse, *utils.HttpErrorResponse)
}

type ProfileService struct {
	httpConfig *utils.HttpConfig
}

type DefaultDestination struct {
	Transport string `json:"transport"`
	Value     string `json:"value"`
}

type AccountSettingsGlobalLimits struct {
	MaxCallDuration int    `json:"max_call_duration"`
	MaxSipChannels  int    `json:"max_sip_channels"`
	MaxCallRate     string `json:"max_call_rate"`
}

type GetCustomerInfoResponse struct {
	Id                  int                  `json:"id"`
	AdditionalInfo      string               `json:"additional_info"`
	AttnContactName     string               `json:"attn_contact_name"`
	BillingAddress      string               `json:"billing_address"`
	CompanyName         string               `json:"company_name"`
	ContactEmail        string               `json:"contact_email"`
	DefaultDestinations []DefaultDestination `json:"default_destinations"`
	Email               string               `json:"email"`
	FirstName           string               `json:"first_name"`
	LastName            string               `json:"last_name"`
	Phone               string               `json:"phone"`
	Timezone            string               `json:"timezone"`
}

type GetAccountSettingsResponse struct {
	Balance      string                      `json:"balance"`
	GlobalLimits AccountSettingsGlobalLimits `json:"global_limits"`
}

type UpdateCustomerInfoResponse = GetCustomerInfoResponse

type UpdateCustomerInfoPayload struct {
	AdditionalInfo      string               `json:"additional_info,omitempty"`
	AttnContactName     string               `json:"attn_contact_name,omitempty"`
	BillingAddress      string               `json:"billing_address,omitempty"`
	CompanyName         string               `json:"company_name,omitempty"`
	ContactEmail        string               `validate:"omitempty,email" json:"contact_email,omitempty"`
	FirstName           string               `json:"first_name,omitempty"`
	LastName            string               `json:"last_name,omitempty"`
	Phone               string               `json:"phone,omitempty"`
	Timezone            string               `validate:"omitempty,timezone" json:"timezone,omitempty"`
	DefaultDestinations []DefaultDestination `json:"default_destinations,omitempty"`
}

func (s *ProfileService) GetCustomerInfo() (*GetCustomerInfoResponse, *utils.HttpErrorResponse) {
	return utils.Get[GetCustomerInfoResponse](*s.httpConfig, "/v1/profile", GetCustomerInfoResponse{})
}

func (s *ProfileService) UpdateCustomerInfo(payload UpdateCustomerInfoPayload) (*UpdateCustomerInfoResponse, *utils.HttpErrorResponse) {
	validate := utils.GetValidate()
	err := validate.Struct(payload)

	if err != nil {
		return nil, &utils.HttpErrorResponse{Message: err.Error()}
	}

	return utils.Put[UpdateCustomerInfoResponse](*s.httpConfig, "/v1/profile", payload, UpdateCustomerInfoResponse{})
}

func (s *ProfileService) GetAccountSettings() (*GetAccountSettingsResponse, *utils.HttpErrorResponse) {
	return utils.Get[GetAccountSettingsResponse](*s.httpConfig, "/v1/profile/config", GetAccountSettingsResponse{})
}
