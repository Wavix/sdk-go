package wavix

import (
	"fmt"

	"github.com/wavix/sdk-go/utils"
)

type NumberValidationServiceInterface interface {
	ValidateSingle(number string, validationType string) (*NumberValidationBody, *utils.HttpErrorResponse)
	ValidateBatch(numbers []string, validationType string) (*NumberValidationResponse, *utils.HttpErrorResponse)
	ValidateBatchAsync(numbers []string, validationType string) (*NumberValidationAsyncResponse, *utils.HttpErrorResponse)
	GetValidationResult(uuid string) (*NumberValidationResponse, *utils.HttpErrorResponse)
}

type ValidationService struct {
	httpConfig *utils.HttpConfig
}

type NumberValidationBody struct {
	PhoneNumber      string `json:"phone_number"`
	Valid            bool   `json:"valid"`
	CountryCode      string `json:"country_code"`
	E164Format       string `json:"e164_format"`
	NationalFormat   string `json:"national_format"`
	Ported           bool   `json:"ported"`
	Mcc              string `json:"mcc"`
	Mnc              string `json:"mnc"`
	NumberType       string `json:"number_type"`
	CarrierName      string `json:"carrier_name"`
	RiskyDestination bool   `json:"risky_destination"`
	UnallocatedRange bool   `json:"unallocated_range"`
	Reachable        bool   `json:"reachable"`
	Roaming          bool   `json:"roaming"`
	Timezone         string `json:"timezone"`
	Charge           string `json:"charge"`
	ErrorCode        string `json:"error_code"`
}

type NumberValidationResponse struct {
	Status  string `json:"status"`
	Count   int    `json:"count"`
	Pending int    `json:"pending"`
	Items   []NumberValidationBody
}

type NumberValidationPayload struct {
	PhoneNumbers []string `json:"phone_numbers"`
	Async        bool     `json:"async"`
	Type         string   `json:"type"`
}

type NumberValidationAsyncResponse struct {
	RequestUUID string `json:"request_uuid"`
}

func (s *ValidationService) ValidateSingle(number string, validationType string) (*NumberValidationBody, *utils.HttpErrorResponse) {
	url := fmt.Sprintf("/v1/validation?phone_number=%s&type=%s", number, validationType)
	return utils.Get[NumberValidationBody](*s.httpConfig, url, NumberValidationBody{})
}

func (s *ValidationService) ValidateBatch(numbers []string, validationType string) (*NumberValidationResponse, *utils.HttpErrorResponse) {
	return utils.Post[NumberValidationResponse](*s.httpConfig, "/v1/validation", &NumberValidationPayload{
		PhoneNumbers: numbers,
		Type:         validationType,
		Async:        false,
	}, NumberValidationResponse{})
}

func (s *ValidationService) ValidateBatchAsync(numbers []string, validationType string) (*NumberValidationAsyncResponse, *utils.HttpErrorResponse) {
	return utils.Post[NumberValidationAsyncResponse](*s.httpConfig, "/v1/validation", &NumberValidationPayload{
		PhoneNumbers: numbers,
		Type:         validationType,
		Async:        true,
	}, NumberValidationAsyncResponse{})
}

func (s *ValidationService) GetValidationResult(uuid string) (*NumberValidationResponse, *utils.HttpErrorResponse) {
	url := fmt.Sprintf("/v1/validation/%s", uuid)
	return utils.Get[NumberValidationResponse](*s.httpConfig, url, NumberValidationResponse{})
}
