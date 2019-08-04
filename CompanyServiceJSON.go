package valitor

import (
	"encoding/json"
	"math/rand"
	"strconv"
	"sync"
	"time"

	"github.com/google/uuid"
)

// NewCompanyServiceUsingJSON ...
// This payment service will use the new JSON api from Valitor
// Documentation: https://uat.valitorpay.com
func NewCompanyServiceUsingJSON(
	agreementNumber string,
	terminalID string,
	url string,
) *CompanyServiceJSON {

	if url == "" {
		// Setting the default url as the test url
		url = "https://uat.valitorpay.com"
	}
	return &CompanyServiceJSON{
		Settings: &SettingsJSON{
			AgreementNumber: agreementNumber,
			TerminalID:      terminalID,
			URL:             url,
		},
	}
}

type CompanyServiceJSON struct {
	Settings *SettingsJSON
	Mux      sync.RWMutex
}
type SettingsJSON struct {
	AgreementNumber string
	TerminalID      string
	URL             string
}

// =====================================================
//
// CREATING A VIRTUAL CARD
//
// Documentation: https://uat.valitorpay.com/index.html#operation/CreateVirtualCard
//
// =====================================================

// VirtualCardRequest ...
// Documentation: https://uat.valitorpay.com/index.html#operation/CreateVirtualCard
type VirtualCardRequest struct {
	CardNumber                string                `json:"cardNumber"`
	ExpirationMonth           int                   `json:"expirationMonth"`
	ExpirationYear            int                   `json:"expirationYear"`
	Cvc                       string                `json:"cvc"`
	AgreementNumber           string                `json:"agreementNumber"`
	TerminalID                string                `json:"terminalId"`
	SubsequentTransactionType string                `json:"subsequentTransactionType"`
	TransactionType           string                `json:"transactionType"`
	TransactionLifecycleID    string                `json:"TransactionLifecycleID,omitempty"`
	CardVerificationData      *CardVerificationData `json:"CardVerificationData,omitempty"`
}

// VirtualCardResponse ...
// Documentation: https://uat.valitorpay.com/index.html#operation/CreateVirtualCard
type VirtualCardResponse struct {
	// The system error is unique to this module
	// It will be (nil) unless something went wrong on the system level
	SystemError            error
	TransactionLifecycleID string
	VirtualCard            string `json:"virtualCard"`
	IsSuccess              bool   `json:"isSuccess"`
	Code                   string `json:"responseCode"`
	Description            string `json:"responseDescription"`
}

// CreateAVirtualCard ...
// Documentation: https://uat.valitorpay.com/index.html#operation/CreateVirtualCard
func (cs *CompanyServiceJSON) CreateAVirtualCard(
	card *Card,
	cardVerificationData *CardVerificationData,
	subsequentTransactionType, transactionType, transactionLifecycleID string,
) (response VirtualCardResponse) {

	Request := &VirtualCardRequest{
		CardNumber:                card.Number,
		ExpirationMonth:           card.ExpMonth,
		ExpirationYear:            card.ExpYear,
		Cvc:                       card.CVC,
		AgreementNumber:           cs.Settings.AgreementNumber,
		TerminalID:                cs.Settings.TerminalID,
		SubsequentTransactionType: subsequentTransactionType,
		TransactionType:           transactionType,
		CardVerificationData:      cardVerificationData,
	}

	// if there is no TransactionLifecycleID we try to make a new one.
	if transactionLifecycleID == "" {
		rand.Seed(time.Now().UnixNano())
		id, err := uuid.NewUUID()
		if err == nil {
			Request.TransactionLifecycleID = id.String()
			response.TransactionLifecycleID = id.String()
		}
	}

	requestAsJSON, err := json.Marshal(Request)

	if err != nil {
		response.SystemError = err
		return
	}
	resp, code, err := SendJSON(requestAsJSON, "POST", cs.Settings.URL+"/VirtualCard/CreateVirtualCard")
	if err != nil {
		response.SystemError = err
		return
	}
	if code != 200 {
		response.Code = strconv.Itoa(code)
		response.IsSuccess = false
		response.Description = getDescriptionForNone200Code(code)
	}
	if err := json.Unmarshal(resp, &response); err != nil {
		response.SystemError = err
		return
	}

	return
}

// =====================================================
//
// UPDATING THE EXPIRATION DATE OF A VIRTUAL CARD
//
// Documentation: https://uat.valitorpay.com/index.html#operation/UpdateExpirationDate
//
// =====================================================

