package wavix

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"path"

	"github.com/gorilla/websocket"
	"github.com/wavix/sdk-go/utils"
)

type CallServiceInterface interface {
	Connect() *utils.HttpErrorResponse
	Disconnect()
	OnEvent(callback EventCallback)
	GetList() (*CallResponse, *utils.HttpErrorResponse)
	GetCall(callId string) (*GetCallResponse, *utils.HttpErrorResponse)
	StartCall(payload StartCallPayload) (*CallEvent, *StartCallErrorResponse)
	PlayAudio(callId string, payload PlayAudioPayload) (*utils.HttpSuccessBasicResponse, *utils.HttpErrorResponse)
	StopAudio(callId string) (*utils.HttpSuccessBasicResponse, *utils.HttpErrorResponse)
	Tts(callId string, payload TtsPayload) (*utils.HttpSuccessBasicResponse, *utils.HttpErrorResponse)
	Transfer(callId string, payload TransferPayload) (*utils.HttpSuccessBasicResponse, *utils.HttpErrorResponse)
	CollectDTMF(callId string, payload CollectDTMFPayload) (*utils.HttpSuccessBasicResponse, *utils.HttpErrorResponse)
	UpdateCall(callId string, payload UpdateCallPayload) (*utils.HttpSuccessBasicResponse, *utils.HttpErrorResponse)
	Hangup(callId string) (*utils.HttpSuccessBasicResponse, *utils.HttpErrorResponse)
}

type CallService struct {
	ws     *websocket.Conn
	http   *utils.HttpConfig
	events chan CallEvent
}

type EventCallback func(event CallEvent)
type CallContext string
type EventType string
type EnglishVoice string
type SpanishVoice string
type GermanVoice string
type RussianVoice string

const (
	IvyEnglishVoice      EnglishVoice = "Ivy"
	JoannaEnglishVoice   EnglishVoice = "Joanna"
	KendraEnglishVoice   EnglishVoice = "Kendra"
	KimberlyEnglishVoice EnglishVoice = "Kimberly"
	SalliEnglishVoice    EnglishVoice = "Salli"
	JoeyEnglishVoice     EnglishVoice = "Joey"
	JustinEnglishVoice   EnglishVoice = "Justin"
	MatthewEnglishVoice  EnglishVoice = "Matthew"
	ConchitaSpanishVoice SpanishVoice = "Conchita"
	LuciaSpanishVoice    SpanishVoice = "Lucia"
	EnriqueSpanishVoice  SpanishVoice = "Enrique"
	MarleneGermanVoice   GermanVoice  = "Marlene"
	VickiGermanVoice     GermanVoice  = "Vicki"
	HansGermanVoice      GermanVoice  = "Hans"
	RussianRussianVoice  RussianVoice = "Russian"
	TatyanaRussianVoice  RussianVoice = "Tatyana"
	MaximRussianVoice    RussianVoice = "Maxim"
)

const (
	AnsweredEventType    EventType = "answered"
	BusyEventType        EventType = "busy"
	CallSetupEventType   EventType = "call_setup"
	CancelledEventType   EventType = "cancelled"
	CompletedEventType   EventType = "completed"
	EarlyMediaEventType  EventType = "early_media"
	OnCallEventEventType EventType = "on_call_event"
	RejectedEventType    EventType = "rejected"
	RingingEventType     EventType = "ringing"
	FailedEventType      EventType = "failed"
	TransferEventType    EventType = "transfer"
)

type Call struct {
	Id         string `json:"uuid"`
	From       string `json:"from"`
	To         string `json:"to"`
	StartedAt  string `json:"call_started"`
	AnsweredAt string `json:"call_answered"`
}

type OnCallEventPayload interface{}

type AudioEventPayload struct {
	PlaybackId string `json:"playback_id"`
	Status     string `json:"status"`
}

type CollectCompletedPayload struct {
	Digits string `json:"digits"`
	Reason string `json:"reason"`
}

type CallEventPayload struct {
	Type    string             `json:"type"`
	Payload OnCallEventPayload `json:"payload"`
}

func (payload *CallEventPayload) UnmarshalJSON(data []byte) error {
	var tmp struct {
		Type    string          `json:"type"`
		Payload json.RawMessage `json:"payload"`
	}

	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}

	payload.Type = tmp.Type

	switch tmp.Type {
	case "audio":
		var audioPayload AudioEventPayload
		if err := json.Unmarshal(tmp.Payload, &audioPayload); err != nil {
			return err
		}
		payload.Payload = audioPayload
	case "collect_completed":
		var collectPayload CollectCompletedPayload
		if err := json.Unmarshal(tmp.Payload, &collectPayload); err != nil {
			return err
		}
		payload.Payload = collectPayload
	default:
		return fmt.Errorf("unknown on_call_event type: %s", tmp.Type)
	}

	return nil
}

