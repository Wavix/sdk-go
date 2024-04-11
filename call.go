package wavix

import (
	"encoding/json"
	"fmt"
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
	GetList() (*CallResponse, *utils.HttpErrorResponse)
	StartCall(payload StartCallPayload) (*CallEvent, *StartCallErrorResponse)
	PlayAudio(callId string, payload PlayAudioPayload) (*utils.HttpSuccessBasicResponse, *utils.HttpErrorResponse)
	Tts(callId string, payload TtsPayload) (*utils.HttpSuccessBasicResponse, *utils.HttpErrorResponse)
	Transfer(callId string, payload TransferPayload) (*utils.HttpSuccessBasicResponse, *utils.HttpErrorResponse)
	CollectDTMF(callId string, payload CollectDTMFPayload) (*utils.HttpSuccessBasicResponse, *utils.HttpErrorResponse)
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
	CallSetupEventType   EventType = "call_setup"
	CompletedEventType   EventType = "completed"
	InCallEventEventType EventType = "in_call_event"
	RingingEventType     EventType = "ringing"
)

type Call struct {
	Id         string `json:"uuid"`
	From       string `json:"from"`
	To         string `json:"to"`
	StartedAt  string `json:"call_started"`
	AnsweredAt string `json:"call_answered"`
}

type InCallEventData interface{}
type DigitsAndReasonEventData struct {
	Digits string `json:"digits"`
	Reason string `json:"reason"`
}

type PlaybackIdEventData struct {
	PlaybackId string `json:"playback_id"`
}

type CallEventPayload struct {
	InCallEvent     string          `json:"in_call_event"`
	InCallEventData InCallEventData `json:"in_call_event_data"`
}

func (payload *CallEventPayload) UnmarshalJSON(data []byte) error {
	var tmp struct {
		InCallEvent     string          `json:"in_call_event"`
		InCallEventData json.RawMessage `json:"in_call_event_data"`
	}

	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}

	payload.InCallEvent = tmp.InCallEvent

	var digitsReasonEventData DigitsAndReasonEventData
	if err := json.Unmarshal(tmp.InCallEventData, &digitsReasonEventData); err == nil {
		payload.InCallEventData = digitsReasonEventData
		return nil
	}

	var playbackIdEventData PlaybackIdEventData
	if err := json.Unmarshal(tmp.InCallEventData, &playbackIdEventData); err == nil {
		payload.InCallEventData = playbackIdEventData
		return nil
	}

	return fmt.Errorf("unknown in_call_event_data type")
}

type StartCallPayload struct {
	From             string `validate:"required" json:"from"`
	To               string `validate:"required" json:"to"`
	StatusCallback   string `validate:"required" json:"status_callback"`
	CallRecording    bool   `json:"call_recording"`
	MachineDetection bool   `json:"machine_detection"`
}

type StartCallErrorResponse struct {
	Success bool              `json:"success"`
	Message string            `json:"message"`
	Error   map[string]string `json:"error"`
}

type PlayAudioPayload struct {
	TimeoutBeforePlaying  int    `json:"timeout_before_playing,omitempty"`
	TimeoutBetweenPlaying int    `json:"timeout_between_playing,omitempty"`
	AudioUrl              string `validate:"required,url" json:"audio_file"`
}

type TtsPayload struct {
	Text               string `validate:"required" json:"text"`
	Voice              string `validate:"required,oneof=Ivy Joanna Kendra Kimberly Salli Joey Justin Matthew Conchita Lucia Enrique Marlene Vicki Hans Russian Tatyana Maxim" json:"voice"`
	DelayBeforePlaying int    `json:"delay_before_playing"`
	MaxRepeatCount     int    `json:"max_repeat_count"`
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

type CollectDTMFAudioPayload struct {
	Url            string `validate:"required" json:"url"`
	StopOnKeypress bool   `json:"stop_on_keypress"`
}
type CollectDTMFPayload struct {
	MinDigits            int                     `json:"min_digits,omitempty"`
	MaxDigits            int                     `json:"max_digits,omitempty"`
	Timeout              int                     `json:"timeout,omitempty"`
	TerminationCharacter string                  `json:"termination_character,omitempty"`
	Audio                CollectDTMFAudioPayload `json:"audio,omitempty"`
	CallbackUrl          string                  `json:"callback_url,omitempty"`
}

type CallEvent struct {
	Uuid            string            `json:"uuid"`
	EventType       EventType         `json:"event_type"`
	EventTime       string            `json:"event_time"`
	EventPayload    *CallEventPayload `json:"event_payload"`
	From            string            `json:"from"`
	To              string            `json:"to"`
	CallStarted     string            `json:"call_started"`
	CallAnswered    string            `json:"call_answered"`
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

func (s *CallService) GetList() (*CallResponse, *utils.HttpErrorResponse) {
	return utils.Get[CallResponse](*s.http, "/v1/call", CallResponse{})
}

func (s *CallService) StartCall(payload StartCallPayload) (*CallEvent, *StartCallErrorResponse) {
	validate := utils.GetValidate()
	err := validate.Struct(payload)

	if err != nil {
		return nil, &StartCallErrorResponse{Success: false, Message: err.Error(), Error: map[string]string{}}
	}

	callEvent, httpError := utils.Post[CallEvent](*s.http, "/v1/call", payload, CallEvent{})

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

	url := path.Join("/v1/call", callId, "play")

	return utils.Post[utils.HttpSuccessBasicResponse](*s.http, url, payload, utils.HttpSuccessBasicResponse{Success: true})
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

	url := path.Join("/v1/call", callId, "tts")

	return utils.Post[utils.HttpSuccessBasicResponse](*s.http, url, payload, utils.HttpSuccessBasicResponse{Success: true})
}

func (s *CallService) Transfer(callId string, payload TransferPayload) (*utils.HttpSuccessBasicResponse, *utils.HttpErrorResponse) {
	validate := utils.GetValidate()
	err := validate.Struct(payload)

	if err != nil {
		return nil, &utils.HttpErrorResponse{Message: err.Error()}
	}

	url := path.Join("/v1/call", callId, "transfer")

	return utils.Post[utils.HttpSuccessBasicResponse](*s.http, url, payload, utils.HttpSuccessBasicResponse{Success: true})
}

func (s *CallService) CollectDTMF(callId string, payload CollectDTMFPayload) (*utils.HttpSuccessBasicResponse, *utils.HttpErrorResponse) {
	validate := utils.GetValidate()
	err := validate.Struct(payload)

	if err != nil {
		return nil, &utils.HttpErrorResponse{Message: err.Error()}
	}

	url := path.Join("/v1/call", callId, "collect")

	return utils.Post[utils.HttpSuccessBasicResponse](*s.http, url, payload, utils.HttpSuccessBasicResponse{Success: true})
}

func (s *CallService) Hangup(callId string) (*utils.HttpSuccessBasicResponse, *utils.HttpErrorResponse) {
	url := path.Join("/v1/call", callId)

	return utils.Delete[utils.HttpSuccessBasicResponse](*s.http, url, utils.HttpSuccessBasicResponse{Success: true})
}
