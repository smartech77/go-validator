package valitor

import (
	"log"
	"testing"
)

// 5304259906522887 2211 749
// 5304259909334470 2211 813
// 5304259902386667 2211 376|

var TestCardXML = &Card{
	Number:        "5304259906522887",
	ExpYear:       22,
	ExpMonth:      11,
	CVC:           "749",
	Test:          true,
	Virtual:       false,
	VirtualNumber: "4999993986001010",
}
var TCSXML = NewCompanyServiceUsingXML(
	"Valitortestfyrirtgr",
	"testadgfyrirgr2010",
	"053128",
	"5006830589",
	"225",
	"https://api.processing.uat.valitor.com/Fyrirtaekjagreidslur/Fyrirtaekjagreidslur.asmx",
)
var VCAuth VirtualCardAuthorization

func Test_CompanyService_GetVirtualNumber(t *testing.T) {

	xmlResponse := TCSXML.GetVirtualNumber(TestCardXML)
	if xmlResponse.ErrorCode != 0 || xmlResponse.SystemError != nil {
		t.Fatalf("Could not get a virtual number for card: " + TestCardXML.Number)
		return
	}

	log.Println(xmlResponse.VirtualNumber)
	// Output
	// 4999993986001010

}

func Test_CompanyService_GetAuthorizationUsingAVirtualCard(t *testing.T) {
	xmlResponse := TCSXML.GetAuthorizationUsingAVirtualCard(TestCardXML, "100", "ISK")
	if xmlResponse.ErrorCode != 0 || xmlResponse.SystemError != nil {
		t.Fatalf("Could not get Authorization for a virtual card: " + TestCardXML.VirtualNumber)
		return
	}
	// set the global auth
	if xmlResponse.ErrorCode == 0 {
		VCAuth = xmlResponse
	}
	jsonReceipt, err := xmlResponse.ReceiptToJSON()
	if err != nil {
		panic(err)
	}

	log.Println(string(jsonReceipt))

}

func Test_CompanyService_Refund(t *testing.T) {
	xmlResponse := TCSXML.Refund(TestCardXML, "100", "ISK")
	if xmlResponse.ErrorCode != 0 || xmlResponse.SystemError != nil {
		t.Fatalf("Could not refund to a virtual card: " + TestCardXML.VirtualNumber)
		return
	}
	jsonReceipt, err := xmlResponse.ReceiptToJSON()
	if err != nil {
		panic(err)
	}
	log.Println(string(jsonReceipt))

}

func Test_CompanyService_InvalidateAuthorization(t *testing.T) {
	xmlResponse := TCSXML.InvalidateAuthorization(TestCardXML, "ISK", "authNumber.")
	if xmlResponse.ErrorCode != 0 || xmlResponse.SystemError != nil {
		t.Fatalf("Could not invalidate authorization number: " + "authnumber.")
		return
	}
	jsonReceipt, err := xmlResponse.ReceiptToJSON()
	if err != nil {
		panic(err)
	}
	log.Println(string(jsonReceipt))

}

func Test_CompanyService_UpdateCardExpirationDate(t *testing.T) {
	TestCardXML.ExpMonth = 13
	TestCardXML.ExpYear = 24
	xmlResponse := TCSXML.UpdateCardExpirationDate(TestCardXML)
	if xmlResponse.ErrorCode != 0 || xmlResponse.SystemError != nil {
		t.Fatalf("Could not update card expiration date")
		return
	}

	log.Println(xmlResponse)

}
