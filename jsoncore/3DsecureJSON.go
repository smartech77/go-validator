package jsoncore

import (
	"encoding/json"
	"log"
)

// =====================================================
//
// 3D SECURE VALIDATION
//
// Documentation: https://uat.valitorpay.com/index.html#operation/CardVerification
//
// =====================================================

type CardVerificationData struct {
	VerifyingEnrollmentResponse              string `json:"verifyingEnrollmentResponse"`
	PayerAuthenticationResponse              string `json:"payerAuthenticationResponse"`
	CardholderAuthenticationVerificationData string `json:"cardholderAuthenticationVerificationData"`
}

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
func (cs CompanyService) VerifyCardUsing3DSecure(cardVerification *CardVerification) (card Card, err error) {

	verificationAsJSON, jsonError := json.Marshal(cardVerification)
	if jsonError != nil {
		err = jsonError
		return
	}
	log.Println(verificationAsJSON)
	return
}