// VirtualCardExpirationUpdateRequest ...
// Documentation: https://uat.valitorpay.com/index.html#operation/UpdateExpirationDate
type VirtualCardExpirationUpdateRequest struct {
	VirtualCardNumber    string                `json:"virtualCardNumber"`
	ExpirationMonth      int                   `json:"expirationMonth"`
	ExpirationYear       int                   `json:"expirationYear"`
	Cvc                  string                `json:"cvc"`
	AgreementNumber      string                `json:"agreementNumber"`
	TerminalID           string                `json:"terminalId"`
	TransactionType      string                `json:"transactionType"`
	CardVerificationData *CardVerificationData `json:"CardVerificationData,omitempty"`
}

// VirtualCardExpirationUpdateResponse ...
// Documentation: https://uat.valitorpay.com/index.html#operation/UpdateExpirationDate
type VirtualCardExpirationUpdateResponse struct {
	// The system error is unique to this module
	// It will be (nil) unless something went wrong on the system level
	SystemError error
	IsSuccess   bool   `json:"isSuccess"`
	Code        string `json:"responseCode"`
	Description string `json:"responseDescription"`
}

// UpdateAVirtualCardsExpirationDate ...
// Documentation: https://uat.valitorpay.com/index.html#operation/UpdateExpirationDate
func (cs *CompanyServiceJSON) UpdateAVirtualCardsExpirationDate(
	card *Card,
	cardVerificationData *CardVerificationData,
	transactionType string,
) (response VirtualCardExpirationUpdateResponse) {

	Request := &VirtualCardExpirationUpdateRequest{
		VirtualCardNumber:    card.VirtualNumber,
		ExpirationMonth:      card.ExpMonth,
		ExpirationYear:       card.ExpYear,
		Cvc:                  card.CVC,
		AgreementNumber:      cs.Settings.AgreementNumber,
		TerminalID:           cs.Settings.TerminalID,
		TransactionType:      transactionType,
		CardVerificationData: cardVerificationData,
	}
	requestAsJSON, err := json.Marshal(Request)

	if err != nil {
		response.SystemError = err
		return
	}
	resp, code, err := SendJSON(requestAsJSON, "POST", cs.Settings.URL+"/VirtualCard/UpdateExpirationDate")
	if err != nil {
		response.SystemError = err
		return
	}
	if code != 200 {
		response.Code = strconv.Itoa(code)
		response.IsSuccess = false
		response.Description = getDescriptionForNone200Code(code)
	}
	if err := json.Unmarshal(resp, &response); err != nil {
		response.SystemError = err
		return
	}

	return
}

// =====================================================
//
// MAKE A DIRRECT PAYMENT WITH A REAL CARD
//
// Documentation: https://uat.valitorpay.com/index.html#operation/CardPayment
//
// =====================================================

// CardPaymentRequest ...
// Documentation: https://uat.valitorpay.com/index.html#operation/CardPayment
type CardPaymentRequest struct {
	Operation                 string                     `json:"operation"`
	TransactionType           string                     `json:"transactionType"`
	Currency                  string                     `json:"currency"`
	Amount                    int                        `json:"amount"`
	TerminalID                string                     `json:"terminalId"`
	AgreementNumber           string                     `json:"agreementNumber"`
	CardNumber                string                     `json:"cardNumber"`
	ExpirationMonth           int                        `json:"expirationMonth"`
	ExpirationYear            int                        `json:"expirationYear"`
	Cvc                       string                     `json:"cvc"`
	ReferenceNumber           string                     `json:"referenceNumber"`
	UseAsFirstTransaction     string                     `json:"useAsFirstTransaction,omitempty"`
	CardVerificationData      *CardVerificationData      `json:"cardVerificationData,omitempty"`
	SubsequentTransactionData *SubsequentTransactionData `json:"subsequentTransactionData,omitempty"`
	DCCData                   *DCCData                   `json:"dccData,omitempty"`
}

// CardPaymentResponse ...
// Documentation: https://uat.valitorpay.com/index.html#operation/CardPayment
type CardPaymentResponse struct {
	SystemError               error
	ReferenceNumber           string `json:"referenceNumber"`
	TransactionID             string `json:"transactionID"`
	AuthorizationCode         string `json:"authorizationCode"`
	TransactionLifecycleID    string `json:"transactionLifecycleId"`
	AuthorizationResponseTime string `json:"authorizationResponseTime"`
	IsSuccess                 bool   `json:"isSuccess"`
	Code                      string `json:"responseCode"`
	Description               string `json:"responseDescription"`
}

