package wavix

import (
	"github.com/wavix/sdk-go/utils"
)

type E911ServiceInterface interface {
	GetList(params GetE911ListQueryParams) (*utils.PaginationResponse[E911ListItem], *utils.HttpErrorResponse)
	ValidateAddress(payload ValidateE911AddressPayload) (*ValidateE911AddressResponse, *utils.HttpErrorResponse)
	Create(payload CreateE911Payload) (*utils.HttpSuccessBasicResponse, *utils.HttpErrorResponse)
	Delete(params DeleteE911QueryParams) (*utils.HttpSuccessBasicResponse, *utils.HttpErrorResponse)
}

type E911Service struct {
	httpConfig *utils.HttpConfig
}

type GetE911ListQueryParams struct {
	utils.PaginationParams
	PhoneNumber string `url:"phone_number,omitempty"`
}

type DeleteE911QueryParams struct {
	PhoneNumber string `url:"phone_number"`
}

type E911Address struct {
	Location     string `validate:"required" json:"location,omitempty"`
	StreetNumber string `validate:"required" json:"street_number,omitempty"`
	Street       string `validate:"required" json:"street,omitempty"`
	City         string `validate:"required" json:"city,omitempty"`
	State        string `validate:"required" json:"state,omitempty"`
	ZipCode      string `validate:"required" json:"zip_code,omitempty"`
	ZipPlusFour  string `validate:"required" json:"zip_plus_four,omitempty"`
}

type CreateE911Payload struct {
	PhoneNumber string      `validate:"required" json:"phone_number,omitempty"`
	Name        string      `validate:"required" json:"name,omitempty"`
	Address     E911Address `validate:"required" json:"address,omitempty"`
}

type ValidateE911AddressPayload struct {
	PhoneNumber string      `validate:"required" json:"phone_number,omitempty"`
	Name        string      `validate:"required" json:"name,omitempty"`
	Address     E911Address `validate:"required" json:"address,omitempty"`
}

type E911ListItem struct {
	PhoneNumber string      `json:"phone_number"`
	Name        string      `json:"name"`
	Address     E911Address `json:"address"`
}

type ValidateE911AddressResponse struct {
	Status           int         `json:"status"`
	Number           string      `json:"number"`
	CorrectedAddress E911Address `json:"corrected_address"`
}

func (s *E911Service) GetList(params GetE911ListQueryParams) (*utils.PaginationResponse[E911ListItem], *utils.HttpErrorResponse) {
	url := utils.BuildUrlWithQueryString("/v1/e911-records", params)

	return utils.Get[utils.PaginationResponse[E911ListItem]](*s.httpConfig, url, utils.PaginationResponse[E911ListItem]{})
}

func (s *E911Service) ValidateAddress(payload ValidateE911AddressPayload) (*ValidateE911AddressResponse, *utils.HttpErrorResponse) {
	validate := utils.GetValidate()
	err := validate.Struct(payload)
	if err != nil {
		return nil, &utils.HttpErrorResponse{Message: err.Error()}
	}

	return utils.Post[ValidateE911AddressResponse](*s.httpConfig, "/v1/e911-records/validate-address", payload, ValidateE911AddressResponse{})
}

func (s *E911Service) Create(payload CreateE911Payload) (*utils.HttpSuccessBasicResponse, *utils.HttpErrorResponse) {
	validate := utils.GetValidate()
	err := validate.Struct(payload)
	if err != nil {
		return nil, &utils.HttpErrorResponse{Message: err.Error()}
	}

	return utils.Post[utils.HttpSuccessBasicResponse](*s.httpConfig, "/v1/e911-records", payload, utils.HttpSuccessBasicResponse{})
}

func (s *E911Service) Delete(params DeleteE911QueryParams) (*utils.HttpSuccessBasicResponse, *utils.HttpErrorResponse) {
	url := utils.BuildUrlWithQueryString("/v1/e911-records", params)

	return utils.Delete[utils.HttpSuccessBasicResponse](*s.httpConfig, url, utils.HttpSuccessBasicResponse{})
}
