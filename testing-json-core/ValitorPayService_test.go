package test

import (
	"log"
	"testing"

	valitor "github.com/opensourcez/go-valitor"
	jsoncore "github.com/opensourcez/go-valitor/jsoncore"
)

var TestCardJSON = &jsoncore.Card{
	Number:        "5304259906522887",
	ExpYear:       2022,
	ExpMonth:      11,
	CVC:           "749",
	VirtualNumber: "4999993986001010",
}
var TCSJSON = valitor.NewValitorPayService(
	"053128",
	"225",
	"",
)

func Test_CompanyService_CreateAVirtualCard(t *testing.T) {

	cardVer := jsoncore.CardVerificationData{}
	virtualCardResponse := TCSJSON.CreateVirtualCard(TestCardJSON, &cardVer, "CardholderInitiatedCredentialOnFile", "ECommerceWithCvc", "")
	log.Println("Test output!")
	log.Println(virtualCardResponse)
}

// func Test_CompanyService_UpdateAVirtualCardsExpirationDate(t *testing.T) {

// 	VirtualCardExpirationUpdateResponse := TCSJSON.UpdateVirtualCardsExpirationDate(TestCardJSON, "ECommerceWithCvc")
// 	log.Println("Test output!")
// 	log.Println(VirtualCardExpirationUpdateResponse)
// }
