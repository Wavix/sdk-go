package wavix

import (
	"github.com/wavix/sdk-go/utils"
)

type CdrServiceInterface interface {
	GetCdrList(queryParams GetCdrListQueryParams) (*utils.PaginationResponse[CdrListItem], *utils.HttpErrorResponse)
}

type CdrService struct {
	httpConfig *utils.HttpConfig
}

type CdrListItem struct {
	Charge      string `json:"charge"`
	Date        string `json:"date"`
	Destination string `json:"destination"`
	Disposition string `json:"disposition"`
	Duration    int    `json:"duration"`
	ForwardFee  string `json:"forward_fee"`
	From        string `json:"from"`
	PerMinute   string `json:"per_minute"`
	To          string `json:"to"`
	Uuid        string `json:"uuid"`
}

type GetCdrListQueryParams struct {
	utils.PaginationParams
	utils.RequiredDateParams
	Type        string `validate:"required,oneof=placed received" url:"type"`
	Disposition string `validate:"omitempty,oneof=answered noanswer busy failed all" url:"disposition,omitempty"`
	FromSearch  string `url:"from_search,omitempty"`
	ToSearch    string `url:"to_search,omitempty"`
	SipTrunk    string `url:"sip_trunk,omitempty"`
	Uuid        string `url:"uuid,omitempty"`
}

func (s *CdrService) GetCdrList(queryParams GetCdrListQueryParams) (*utils.PaginationResponse[CdrListItem], *utils.HttpErrorResponse) {
	validate := utils.GetValidate()
	err := validate.Struct(queryParams)

	if err != nil {
		return nil, &utils.HttpErrorResponse{Message: err.Error()}
	}

	url := utils.BuildUrlWithQueryString("/v1/cdr", queryParams)

	return utils.Get[utils.PaginationResponse[CdrListItem]](*s.httpConfig, url, utils.PaginationResponse[CdrListItem]{})
}
