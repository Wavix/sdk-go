package wavix

import (
	"path"

	"github.com/wavix/sdk-go/utils"
)

type SpeechAnalyticsServiceInterface interface {
	GetSpeechAnalyticsCalls(payload GetSpeechAnalyticsCallsPayload) (*utils.PaginationResponse[SpeechAnalyticsCallItem], *utils.HttpErrorResponse)
	TranscribeCallById(callId string, payload TranscribeCallByIdPayload) (*utils.HttpSuccessBasicResponse, *utils.HttpErrorResponse)
	RequestTranscriptionByCallId(callId string) (*RequestTranscriptionByCallIdResponse, *utils.HttpErrorResponse)
}

type SpeechAnalyticsLanguage string
type SpeechAnalyticsCallType string

const (
	EnglishSpeechAnalyticsLanguage SpeechAnalyticsLanguage = "en"
	GermanSpeechAnalyticsLanguage  SpeechAnalyticsLanguage = "de"
	SpanishSpeechAnalyticsLanguage SpeechAnalyticsLanguage = "es"
)

const (
	ReceivedSpeechAnalyticsCallType SpeechAnalyticsCallType = "received"
	PlacedSpeechAnalyticsCallType   SpeechAnalyticsCallType = "placed"
)

type SpeechAnalyticsService struct {
	httpConfig *utils.HttpConfig
}

type SpeechAnalyticsCallTranscription struct {
	Uuid string `json:"uuid"`
	Url  string `json:"url"`
}

type TranscriptionTurn struct {
	PhoneNumber string `json:"phone_number"`
	S           string `json:"s"`
	E           string `json:"e"`
	Text        string `json:"text"`
}

type SpeechAnalyticsCallItem struct {
	Charge        string                            `json:"charge"`
	Date          string                            `json:"date"`
	Destination   string                            `json:"destination"`
	Disposition   string                            `json:"disposition"`
	Duration      int                               `json:"duration"`
	ForwardFee    string                            `json:"forward_fee"`
	From          string                            `json:"from"`
	To            string                            `json:"to"`
	PerMinute     string                            `json:"per_minute"`
	Uuid          string                            `json:"uuid"`
	SipTrunk      string                            `json:"sip_trunk"`
	Transcription *SpeechAnalyticsCallTranscription `json:"transcription"`
}

type TranscriptionPayloadItem struct {
	Must    []string `json:"must"`
	Match   []string `json:"match"`
	Exclude []string `json:"exclude"`
}
type TranscriptionPayload struct {
	Agent  TranscriptionPayloadItem `json:"agent,omitempty"`
	Client TranscriptionPayloadItem `json:"client,omitempty"`
	Any    TranscriptionPayloadItem `json:"any,omitempty"`
}

type GetSpeechAnalyticsCallsPayload struct {
	utils.PaginationPayload
	utils.RequiredDatePayload
	Type          SpeechAnalyticsCallType `validate:"oneof=received placed" json:"type,omitempty"`
	FromSearch    string                  `json:"from_search,omitempty"`
	ToSearch      string                  `json:"to_search,omitempty"`
	SipTrunk      string                  `json:"sip_trunk,omitempty"`
	MinDuration   int                     `json:"min_duration,omitempty"`
	Transcription *TranscriptionPayload   `json:"transcription,omitempty"`
}

type TranscribeCallByIdPayload struct {
	Language   SpeechAnalyticsLanguage `validate:"required,oneof=en de es" json:"language,omitempty"`
	WebhookUrl string                  `validate:"required,url" json:"webhook_url,omitempty"`
}

type RequestTranscriptionByCallIdResponse struct {
	Transcription     map[string]interface{}  `json:"transcription"`
	Turns             []TranscriptionTurn     `json:"turns"`
	Uuid              string                  `json:"uuid"`
	Language          SpeechAnalyticsLanguage `json:"language"`
	Duration          int                     `json:"duration"`
	Charge            string                  `json:"charge"`
	Status            string                  `json:"status"`
	TranscriptionDate string                  `json:"transcription_date"`
}

func (s *SpeechAnalyticsService) GetSpeechAnalyticsCalls(payload GetSpeechAnalyticsCallsPayload) (*utils.PaginationResponse[SpeechAnalyticsCallItem], *utils.HttpErrorResponse) {
	validate := utils.GetValidate()
	err := validate.Struct(payload)

	if err != nil {
		return nil, &utils.HttpErrorResponse{Message: err.Error()}
	}

	return utils.Post[utils.PaginationResponse[SpeechAnalyticsCallItem]](*s.httpConfig, "/v1/cdr", payload, utils.PaginationResponse[SpeechAnalyticsCallItem]{})
}

func (s *SpeechAnalyticsService) TranscribeCallById(callId string, payload TranscribeCallByIdPayload) (*utils.HttpSuccessBasicResponse, *utils.HttpErrorResponse) {
	validate := utils.GetValidate()
	err := validate.Struct(payload)

	if err != nil {
		return nil, &utils.HttpErrorResponse{Message: err.Error()}
	}

	url := path.Join("/v1/cdr", callId, "retranscribe")

	return utils.Put[utils.HttpSuccessBasicResponse](*s.httpConfig, url, payload, utils.HttpSuccessBasicResponse{})
}

func (s *SpeechAnalyticsService) RequestTranscriptionByCallId(callId string) (*RequestTranscriptionByCallIdResponse, *utils.HttpErrorResponse) {
	url := path.Join("/v1/cdr", callId, "transcription")

	return utils.Get[RequestTranscriptionByCallIdResponse](*s.httpConfig, url, RequestTranscriptionByCallIdResponse{})
}