type StartCallPayload struct {
	From               string `validate:"required" json:"from"`
	To                 string `validate:"required" json:"to"`
	CallbackUrl        string `validate:"required" json:"callback_url"`
	Recording          bool   `json:"recording"`
	VoicemailDetection bool   `json:"voicemail_detection"`
	Timeout            int    `json:"timeout,omitempty"`
	Tag                string `json:"tag,omitempty"`
}

type StartCallErrorResponse struct {
	Success bool              `json:"success"`
	Message string            `json:"message"`
	Error   map[string]string `json:"error"`
}

type PlayAudioPayload struct {
	AudioUrl string `validate:"required,url" json:"audio_file"`
}

type TtsPayload struct {
	Text               string `validate:"required" json:"text"`
	Voice              string `json:"voice,omitempty"`
	DelayBeforePlaying int    `json:"delay_before_playing,omitempty"`
	MaxRepeatCount     int    `json:"max_repeat_count,omitempty"`
}

type TransferPayload struct {
	From                 string `validate:"required" json:"from"`
	To                   string `validate:"required" json:"to"`
	CallRecording        bool   `json:"call_recording"`
	DualChannelRecording bool   `json:"dual_channel_recording"`
	MachineDetection     bool   `json:"machine_detection"`
	APlaybackAudio       string `json:"a_playback_audio"`
	BPlaybackAudio       string `json:"b_playback_audio"`
}

type CollectPromptSay struct {
	Text     string `validate:"required" json:"text"`
	Voice    string `validate:"required" json:"voice"`
	Language string `json:"language,omitempty"`
}

type CollectPrompt struct {
	Play string            `json:"play,omitempty"`
	Say  *CollectPromptSay `json:"say,omitempty"`
}

type CollectDTMFPayload struct {
	MaxDigits            int            `json:"max_digits,omitempty"`
	Timeout              int            `json:"timeout,omitempty"`
	TerminationCharacter string         `json:"termination_character,omitempty"`
	MaxAttempts          int            `json:"max_attempts,omitempty"`
	StopOnKeypress       *bool          `json:"stop_on_keypress,omitempty"`
	Prompt               *CollectPrompt `json:"prompt,omitempty"`
}

type CallEvent struct {
	Uuid            string            `json:"uuid"`
	Direction       string            `json:"direction"`
	EventType       EventType         `json:"event_type"`
	EventTime       string            `json:"event_time"`
	EventPayload    *CallEventPayload `json:"event_payload"`
	From            string            `json:"from"`
	To              string            `json:"to"`
	CallStarted     string            `json:"call_started"`
	CallAnswered    *string           `json:"call_answered"`
	CallCompleted   *string           `json:"call_completed"`
	MachineDetected bool              `json:"machine_detected"`
	Tag             string            `json:"tag"`
}

type CallResponse struct {
	Calls []Call `json:"calls"`
}

func (s *CallService) Connect() *utils.HttpErrorResponse {
	parsedUrl, err := url.Parse(s.http.BaseUrl)

	if err != nil {
		return &utils.HttpErrorResponse{Message: err.Error()}
	}

	wsUrl := "wss://" + parsedUrl.Host + "/sip"

	headers := make(http.Header)
	if s.http.AppId != "" {
		headers.Set("Authorization", "Bearer "+s.http.AppId)
	}

	s.ws, _, err = websocket.DefaultDialer.Dial(wsUrl, headers)

	if err != nil {
		return &utils.HttpErrorResponse{Message: err.Error()}
	}

	go func() {
		for {
			_, message, err := s.ws.ReadMessage()

			if err != nil {
				log.Println("Error when read socket message:", err)
				continue
			}

			var event CallEvent
			err = json.Unmarshal(message, &event)

			if err != nil {
				log.Println("Error when parse socket JSON message:", err)
				continue
			}

			s.events <- event
		}
	}()

	return nil
}

func (s *CallService) OnEvent(callback EventCallback) {
	go func() {
		for event := range s.events {
			callback(event)
		}
	}()
}

func (s *CallService) Disconnect() {
	if s.ws != nil {
		s.ws.Close()
	}

	close(s.events)
}

