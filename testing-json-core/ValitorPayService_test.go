package test

import (
	valitor "github.com/zkynetio/go-valitor"
	jsoncore "github.com/zkynetio/go-valitor/jsoncore"
)

var TestCardJSON = &jsoncore.Card{
	Number:        "5304259906522887",
	ExpYear:       22,
	ExpMonth:      11,
	CVC:           "749",
	VirtualNumber: "4999993986001010",
}
var TCSJSON = valitor.NewValitorPayService(
	"053128",
	"225",
	"",
)

// func Test_CompanyService_CreateAVirtualCard(t *testing.T) {

// 	virtualCardResponse := TCSJSON.CreateAVirtualCard(TestCardJSON, "CardholderInitiatedCredentialOnFile", "ECommerceWithCvc", "")
// 	log.Println("Test output!")
// 	log.Println(virtualCardResponse)
// }
// func Test_CompanyService_UpdateAVirtualCardsExpirationDate(t *testing.T) {

// 	VirtualCardExpirationUpdateResponse := TCSJSON.UpdateAVirtualCardsExpirationDate(TestCardJSON, "ECommerceWithCvc")
// 	log.Println("Test output!")
// 	log.Println(VirtualCardExpirationUpdateResponse)
// }
