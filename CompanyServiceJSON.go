package valitor

import (
	"encoding/json"
	"log"
	"strconv"
	"sync"

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
// 3D SECURE VALIDATION
//
// Documentation: https://uat.valitorpay.com/index.html#operation/CardVerification
//
// =====================================================

// CardVerification ...
// Documentation: https://uat.valitorpay.com/index.html#operation/CardVerification
type CardVerification struct {
	AgreementNumber         string `json:"agreementNumber"`
	TerminalID              string `json:"terminalId"`
	CardType                string `json:"cardType"`
	CardNumber              string `json:"cardNumber"`
	ExpirationMonth         int    `json:"expirationMonth"`
	ExpirationYear          int    `json:"expirationYear"`
	CardholderDeviceType    string `json:"cardholderDeviceType"`
	Amount                  int    `json:"amount"`
	Currency                string `json:"currency"`
	AuthorizationSuccessURL string `json:"authorizationSuccessUrl"`
	AuthorizationFailedURL  string `json:"authorizationFailedUrl"`
	MerchantData            string `json:"merchantData"`
}

// TODO: Finish this and get some info.
// VerifyCardUsing3DSecure ...
// Documentation: https://uat.valitorpay.com/index.html#operation/CardVerification
func (cs CompanyServiceJSON) VerifyCardUsing3DSecure(cardVerification *CardVerification) (card Card, err error) {

	verificationAsJSON, jsonError := json.Marshal(cardVerification)
	if jsonError != nil {
		err = jsonError
		return
	}
	log.Println(verificationAsJSON)
	return
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
	CardNumber                string `json:"cardNumber"`
	ExpirationMonth           int    `json:"expirationMonth"`
	ExpirationYear            int    `json:"expirationYear"`
	Cvc                       string `json:"cvc"`
	AgreementNumber           string `json:"agreementNumber"`
	TerminalID                string `json:"terminalId"`
	SubsequentTransactionType string `json:"subsequentTransactionType"`
	TransactionType           string `json:"transactionType"`
	TransactionLifecycleID    string `json:"TransactionLifecycleID,omitempty"`
	CardVerificationData      string `json:"CardVerificationData,omitempty"`
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
func (cs *CompanyServiceJSON) CreateAVirtualCard(card *Card, subsequentTransactionType, transactionType, transactionLifecycleID string) (response VirtualCardResponse) {

	Request := &VirtualCardRequest{
		CardNumber:                card.Number,
		ExpirationMonth:           card.ExpMonth,
		ExpirationYear:            card.ExpYear,
		Cvc:                       card.CVC,
		AgreementNumber:           cs.Settings.AgreementNumber,
		TerminalID:                cs.Settings.TerminalID,
		SubsequentTransactionType: subsequentTransactionType,
		TransactionType:           transactionType,
	}

	if card.VerificationDataFrom3DSecure != "" {
		Request.CardVerificationData = card.VerificationDataFrom3DSecure
	}
	// if there is no TransactionLifecycleID we try to make a new one.
	if transactionLifecycleID == "" {
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
	VirtualCardNumber    string `json:"virtualCardNumber"`
	ExpirationMonth      int    `json:"expirationMonth"`
	ExpirationYear       int    `json:"expirationYear"`
	Cvc                  string `json:"cvc"`
	AgreementNumber      string `json:"agreementNumber"`
	TerminalID           string `json:"terminalId"`
	TransactionType      string `json:"transactionType"`
	CardVerificationData string `json:"CardVerificationData,omitempty"`
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
// Documentation: https://uat.valitorpay.com/index.html#operation/CreateVirtualCard
func (cs *CompanyServiceJSON) UpdateAVirtualCardsExpirationDate(card *Card, transactionType string) (response VirtualCardExpirationUpdateResponse) {

	Request := &VirtualCardExpirationUpdateRequest{
		VirtualCardNumber: card.VirtualNumber,
		ExpirationMonth:   card.ExpMonth,
		ExpirationYear:    card.ExpYear,
		Cvc:               card.CVC,
		AgreementNumber:   cs.Settings.AgreementNumber,
		TerminalID:        cs.Settings.TerminalID,
		TransactionType:   transactionType,
	}
	if card.VerificationDataFrom3DSecure != "" {
		Request.CardVerificationData = card.VerificationDataFrom3DSecure
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

// For when we get a code other then 200 from valitor.
func getDescriptionForNone200Code(code int) string {
	switch code {
	case 401:
		return "Unauthorized..??"
	default:
		return "Unknown error from valitor, code: " + strconv.Itoa(code)
	}

}
