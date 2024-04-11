package wavix

import (
	"github.com/wavix/sdk-go/utils"
)

type SmsServiceInterface interface {
	SendMessage(payload SendMessagePayload) (*MessageResponseBody, *utils.HttpErrorResponse)
}

type MessageBody struct {
	Text  string    `json:"text"`
	Media *[]string `json:"media,omitempty"`
}

type SendMessagePayload struct {
	From        string      `json:"from"`
	To          string      `json:"to"`
	MessageBody MessageBody `json:"message_body"`
	CallbackUrl *string     `json:"callback_url,omitempty"`
	Validity    *int        `json:"validity,omitempty"`
	ExternalId  *string     `json:"external_id,omitempty"`
}

type MessageResponseBody struct {
	Charge       string      `json:"charge"`
	DeliveredAt  *string     `json:"delivered_at"`
	Direction    string      `json:"direction"`
	ErrorMessage *string     `json:"error_message"`
	From         string      `json:"from"`
	To           string      `json:"to"`
	Mcc          string      `json:"mcc"`
	Mnc          string      `json:"mnc"`
	MessageBody  MessageBody `json:"message_body"`
	MessageId    string      `json:"message_id"`
	MessageType  string      `json:"message_type"`
	Segments     int         `json:"segments"`
	SentAt       *string     `json:"sent_at"`
	Status       string      `json:"status"`
	SubmittedAt  string      `json:"submitted_at"`
	Tag          *string     `json:"tag"`
}

type SmsService struct {
	httpConfig *utils.HttpConfig
}

func (s *SmsService) SendMessage(payload SendMessagePayload) (*MessageResponseBody, *utils.HttpErrorResponse) {
	return utils.Post[MessageResponseBody](*s.httpConfig, "/v2/messages", payload, MessageResponseBody{})
}
