package wavix

import "github.com/wavix/sdk-go/utils"

type LinkShortenerServiceInterface interface {
	GetShortLinkMetrics(queryParams GetShortLinksMetricsQueryParams) (*GetShortLinksMetricsResponse, *utils.HttpErrorResponse)
	CreateShortLink(payload CreateShortLinkPayload) (*CreateShortLinkResponse, *utils.HttpErrorResponse)
}

type LinkShortenerService struct {
	httpConfig *utils.HttpConfig
}

type GetShortLinksMetricsQueryParams struct {
	utils.RequiredDateParams
	Phone       string `url:"phone,omitempty"`
	UtmCampaign string `url:"utm_campaign,omitempty"`
}

type CreateShortLinkPayload struct {
	Link           string `validate:"required" json:"link,omitempty"`
	ExpirationTime string `validate:"required" json:"expiration_time,omitempty"`
	FallbackUrl    string `validate:"required" json:"fallback_url,omitempty"`
	Phone          string `validate:"required" json:"phone,omitempty"`
	UtmCampaign    string `validate:"required" json:"utm_campaign,omitempty"`
}

type ShortLinkMetricListItem struct {
	Latitude        float64 `json:"latitude"`
	Longitude       float64 `json:"longitude"`
	OperatingSystem string  `json:"operating_system"`
	Browser         string  `json:"browser"`
	Language        string  `json:"language"`
	Phone           string  `json:"phone"`
	UtmCampaign     string  `json:"utm_campaign"`
	CreatedAt       string  `json:"created_at"`
	UserId          int     `json:"user_id"`
}

type GetShortLinksMetricsResponse struct {
	Metrics []ShortLinkMetricListItem `json:"metrics"`
}

type CreateShortLinkResponse struct {
	ShortLink string `json:"short_link"`
}

func (s *LinkShortenerService) GetShortLinkMetrics(queryParams GetShortLinksMetricsQueryParams) (*GetShortLinksMetricsResponse, *utils.HttpErrorResponse) {
	validate := utils.GetValidate()
	err := validate.Struct(queryParams)

	if err != nil {
		return nil, &utils.HttpErrorResponse{Message: err.Error()}
	}

	url := utils.BuildUrlWithQueryString("/v1/short-links/metrics", queryParams)

	return utils.Get[GetShortLinksMetricsResponse](*s.httpConfig, url, GetShortLinksMetricsResponse{})
}

func (s *LinkShortenerService) CreateShortLink(payload CreateShortLinkPayload) (*CreateShortLinkResponse, *utils.HttpErrorResponse) {
	validate := utils.GetValidate()
	err := validate.Struct(payload)

	if err != nil {
		return nil, &utils.HttpErrorResponse{Message: err.Error()}
	}

	return utils.Post[CreateShortLinkResponse](*s.httpConfig, "/v1/short-links", payload, CreateShortLinkResponse{})
}
