package wavix

import (
	"github.com/wavix/sdk-go/utils"
)

type CartServiceInterface interface {
	GetCartContent() (*GetCartContentResponse, *utils.HttpErrorResponse)
	AddDidToCart(ids []string) (*AddDidToCartResponse, *utils.HttpErrorResponse)
	Checkout(ids []string) (*CheckoutResponse, *utils.HttpErrorResponse)
}

type CartService struct {
	httpConfig *utils.HttpConfig
}

type CartDidItem struct {
	Id               int      `json:"id"`
	ActivationFee    string   `json:"activation_fee"`
	Channels         int      `json:"channels"`
	City             string   `json:"city"`
	Country          string   `json:"country"`
	Cnam             bool     `json:"cnam"`
	CountryShortName string   `json:"country_short_name"`
	FreeMin          int      `json:"free_min"`
	MonthlyFee       string   `json:"monthly_fee"`
	Number           string   `json:"number"`
	PerMin           string   `json:"per_min"`
	RequireDocs      []string `json:"require_docs"`
	SmsEnabled       bool     `json:"sms_enabled"`
	SmsPrice         string   `json:"sms_price"`
}

type CartContentDocType struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Title string `json:"title"`
}

type AddDidToCartPayload struct {
	Ids []string `json:"ids"`
}

type CheckoutPayload struct {
	Ids []string `json:"ids"`
}

type GetCartContentResponse struct {
	Dids     []CartDidItem        `json:"dids"`
	DocTypes []CartContentDocType `json:"doc_types"`
}

type AddDidToCartResponse []CartDidItem

type CheckoutResponse struct {
	Success bool `json:"success"`
}

func (s *CartService) GetCartContent() (*GetCartContentResponse, *utils.HttpErrorResponse) {
	return utils.Get[GetCartContentResponse](*s.httpConfig, "/v1/buy/cart", GetCartContentResponse{})
}

func (s *CartService) AddDidToCart(ids []string) (*AddDidToCartResponse, *utils.HttpErrorResponse) {
	return utils.Put[AddDidToCartResponse](*s.httpConfig, "/v1/buy/cart",
		AddDidToCartPayload{Ids: ids}, AddDidToCartResponse{})
}

func (s *CartService) Checkout(ids []string) (*CheckoutResponse, *utils.HttpErrorResponse) {
	return utils.Post[CheckoutResponse](*s.httpConfig, "/v1/buy/cart/checkout", CheckoutPayload{Ids: ids}, CheckoutResponse{})
}
