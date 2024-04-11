package wavix

import (
	"fmt"
	"path"

	utils "github.com/wavix/sdk-go/utils"
)

type BuyServiceInterface interface {
	GetCountryList() (*GetCountryListResponse, *utils.HttpErrorResponse)
	GetRegionList(countryId int) (*GetRegionListResponse, *utils.HttpErrorResponse)
	GetCountryCitiesList(countryId int) (*GetCityListResponse, *utils.HttpErrorResponse)
	GetRegionCitiesList(countryId int, regionId int) (*GetCityListResponse, *utils.HttpErrorResponse)
	GetAvailableDids(countryId int, cityId int, queryParams GetAvailableDidsQueryParams) (*GetAvailableDidsPaginatedResponse, *utils.HttpErrorResponse)
}

type BuyService struct {
	httpConfig *utils.HttpConfig
}

type Country struct {
	Id                   int    `json:"id"`
	Name                 string `json:"name"`
	HasProvincesOrStates bool   `json:"has_provinces_or_states"`
}

type Region struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type City struct {
	Id       int    `json:"id"`
	AreaCode int    `json:"area_code"`
	Name     string `json:"name"`
}

type GetAvailableDidsQueryParams struct {
	utils.PaginationParams
	TextEnabledOnly bool   `url:"text_enabled_only,omitempty"`
	TypeFilter      string `url:"type_filter,omitempty"`
}

type GetCountryListResponse struct {
	Countries []Country `json:"countries"`
}

type GetCityListResponse struct {
	Cities []City `json:"cities"`
}

type GetRegionListResponse struct {
	Regions []Region `json:"regions"`
}

type GetAvailableDidsPaginatedResponse struct {
	Items      []CartDidItem    `json:"dids"`
	Pagination utils.Pagination `json:"pagination"`
}

func (s *BuyService) GetCountryList() (*GetCountryListResponse, *utils.HttpErrorResponse) {
	return utils.Get[GetCountryListResponse](*s.httpConfig, "/v1/buy/countries", GetCountryListResponse{})
}

func (s *BuyService) GetRegionList(countryId int) (*GetRegionListResponse, *utils.HttpErrorResponse) {
	url := path.Join("/v1/buy/countries", fmt.Sprintf("%d", countryId), "regions")

	return utils.Get[GetRegionListResponse](*s.httpConfig, url, GetRegionListResponse{})
}

func (s *BuyService) GetCountryCitiesList(countryId int) (*GetCityListResponse, *utils.HttpErrorResponse) {
	url := path.Join("/v1/buy/countries/", fmt.Sprintf("%d", countryId), "cities")

	return utils.Get[GetCityListResponse](*s.httpConfig, url, GetCityListResponse{})
}

func (s *BuyService) GetRegionCitiesList(countryId int, regionId int) (*GetCityListResponse, *utils.HttpErrorResponse) {
	url := path.Join("/v1/buy/countries", fmt.Sprintf("%d", countryId), "regions", fmt.Sprintf("%d", regionId), "cities")

	return utils.Get[GetCityListResponse](*s.httpConfig, url, GetCityListResponse{})
}

func (s *BuyService) GetAvailableDids(countryId int, cityId int, queryParams GetAvailableDidsQueryParams) (*GetAvailableDidsPaginatedResponse, *utils.HttpErrorResponse) {
	baseUrl := path.Join("/v1/buy/countries", fmt.Sprintf("%d", countryId), "cities", fmt.Sprintf("%d", cityId), "dids")
	url := utils.BuildUrlWithQueryString(baseUrl, queryParams)

	return utils.Get[GetAvailableDidsPaginatedResponse](*s.httpConfig, url, GetAvailableDidsPaginatedResponse{})
}
