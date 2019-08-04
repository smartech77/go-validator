package valitor

import (
	"log"
	"testing"
)

// 5304259906522887 2211 749
// 5304259909334470 2211 813
// 5304259902386667 2211 376|

var TestCard1 = &Card{
	Number:        "5304259906522887",
	ExpYear:       "22",
	ExpMonth:      "11",
	CVC:           "749",
	Test:          true,
	Virtual:       false,
	VirtualNumber: "4999993986001010",
}
var TCS = NewCompanyService(
	"Valitortestfyrirtgr",
	"testadgfyrirgr2010",
	"053128",
	"5006830589",
	"225",
	"https://api.processing.uat.valitor.com/Fyrirtaekjagreidslur/Fyrirtaekjagreidslur.asmx",
)

func Test_CompanyService_GetVirtualNumber(t *testing.T) {

	xmlResponse := TCS.GetVirtualNumber(TestCard1)
	if xmlResponse.ErrorCode != 0 || xmlResponse.SystemError != nil {
		t.Fatalf("Could not get a virtual number for card:" + TestCard1.Number)
		return
	}

	// Next Step
	t.Run("VirtualCard=GetAuthorization", Test_CompanyService_GetAuthorizationUsingAVirtualCard)

	log.Println(xmlResponse.VirtualNumber)
	// Output
	// 4999993986001010

}

func Test_CompanyService_GetAuthorizationUsingAVirtualCard(t *testing.T) {
	xmlResponse := TCS.GetAuthorizationUsingAVirtualCard(TestCard1, "100", "ISK")
	if xmlResponse.ErrorCode != 0 || xmlResponse.SystemError != nil {
		t.Fatalf("Could not get Authorization for a virtual card:" + TestCard1.VirtualNumber)
		return
	}
	jsonReceipt, err := xmlResponse.ReceiptToJSON()
	if err != nil {
		panic(err)
	}

	// Next Step
	t.Run("VirtualCard=Refund", Test_CompanyService_GetAuthorizationUsingAVirtualCard)

	log.Println(string(jsonReceipt))

}

func Test_CompanyService_Refund(t *testing.T) {
	xmlResponse := TCS.Refund(TestCard1, "100", "ISK")
	if xmlResponse.ErrorCode != 0 || xmlResponse.SystemError != nil {
		t.Fatalf("Could not refund to a virtual card:" + TestCard1.VirtualNumber)
		return
	}
	jsonReceipt, err := xmlResponse.ReceiptToJSON()
	if err != nil {
		panic(err)
	}
	log.Println(string(jsonReceipt))

}