// CardPayment ...
// Documentation: https://uat.valitorpay.com/index.html#operation/CardPayment
func (cs *CompanyServiceJSON) CardPayment(
	card *Card,
	operation string,
	transactionType string,
	currency string,
	amount int,
	referenceNumer string,
	useAsFirstTransaction string,
	subsequentTransactionData *SubsequentTransactionData,
	cardVerificationData *CardVerificationData,
	dccData *DCCData,
) (response CardPaymentResponse) {

	Request := &CardPaymentRequest{
		Operation:                 operation,
		CardNumber:                card.Number,
		ExpirationMonth:           card.ExpMonth,
		ExpirationYear:            card.ExpYear,
		Cvc:                       card.CVC,
		AgreementNumber:           cs.Settings.AgreementNumber,
		TerminalID:                cs.Settings.TerminalID,
		Amount:                    amount,
		Currency:                  currency,
		ReferenceNumber:           referenceNumer,
		UseAsFirstTransaction:     useAsFirstTransaction,
		TransactionType:           transactionType,
		CardVerificationData:      cardVerificationData,
		SubsequentTransactionData: subsequentTransactionData,
		DCCData:                   dccData,
	}
	requestAsJSON, err := json.Marshal(Request)

	if err != nil {
		response.SystemError = err
		return
	}
	resp, code, err := SendJSON(requestAsJSON, "POST", cs.Settings.URL+"/Payment/CardPayment")
	if err != nil {
		response.SystemError = err
		return
	}
	// code 400 will send back a json object describing the error.
	// it is sufficient to forward this object.
	if code == 400 {
		response.Code = strconv.Itoa(400)
		response.IsSuccess = false
		response.Description = string(resp)
	} else if code != 200 {
		response.Code = strconv.Itoa(code)
		response.IsSuccess = false
		response.Description = getDescriptionForNone200Code(code)
	}
	if err := json.Unmarshal(resp, &response); err != nil {
		response.SystemError = err
		return
	}

	return
}

// =====================================================
//
// MAKE A PAYMENT WITH A VIRTUAL CARD
//
// Documentation: https://uat.valitorpay.com/index.html#operation/VirtualCardPayment
//
// =====================================================

// VirtualCardPaymentRequest ...
// Documentation: https://uat.valitorpay.com/index.html#operation/VirtualCardPayment
type VirtualCardPaymentRequest struct {
	Operation         string `json:"operation"`
	Currency          string `json:"currency"`
	Amount            int    `json:"amount"`
	TerminalID        string `json:"terminalId"`
	AgreementNumber   string `json:"agreementNumber"`
	VirtualCardNumber string `json:"virtualCardNumber"`
	ReferenceNumber   string `json:"referenceNumber"`
	InitiationReason  string `json:"initiationReason,omitempty"`
}

// VirtualCardPaymentResponse ...
// Documentation: https://uat.valitorpay.com/index.html#operation/VirtualCardPayment
type VirtualCardPaymentResponse struct {
	SystemError               error
	ReferenceNumber           string `json:"referenceNumber"`
	TransactionID             string `json:"transactionID"`
	AuthorizationCode         string `json:"authorizationCode"`
	TransactionLifecycleID    string `json:"transactionLifecycleId"`
	AuthorizationResponseTime string `json:"authorizationResponseTime"`
	IsSuccess                 bool   `json:"isSuccess"`
	Code                      string `json:"responseCode"`
	Description               string `json:"responseDescription"`
}

