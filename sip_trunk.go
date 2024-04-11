package wavix

import (
	"fmt"
	"path"

	"github.com/wavix/sdk-go/utils"
)

type SipTrunkServiceInterface interface {
	GetAccountSipTrunks(params GetAccountSipTrunksQueryParams) (*GetAccountSipTrunksPaginatedResponse, *utils.HttpErrorResponse)
	GetSipTrunkConfiguration(sipTrunkId int) (*SipTrunkConfigurationItem, *utils.HttpErrorResponse)
	CreateSipTrunk(payload CreateSipTrunkPayload) (*SipTrunkConfigurationItem, *utils.HttpErrorResponse)
	UpdateSipTrunk(sipTrunkId int, payload UpdateSipTrunkPayload) (*SipTrunkConfigurationItem, *utils.HttpErrorResponse)
	DeleteSipTrunk(sipTrunkId int) (*utils.HttpSuccessBasicResponse, *utils.HttpErrorResponse)
}

type SipTrunkService struct {
	httpConfig *utils.HttpConfig
}

type GetAccountSipTrunksQueryParams struct {
	utils.PaginationParams
}

type GetAccountSipTrunksPaginatedResponse struct {
	Pagination utils.Pagination   `json:"pagination"`
	Items      []SipTrunkListItem `json:"sip_trunks"`
}

type HostRequest struct {
	Host   string `json:"host"`
	Status string `json:"status"`
}

type SipTrunkListItem struct {
	Id                      int          `json:"id"`
	TalkTime                int          `json:"talk_time"`
	TranscriptionThreshold  int          `json:"transcription_threshold"`
	HostRequest             *HostRequest `json:"host_request"`
	AuthMethod              string       `json:"auth_method"`
	CallerId                string       `json:"callerid"`
	Charge                  string       `json:"charge"`
	Label                   string       `json:"label"`
	Name                    string       `json:"name"`
	Status                  string       `json:"status"`
	TranscriptionEnabled    bool         `json:"transcription_enabled"`
	MultipleNumbers         bool         `json:"multiple_numbers"`
	Passthrough             bool         `json:"passthrough"`
	MachineDetectionEnabled bool         `json:"machine_detection_enabled"`
	CallRecordingEnabled    bool         `json:"call_recording_enabled"`
}

type AllowedIps struct {
	Id int    `json:"id"`
	Ip string `json:"ip"`
}

type SipTrunkConfigurationItem struct {
	Id                      int          `json:"id"`
	MaxChannels             int          `json:"max_channels"`
	CallLimit               int          `json:"call_limit"`
	TranscriptionThreshold  int          `json:"transcription_threshold"`
	CreatedAt               string       `json:"created_at"`
	Name                    string       `json:"string"`
	CallerId                string       `json:"callerid"`
	Label                   string       `json:"label"`
	AuthMethod              string       `json:"auth_method"`
	Host                    string       `json:"host"`
	RewritePrefix           string       `json:"rewrite_prefix"`
	RewriteCond             string       `json:"rewrite_cond"`
	MaxCallCost             string       `json:"max_call_cost"`
	AllowedIps              []AllowedIps `json:"allowed_ips"`
	CallRestrict            bool         `json:"call_restrict"`
	ChannelsRestrict        bool         `json:"channels_restrict"`
	IpRestrict              bool         `json:"ip_restrict"`
	CostLimit               bool         `json:"cost_limit"`
	RewriteEnabled          bool         `json:"rewrite_enabled"`
	CallRecordingEnabled    bool         `json:"call_recording_enabled"`
	MachineDetectionEnabled bool         `json:"machine_detection_enabled"`
	DidInfoEnabled          bool         `json:"didinfo_enabled"`
	TranscriptionEnabled    bool         `json:"transcription_enabled"`
}

type CreateSipTrunkPayload struct {
	Label                   string   `validate:"required" json:"label,omitempty"`
	Password                string   `validate:"required" json:"password"`
	CallerId                string   `validate:"required" json:"callerid"`
	MaxCallCost             string   `validate:"required" json:"max_call_cost"`
	Host                    string   `json:"host,omitempty"`
	RewritePrefix           string   `json:"rewrite_prefix,omitempty"`
	RewriteCond             string   `json:"rewrite_cond,omitempty"`
	AllowedIps              []string `json:"allowed_ips,omitempty"`
	MaxChannels             int      `json:"max_channels,omitempty"`
	CallLimit               int      `json:"call_limit,omitempty"`
	TranscriptionThreshold  int      `json:"transcription_threshold"`
	CostLimit               bool     `json:"cost_limit"`
	IpRestrict              bool     `json:"ip_restrict"`
	ChannelsRestrict        bool     `json:"channels_restrict"`
	CallRestrict            bool     `json:"call_restrict"`
	DidInfoEnabled          bool     `json:"didinfo_enabled"`
	TranscriptionEnabled    bool     `json:"transcription_enabled"`
	RewriteEnabled          bool     `json:"rewrite_enabled"`
	MachineDetectionEnabled bool     `json:"machine_detection_enabled,omitempty"`
	CallRecordingEnabled    bool     `json:"call_recording_enabled,omitempty"`
}
type UpdateSipTrunkPayload = CreateSipTrunkPayload

func (s *SipTrunkService) GetAccountSipTrunks(params GetAccountSipTrunksQueryParams) (*GetAccountSipTrunksPaginatedResponse, *utils.HttpErrorResponse) {
	return utils.Get[GetAccountSipTrunksPaginatedResponse](*s.httpConfig, "/v1/trunks", GetAccountSipTrunksPaginatedResponse{})
}

func (s *SipTrunkService) GetSipTrunkConfiguration(sipTrunkId int) (*SipTrunkConfigurationItem, *utils.HttpErrorResponse) {
	url := path.Join("/v1/trunks", fmt.Sprintf("%d", sipTrunkId))
	return utils.Get[SipTrunkConfigurationItem](*s.httpConfig, url, SipTrunkConfigurationItem{})
}

func (s *SipTrunkService) CreateSipTrunk(payload CreateSipTrunkPayload) (*SipTrunkConfigurationItem, *utils.HttpErrorResponse) {
	validate := utils.GetValidate()
	err := validate.Struct(payload)

	if err != nil {
		return nil, &utils.HttpErrorResponse{Message: err.Error()}
	}

	return utils.Post[SipTrunkConfigurationItem](*s.httpConfig, "/v1/trunks", payload, SipTrunkConfigurationItem{})
}

func (s *SipTrunkService) UpdateSipTrunk(sipTrunkId int, payload UpdateSipTrunkPayload) (*SipTrunkConfigurationItem, *utils.HttpErrorResponse) {
	validate := utils.GetValidate()
	err := validate.Struct(payload)

	if err != nil {
		return nil, &utils.HttpErrorResponse{Message: err.Error()}
	}

	url := path.Join("/v1/trunks", fmt.Sprintf("%d", sipTrunkId))
	return utils.Put[SipTrunkConfigurationItem](*s.httpConfig, url, payload, SipTrunkConfigurationItem{})
}

func (s *SipTrunkService) DeleteSipTrunk(sipTrunkId int) (*utils.HttpSuccessBasicResponse, *utils.HttpErrorResponse) {
	url := path.Join("/v1/trunks", fmt.Sprintf("%d", sipTrunkId))
	return utils.Delete[utils.HttpSuccessBasicResponse](*s.httpConfig, url, utils.HttpSuccessBasicResponse{Success: true})
}
