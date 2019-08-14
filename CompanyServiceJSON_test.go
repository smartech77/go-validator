package valitor

// 5304259906522887 2211 749
// 5304259909334470 2211 813
// 5304259902386667 2211 376|

var TestCardJSON = &Card{
	Number:        "5304259906522887",
	ExpYear:       22,
	ExpMonth:      11,
	CVC:           "749",
	Test:          true,
	Virtual:       false,
	VirtualNumber: "4999993986001010",
}
var TCSJSON = NewCompanyServiceUsingJSON(
	"053128",
	"225",
	"",
)

// func Test_CompanyServiceJSON_CreateAVirtualCard(t *testing.T) {

// 	virtualCardResponse := TCSJSON.CreateAVirtualCard(TestCardJSON, "CardholderInitiatedCredentialOnFile", "ECommerceWithCvc", "")
// 	log.Println("Test output!")
// 	log.Println(virtualCardResponse)
// }
// func Test_CompanyServiceJSON_UpdateAVirtualCardsExpirationDate(t *testing.T) {

// 	VirtualCardExpirationUpdateResponse := TCSJSON.UpdateAVirtualCardsExpirationDate(TestCardJSON, "ECommerceWithCvc")
// 	log.Println("Test output!")
// 	log.Println(VirtualCardExpirationUpdateResponse)
// }
