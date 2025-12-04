package wavix

import (
	"encoding/json"
	"log"
	"net/url"
	"path"

	"github.com/gorilla/websocket"
	"github.com/wavix/sdk-go/utils"
)

type CallServiceInterface interface {
	Connect() *utils.HttpErrorResponse
	Disconnect()
	OnEvent(callback EventCallback)
	GetList() (*CallsResponse, *utils.HttpErrorResponse)
	Get(callId string) (*CallResponse, *utils.HttpErrorResponse)
	StartCall(payload StartCallPayload) (*CallEvent, *StartCallErrorResponse)
	Answer(callId string, payload AnswerCallPayload) (*utils.HttpSuccessBasicResponse, *utils.HttpErrorResponse)
	UpdateTag(callId string, tag string) (*utils.HttpSuccessBasicResponse, *utils.HttpErrorResponse)
	PlayAudio(callId string, audioUrl string) (*utils.HttpSuccessBasicResponse, *utils.HttpErrorResponse)
	Hangup(callId string) (*utils.HttpSuccessBasicResponse, *utils.HttpErrorResponse)
}

type CallService struct {
	ws     *websocket.Conn
	http   *utils.HttpConfig
	events chan CallEvent
}

type EventCallback func(event CallEvent)
type CallDirection string
type EventType string

const (
	InboundCallDirection  CallDirection = "inbound"
	OutboundCallDirection CallDirection = "outbound"
)

const (
	CallSetupEventType   EventType = "call_setup"
	RingingEventType     EventType = "ringing"
	EarlyMediaEventType  EventType = "early_media"
	AnsweredEventType    EventType = "answered"
	CompletedEventType   EventType = "completed"
	BusyEventType        EventType = "busy"
	CancelledEventType   EventType = "cancelled"
	RejectedEventType    EventType = "rejected"
	OnCallEventEventType EventType = "on_call_event"
)

type AudioEventPayload struct {
	Type    string `json:"type"`
	Payload struct {
		Status string `json:"status"`
	} `json:"payload"`
}

type CallEvent struct {
	Uuid            string             `json:"uuid"`
	Direction       CallDirection      `json:"direction"`
	EventType       EventType          `json:"event_type"`
	EventTime       string             `json:"event_time"`
	EventPayload    *AudioEventPayload `json:"event_payload"`
	From            string             `json:"from"`
	To              string             `json:"to"`
	CallStarted     string             `json:"call_started"`
	CallAnswered    *string            `json:"call_answered"`
	CallCompleted   *string            `json:"call_completed"`
	MachineDetected bool               `json:"machine_detected"`
	Tag             *string            `json:"tag"`
}

type CallsResponse struct {
	Success bool        `json:"success"`
	Calls   []CallEvent `json:"calls"`
}

type CallResponse struct {
	Success bool      `json:"success"`
	Call    CallEvent `json:"call"`
}

type StartCallPayload struct {
	From               string `validate:"required" json:"from"`
	To                 string `validate:"required" json:"to"`
	CallbackUrl        string `validate:"required" json:"callback_url"`
	Recording          bool   `json:"recording"`
	VoicemailDetection bool   `json:"voicemail_detection"`
	Timeout            *int   `json:"timeout,omitempty"`
	Tag                string `json:"tag,omitempty"`
}

type StartCallErrorResponse struct {
	Success bool              `json:"success"`
	Message string            `json:"message"`
	Errors  map[string]string `json:"errors"`
}

type AnswerCallPayload struct {
	CallRecording     bool `json:"call_recording"`
	CallTranscription bool `json:"call_transcription"`
}

type UpdateTagPayload struct {
	Tag string `json:"tag"`
}

type PlayAudioPayload struct {
	AudioFile string `json:"audio_file"`
}

type TransferPayload struct {
	From                 string `validate:"required" json:"from"`
	To                   string `validate:"required" json:"to"`
	CallRecording        bool   `json:"call_recording"`
	DualChannelRecording bool   `json:"dual_channel_recording"`
	MachineDetection     bool   `json:"machine_detection"`
	APlaybackAudio       string `json:"a_playback_audio,omitempty"`
	BPlaybackAudio       string `json:"b_playback_audio,omitempty"`
}

func (s *CallService) Connect() *utils.HttpErrorResponse {
	parsedUrl, err := url.Parse(s.http.BaseUrl)

	if err != nil {
		return &utils.HttpErrorResponse{Message: err.Error()}
	}

	wsUrl := "wss://" + parsedUrl.Host + "/sip?appid=" + s.http.AppId

	s.ws, _, err = websocket.DefaultDialer.Dial(wsUrl, nil)

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

func (s *CallService) GetList() (*CallsResponse, *utils.HttpErrorResponse) {
	return utils.Get(*s.http, "/v2/calls", CallsResponse{})
}

func (s *CallService) Get(callId string) (*CallResponse, *utils.HttpErrorResponse) {
	url := path.Join("/v2/calls", callId)
	return utils.Get(*s.http, url, CallResponse{})
}

func (s *CallService) StartCall(payload StartCallPayload) (*CallEvent, *StartCallErrorResponse) {
	validate := utils.GetValidate()
	err := validate.Struct(payload)

	if err != nil {
		return nil, &StartCallErrorResponse{Success: false, Message: err.Error(), Errors: map[string]string{}}
	}

	callEvent, httpError := utils.Post(*s.http, "/v2/calls", payload, CallEvent{})

	if httpError != nil {
		return nil, &StartCallErrorResponse{Success: false, Message: httpError.Message, Errors: map[string]string{}}
	}

	return callEvent, nil
}

func (s *CallService) Answer(callId string, payload AnswerCallPayload) (*utils.HttpSuccessBasicResponse, *utils.HttpErrorResponse) {
	url := path.Join("/v2/calls", callId, "answer")
	return utils.Post(*s.http, url, payload, utils.HttpSuccessBasicResponse{Success: true})
}

func (s *CallService) UpdateTag(callId string, tag string) (*utils.HttpSuccessBasicResponse, *utils.HttpErrorResponse) {
	url := path.Join("/v2/calls", callId)
	payload := UpdateTagPayload{Tag: tag}
	return utils.Patch(*s.http, url, payload, utils.HttpSuccessBasicResponse{Success: true})
}

func (s *CallService) PlayAudio(callId string, audioUrl string) (*utils.HttpSuccessBasicResponse, *utils.HttpErrorResponse) {
	url := path.Join("/v2/calls", callId, "play")
	payload := PlayAudioPayload{AudioFile: audioUrl}
	return utils.Post(*s.http, url, payload, utils.HttpSuccessBasicResponse{Success: true})
}

func (s *CallService) Hangup(callId string) (*utils.HttpSuccessBasicResponse, *utils.HttpErrorResponse) {
	url := path.Join("/v2/calls", callId)

	return utils.Delete(*s.http, url, utils.HttpSuccessBasicResponse{Success: true})
}
