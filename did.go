package wavix

import (
	"bytes"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/wavix/sdk-go/utils"
)

type DidDestinationTransport int
type DidDocumentId int

const (
	DidDestinationTransportSipURI   DidDestinationTransport = 1
	DidDestinationTransportPSTN     DidDestinationTransport = 4
	DidDestinationTransportSIPTrunk DidDestinationTransport = 5
)

const (
	DidDocumentIdGeneral    DidDocumentId = 1
	DidDocumentAddress      DidDocumentId = 2
	DidDocumentLocalAddress DidDocumentId = 3
)

type DidServiceInterface interface {
	GetAccountDids(params GetAccountDidsQueryParams) (*utils.PaginationResponse[DidItem], *utils.HttpErrorResponse)
	UpdateDidDestinations(payload UpdateDidDestinationsPayload) (*utils.HttpSuccessBasicResponse, *utils.HttpErrorResponse)
	UploadDidDocument(payload UploadDidDocumentPayload) (*utils.HttpSuccessBasicResponse, *utils.HttpErrorResponse)
	ReturnDidsToStock(ids []string) (*utils.HttpSuccessBasicResponse, *utils.HttpErrorResponse)
}

type DidService struct {
	httpConfig *utils.HttpConfig
}

type DidDestination struct {
	Id          int    `json:"id"`
	Destination string `json:"destination"`
	Priority    int    `json:"priority"`
	Transport   int    `json:"transport"`
	TrunkId     int    `json:"trunk_id"`
	TrunkLabel  string `json:"trunk_label"`
}

type DidDocument struct {
	Id             int    `json:"id"`
	AllowReplace   bool   `json:"allow_replace"`
	DidNumber      string `json:"did_number"`
	DocContentType string `json:"doc_content_type"`
	DocFileName    string `json:"doc_file_name"`
	DocTypeId      int    `json:"doc_type_id"`
	Status         string `json:"status"`
	Url            string `json:"url"`
}

type DidItem struct {
	Id                     int              `json:"id"`
	ActivationFee          string           `json:"activation_fee"`
	Added                  string           `json:"added"`
	CallRecordingEnabled   bool             `json:"call_recording_enabled"`
	Channels               int              `json:"channels"`
	City                   string           `json:"city"`
	Cnam                   bool             `json:"cnam"`
	Country                string           `json:"country"`
	CountryShortName       string           `json:"country_short_name"`
	Destination            []DidDestination `json:"destination"`
	Documents              []DidDocument    `json:"documents"`
	Label                  string           `json:"label"`
	MonthlyFee             string           `json:"monthly_fee"`
	Number                 string           `json:"number"`
	PaidUntil              string           `json:"paid_until"`
	PerMin                 string           `json:"per_min"`
	RequireDocs            []string         `json:"require_docs"`
	Seconds                string           `json:"seconds"`
	SmsEnabled             bool             `json:"sms_enabled"`
	SmsRelayUrl            string           `json:"sms_relay_url"`
	Status                 string           `json:"status"`
	TranscriptionEnabled   bool             `json:"transcription_enabled"`
	TranscriptionThreshold int              `json:"transcription_threshold"`
}

type GetAccountDidsQueryParams struct {
	utils.PaginationParams
	CityId       int    `url:"city_id,omitempty"`
	Search       string `url:"search,omitempty"`
	Label        string `url:"label,omitempty"`
	LabelPresent string `url:"label_present,omitempty"`
}

type DidDestinationPayload struct {
	Destination string `validate:"required" json:"destination"`
	Transport   int    `validate:"required,oneof=1 4 5" json:"transport"`
	TrunkId     int    `validate:"required" json:"trunk_id"`
	Priority    int    `json:"priority,omitempty"`
}

type UpdateDidDestinationsPayload struct {
	Ids          []int                   `validate:"min=1" json:"ids"`
	SmsRelayUrl  string                  `validate:"omitempty,url" json:"sms_relay_url,omitempty"`
	Destinations []DidDestinationPayload `validate:"omitempty,min=1,dive" json:"destinations,omitempty"`
}

type DidDocumentFile struct {
	Data []byte `json:"data"`
	Name string `json:"name"`
}

type UploadDidDocumentPayload struct {
	DidIds []string        `validate:"min=1" json:"did_ids"`
	File   DidDocumentFile `validate:"required" json:"file"`
	DocId  int             `validate:"required,oneof=1 2 3" json:"doc_id"`
}

func (d UploadDidDocumentPayload) GetFileData() utils.File {
	return utils.File{Reader: bytes.NewReader(d.File.Data), FileName: d.File.Name, FileKey: "doc_attachment"}
}

func (d UploadDidDocumentPayload) GetFormValues() url.Values {
	return url.Values{
		"did_ids": []string{strings.Join(d.DidIds, ",")},
		"doc_id":  []string{strconv.Itoa(d.DocId)},
	}
}

func (s *DidService) GetAccountDids(params GetAccountDidsQueryParams) (*utils.PaginationResponse[DidItem], *utils.HttpErrorResponse) {
	url := utils.BuildUrlWithQueryString("/v1/mydids", params)

	return utils.Get[utils.PaginationResponse[DidItem]](*s.httpConfig, url, utils.PaginationResponse[DidItem]{})
}

func (s *DidService) UpdateDidDestinations(payload UpdateDidDestinationsPayload) (*utils.HttpSuccessBasicResponse, *utils.HttpErrorResponse) {
	validate := utils.GetValidate()
	err := validate.Struct(payload)

	if err != nil {
		return nil, &utils.HttpErrorResponse{Message: err.Error()}
	}

	return utils.Post[utils.HttpSuccessBasicResponse](*s.httpConfig, "/v1/mydids/update-destinations", payload, utils.HttpSuccessBasicResponse{})
}

func (s *DidService) UploadDidDocument(payload UploadDidDocumentPayload) (*utils.HttpSuccessBasicResponse, *utils.HttpErrorResponse) {
	validate := utils.GetValidate()
	err := validate.Struct(payload)

	if err != nil {
		return nil, &utils.HttpErrorResponse{Message: err.Error()}
	}

	return utils.Upload(*s.httpConfig, "/v1/mydids/papers", payload)
}

func (s *DidService) ReturnDidsToStock(ids []string) (*utils.HttpSuccessBasicResponse, *utils.HttpErrorResponse) {
	queryStringSlice := make([]string, len(ids))
	for index, id := range ids {
		queryStringSlice[index] = fmt.Sprintf("ids[]=%v", id)
	}
	url := "/v1/mydids?" + strings.Join(queryStringSlice, "&")

	return utils.Delete[utils.HttpSuccessBasicResponse](*s.httpConfig, url, utils.HttpSuccessBasicResponse{})
}