func (s *CallService) GetList() (*CallResponse, *utils.HttpErrorResponse) {
	return utils.Get(*s.http, "/v1/calls", CallResponse{})
}

type GetCallResponse struct {
	Success bool      `json:"success"`
	Call    CallEvent `json:"call"`
}

func (s *CallService) GetCall(callId string) (*GetCallResponse, *utils.HttpErrorResponse) {
	url := path.Join("/v1/calls", callId)
	return utils.Get(*s.http, url, GetCallResponse{})
}

func (s *CallService) StartCall(payload StartCallPayload) (*CallEvent, *StartCallErrorResponse) {
	validate := utils.GetValidate()
	err := validate.Struct(payload)

	if err != nil {
		return nil, &StartCallErrorResponse{Success: false, Message: err.Error(), Error: map[string]string{}}
	}

	callEvent, httpError := utils.Post(*s.http, "/v1/calls", payload, CallEvent{})

	if httpError != nil {
		return nil, &StartCallErrorResponse{Success: false, Message: httpError.Message, Error: map[string]string{}}
	}

	return callEvent, nil
}

func (s *CallService) PlayAudio(callId string, payload PlayAudioPayload) (*utils.HttpSuccessBasicResponse, *utils.HttpErrorResponse) {
	validate := utils.GetValidate()
	err := validate.Struct(payload)

	if err != nil {
		return nil, &utils.HttpErrorResponse{Message: err.Error()}
	}

	url := path.Join("/v1/calls", callId, "play")

	return utils.Post(*s.http, url, payload, utils.HttpSuccessBasicResponse{Success: true})
}

/*
Voice list
https://docs.aws.amazon.com/polly/latest/dg/voicelist.html
*/
func (s *CallService) Tts(callId string, payload TtsPayload) (*utils.HttpSuccessBasicResponse, *utils.HttpErrorResponse) {
	validate := utils.GetValidate()
	err := validate.Struct(payload)

	if err != nil {
		return nil, &utils.HttpErrorResponse{Message: err.Error()}
	}

	if payload.Voice == "" {
		payload.Voice = "Joey"
	}

	url := path.Join("/v1/calls", callId, "tts")

	return utils.Post(*s.http, url, payload, utils.HttpSuccessBasicResponse{Success: true})
}

func (s *CallService) Transfer(callId string, payload TransferPayload) (*utils.HttpSuccessBasicResponse, *utils.HttpErrorResponse) {
	validate := utils.GetValidate()
	err := validate.Struct(payload)

	if err != nil {
		return nil, &utils.HttpErrorResponse{Message: err.Error()}
	}

	url := path.Join("/v1/calls", callId, "transfer")

	return utils.Post(*s.http, url, payload, utils.HttpSuccessBasicResponse{Success: true})
}

func (s *CallService) CollectDTMF(callId string, payload CollectDTMFPayload) (*utils.HttpSuccessBasicResponse, *utils.HttpErrorResponse) {
	validate := utils.GetValidate()
	err := validate.Struct(payload)

	if err != nil {
		return nil, &utils.HttpErrorResponse{Message: err.Error()}
	}

	url := path.Join("/v1/calls", callId, "collect")

	return utils.Post(*s.http, url, payload, utils.HttpSuccessBasicResponse{Success: true})
}

func (s *CallService) Hangup(callId string) (*utils.HttpSuccessBasicResponse, *utils.HttpErrorResponse) {
	url := path.Join("/v1/calls", callId)

	return utils.Delete(*s.http, url, utils.HttpSuccessBasicResponse{Success: true})
}

type UpdateCallPayload struct {
	Tag string `validate:"required" json:"tag"`
}

func (s *CallService) UpdateCall(callId string, payload UpdateCallPayload) (*utils.HttpSuccessBasicResponse, *utils.HttpErrorResponse) {
	validate := utils.GetValidate()
	err := validate.Struct(payload)

	if err != nil {
		return nil, &utils.HttpErrorResponse{Message: err.Error()}
	}

	url := path.Join("/v1/calls", callId)

	return utils.Patch(*s.http, url, payload, utils.HttpSuccessBasicResponse{Success: true})
}

func (s *CallService) StopAudio(callId string) (*utils.HttpSuccessBasicResponse, *utils.HttpErrorResponse) {
	url := path.Join("/v1/calls", callId, "audio")

	return utils.Delete(*s.http, url, utils.HttpSuccessBasicResponse{Success: true})
}