// VirtualCardPayment ...
// Documentation: https://uat.valitorpay.com/index.html#operation/VirtualCardPayment
func (cs *CompanyServiceJSON) VirtualCardPayment(
	card *Card,
	initialReason string,
	currency string,
	amount int,
	referenceNumer string,

) (response VirtualCardPaymentResponse) {

	Request := &VirtualCardPaymentRequest{
		VirtualCardNumber: card.VirtualNumber,
		AgreementNumber:   cs.Settings.AgreementNumber,
		TerminalID:        cs.Settings.TerminalID,
		Amount:            amount,
		Currency:          currency,
		ReferenceNumber:   referenceNumer,
		InitiationReason:  initialReason,
	}
	requestAsJSON, err := json.Marshal(Request)

	if err != nil {
		response.SystemError = err
		return
	}
	resp, code, err := SendJSON(requestAsJSON, "POST", cs.Settings.URL+"/Payment/VirtualCardPayment")
	if err != nil {
		response.SystemError = err
		return
	}
	if code != 200 {
		response.Code = strconv.Itoa(code)
		response.IsSuccess = false
		response.Description = getDescriptionForNone200Code(code)
	}
	if err := json.Unmarshal(resp, &response); err != nil {
		response.SystemError = err
		return
	}

	return
}

// =====================================================
//
// GET A DCC OFFER
//
// This is an offer of currency convertion.
//
// Documentation: https://uat.valitorpay.com/index.html#operation/DccOffer
//
// =====================================================

// DCCOfferRequest ...
// Documentation:  https://uat.valitorpay.com/index.html#operation/DccOffer
type DCCOfferRequest struct {
	CardNumber      string `json:"cardNumber"`
	Amount          int    `json:"amount"`
	Currency        string `json:"currency"`
	AgreementNumber string `json:"agreementNumber"`
	TerminalID      string `json:"terminalId"`
}

// DCCOfferResponse ...
// Documentation:  https://uat.valitorpay.com/index.html#operation/DccOffer
type DCCOfferResponse struct {
	SystemError                  error
	Currency                     string  `json:"currency"`
	Amount                       int     `json:"amount"`
	OfferCurrency                string  `json:"offerCurrency"`
	OfferAmount                  int     `json:"offerAmount"`
	DccCardholderBillingFee      int     `json:"dccCardholderBillingFee"`
	ExchangeRate                 float64 `json:"exchangeRate"`
	DccInformationEncryptedValue string  `json:"dccInformationEncryptedValue"`
	ResponseTimestamp            string  `json:"responseTimestamp"`
	IsSuccess                    bool    `json:"isSuccess"`
	Code                         string  `json:"responseCode"`
	Description                  string  `json:"responseDescription"`
}

// DCCOffer ...
// Documentation: https://uat.valitorpay.com/index.html#operation/VirtualCardPayment
func (cs *CompanyServiceJSON) DCCOffer(
	card *Card,
	currency string,
	amount int,
) (response DCCOfferResponse) {

	Request := &DCCOfferRequest{
		CardNumber:      card.Number,
		AgreementNumber: cs.Settings.AgreementNumber,
		TerminalID:      cs.Settings.TerminalID,
		Amount:          amount,
		Currency:        currency,
	}
	requestAsJSON, err := json.Marshal(Request)

	if err != nil {
		response.SystemError = err
		return
	}
	resp, code, err := SendJSON(requestAsJSON, "POST", cs.Settings.URL+"/Dcc")
	if err != nil {
		response.SystemError = err
		return
	}
	if code != 200 {
		response.Code = strconv.Itoa(code)
		response.IsSuccess = false
		response.Description = getDescriptionForNone200Code(code)
	}
	if err := json.Unmarshal(resp, &response); err != nil {
		response.SystemError = err
		return
	}

	return
}

// =====================================================
//
// GENERAL STRUCTS USED FOR REQUESTS AND RESPONSES
//
// =====================================================
type SubsequentTransactionData struct {
	IsStoredCredential        string `json:"isStoredCredential"`
	TransactionLifecycleID    string `json:"transactionLifecycleId,omitempty"`
	SubsequentTransactionType string `json:"subsequentTransactionType"`
}

type DCCData struct {
	originalTransAmount          int64   `json:"originalTransAmount"`
	originalTransCurrency        string  `json:"originalTransCurrency"`
	dccCardholderBillingFee      int64   `json:"dccCardholderBillingFee"`
	dccExchangeRate              float64 `json:"dccExchangeRate"`
	dccOfferCreationDate         string  `json:"dccOfferCreationDate"`
	dccInformationEncryptedValue string  `json:"dccInformationEncryptedValue"`
}

// =====================================================
//
// HELPER FUNCTIONS
//
// =====================================================

// For when we get a code other then 200 from valitor.
func getDescriptionForNone200Code(code int) string {
	switch code {
	case 401:
		return "Unauthorized..??"
	default:
		return "Unknown error from valitor, code: " + strconv.Itoa(code)
	}

}
