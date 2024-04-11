package wavix

import (
	"fmt"

	"github.com/wavix/sdk-go/utils"
)

type TransactionType int

const (
	AdjustmentsTransactionType                     TransactionType = 0
	DealTransactionType                            TransactionType = 1
	ActivationTransactionType                      TransactionType = 2
	MonthTransactionType                           TransactionType = 3
	ActivationFeeTransactionType                   TransactionType = 4
	MonthFeeTransactionType                        TransactionType = 5
	CallTransactionType                            TransactionType = 6
	CallFeeTransactionType                         TransactionType = 7
	CallFixFeeTransactionType                      TransactionType = 11
	PaypalInTransactionType                        TransactionType = 8
	PaypalOutTransactionType                       TransactionType = 9
	TaxTransactionType                             TransactionType = 10
	WebcallTransactionType                         TransactionType = 12
	SipTransactionType                             TransactionType = 14
	SmsTransactionType                             TransactionType = 15
	ChannelTransactionType                         TransactionType = 16
	ChannelFeeTransactionType                      TransactionType = 17
	CallSkypeFeeTransactionType                    TransactionType = 18
	CcInTransactionType                            TransactionType = 19
	PaymentFeeTransactionType                      TransactionType = 20
	ConnectionTransactionType                      TransactionType = 21
	ConnectionFeeTransactionType                   TransactionType = 22
	PortingTransactionType                         TransactionType = 23
	InboundSmsTransactionType                      TransactionType = 24
	WireTransferTransactionType                    TransactionType = 25
	SubscriptionTransactionType                    TransactionType = 26
	SurchargeTransactionType                       TransactionType = 27
	HlrTransactionType                             TransactionType = 28
	NumberValidationTransactionType                TransactionType = 29
	CallRecordingTransactionType                   TransactionType = 30
	CallRecordingStorageTransactionType            TransactionType = 31
	CampaignBuilderRunTransactionType              TransactionType = 32
	VoicemailDetectionTransactionType              TransactionType = 33
	SenderIdDestinationRegistrationTransactionType TransactionType = 34
	SenderIdDestinationFeeTransactionType          TransactionType = 35
	TwoFaServiceTransactionType                    TransactionType = 36
	IvrTransactionType                             TransactionType = 37
	E911ActivationTransactionType                  TransactionType = 38
	MmsTransactionType                             TransactionType = 39
	InboundMmsTransactionType                      TransactionType = 40
	CallTranscriptionTransactionType               TransactionType = 41
	TendlcBrandsTransactionType                    TransactionType = 42
	TendlcCampaignFeeTransactionType               TransactionType = 43
	DidOrderTransactionType                        TransactionType = 44
	AdjustmentsInTransactionType                   TransactionType = 45
)

type BillingServiceInterface interface {
	GetAccountTransactions(params AccountTransactionsParams) (*AccountTransactionsPaginatedResponse, *utils.HttpErrorResponse)
	GetAccountInvoices(params AccountInvoicesParams) (*AccountInvoicesPaginatedResponse, *utils.HttpErrorResponse)
	DownloadInvoiceById(id int) ([]byte, *utils.HttpErrorResponse)
}

type BillingService struct {
	httpConfig *utils.HttpConfig
}

type AccountTransactionsPaginatedResponse struct {
	Items      []AccountTransactionsItem `json:"transactions"`
	Pagination utils.Pagination          `json:"pagination"`
}

type AccountInvoicesPaginatedResponse struct {
	Items      []AccountInvoiceItem `json:"invoices"`
	Pagination utils.Pagination     `json:"pagination"`
}

type AccountTransactionsParams struct {
	utils.OptionalDateParams
	utils.PaginationParams
}

type AccountInvoicesParams struct {
	utils.PaginationParams
}

type AccountTransactionsItem struct {
	Id           int             `json:"id"`
	Amount       string          `json:"amount"`
	BalanceAfter string          `json:"balance_after"`
	Date         string          `json:"date"`
	Details      string          `json:"details"`
	ShowInvoice  bool            `json:"show_invoice"`
	Status       string          `json:"status"`
	Type         TransactionType `json:"type"`
}

type DownloadInvoiceItem []byte

type AccountInvoiceItem struct {
	Id       int    `json:"id"`
	Amount   string `json:"amount"`
	FromDate string `json:"from_date"`
	ToDate   string `json:"to_date"`
}

func (s *BillingService) GetAccountTransactions(params AccountTransactionsParams) (*AccountTransactionsPaginatedResponse, *utils.HttpErrorResponse) {
	validate := utils.GetValidate()
	err := validate.Struct(params)

	if err != nil {
		return nil, &utils.HttpErrorResponse{Message: err.Error()}
	}

	url := utils.BuildUrlWithQueryString("/v1/billing/transactions", params)

	return utils.Get[AccountTransactionsPaginatedResponse](*s.httpConfig, url, AccountTransactionsPaginatedResponse{})
}

func (s *BillingService) GetAccountInvoices(params AccountInvoicesParams) (*AccountInvoicesPaginatedResponse, *utils.HttpErrorResponse) {
	url := utils.BuildUrlWithQueryString("/v1/billing/invoices", params)

	return utils.Get[AccountInvoicesPaginatedResponse](*s.httpConfig, url, AccountInvoicesPaginatedResponse{})
}

func (s *BillingService) DownloadInvoiceById(id int) ([]byte, *utils.HttpErrorResponse) {
	url := fmt.Sprintf("/v1/billing/invoices/%d", id)

	file, err := utils.Download(*s.httpConfig, url)

	if err != nil {
		return nil, err
	}

	return file, nil
}
