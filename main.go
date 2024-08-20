package wavix

import (
	"github.com/wavix/sdk-go/utils"
)

type Instance struct {
	NumberValidation NumberValidationServiceInterface
	Sms              SmsServiceInterface
	Billing          BillingServiceInterface
	Cart             CartServiceInterface
	Buy              BuyServiceInterface
	Cdr              CdrServiceInterface
	Profile          ProfileServiceInterface
	SipTrunk         SipTrunkServiceInterface
	Did              DidServiceInterface
	E911             E911ServiceInterface
	LinkShortener    LinkShortenerServiceInterface
	TwoFa            TwoFaServiceInterface
	SpeechAnalytics  SpeechAnalyticsServiceInterface
	VoiceCampaign    VoiceCampaignServiceInterface
	Call             CallServiceInterface
}

type ClientOptions struct {
	Appid   string
	BaseURL string
}

func Init(options ClientOptions) *Instance {

	baseURL := getBaseURL(options.BaseURL)
	httpConfig := utils.InitHttpConfig(baseURL, options.Appid)

	return &Instance{
		NumberValidation: &ValidationService{httpConfig},
		Sms:              &SmsService{httpConfig},
		Billing:          &BillingService{httpConfig},
		Cart:             &CartService{httpConfig},
		Buy:              &BuyService{httpConfig},
		Cdr:              &CdrService{httpConfig},
		Profile:          &ProfileService{httpConfig},
		SipTrunk:         &SipTrunkService{httpConfig},
		Did:              &DidService{httpConfig},
		E911:             &E911Service{httpConfig},
		LinkShortener:    &LinkShortenerService{httpConfig},
		TwoFa:            &TwoFaService{httpConfig},
		SpeechAnalytics:  &SpeechAnalyticsService{httpConfig},
		VoiceCampaign:    &VoiceCampaignService{httpConfig},
		Call:             &CallService{http: httpConfig, events: make(chan CallEvent)},
	}
}

func getBaseURL(baseURL string) string {
	if baseURL == "" {
		return "https://api.wavix.com"
	}

	return baseURL
}
