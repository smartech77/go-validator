package valitor

import (
	"log"
	"testing"
)

// 5304259906522887 2211 749
// 5304259909334470 2211 813
// 5304259902386667 2211 376|

var TestCardXML = &Card{
	Number:   "5304259906522887",
	ExpYear:  22,
	ExpMonth: 11,
	CVC:      "749",
	Test:     true,
	Virtual:  false,
	// 5999993615731195
	VirtualNumber: "",
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
var VCAuthWithoutPayment VirtualCardAuthorizationWithoutPayment
var hasBeenInitialized bool

func InitTestThings() {
	if !hasBeenInitialized {
		DebugMode = true
		hasBeenInitialized = true
	}
}
func Test_CompanyServiceXML_GetVirtualNumber(t *testing.T) {
	InitTestThings()
	xmlResponse := TCSXML.GetVirtualNumber(TestCardXML)
	if xmlResponse.ErrorCode != 0 || xmlResponse.SystemError != nil {
		t.Log("Could not get a virtual number for card: "+TestCardXML.Number, " .. not running more tests ...")
		t.Log("System Error:", xmlResponse.SystemError)
		t.Log("Full Error:", xmlResponse)
		t.Fatal()
		return
	}

	TestCardXML.VirtualNumber = xmlResponse.VirtualNumber
	t.Log("Got virtual number :", xmlResponse.VirtualNumber, "running more test...")

	// t.Run("VirtualCard=GetAuthorization", CompanyServiceXML_GetAuthorization)
	// t.Run("VirtualCard=Refund", CompanyServiceXML_Refund)
	// t.Run("VirtualCard=GetAuthorization", CompanyServiceXML_GetAuthorization)
	// t.Run("VirtualCard=InvalidateAuthorization", CompanyServiceXML_InvalidateAuthorization)
	// t.Run("VirtualCard=UpdateExpiration", CompanyServiceXML_UpdateCardExpirationDate)Card
	t.Run("VirtualCard=GetAuthorizationWithoutPayment", CompanyServiceXML_GetAuthorizationWithoutPayment)
	t.Run("VirtualCard=UseAuthorization", CompanyServiceXML_UseAuthorization)

}

func CompanyServiceXML_GetAuthorization(t *testing.T) {
	InitTestThings()
	xmlResponse := TCSXML.GetAuthorization(TestCardXML, "100", "ISK")
	if xmlResponse.ErrorCode != 0 || xmlResponse.SystemError != nil {
		t.Log("Could not get Authorization: " + TestCardXML.VirtualNumber)
		t.Log("System Error:", xmlResponse.SystemError)
		t.Log("Full Error:", xmlResponse)
		t.Fatal()
		return
	}
	VCAuth = xmlResponse
	jsonReceipt, err := xmlResponse.Receipt.ReceiptToJSON()
	if err != nil {
		panic(err)
	}

	t.Log("Got an authorization receipt:", string(jsonReceipt))

}

func CompanyServiceXML_Refund(t *testing.T) {
	InitTestThings()
	xmlResponse := TCSXML.Refund(TestCardXML, "100", "ISK")
	if xmlResponse.ErrorCode != 0 || xmlResponse.SystemError != nil {
		t.Log("Could not refund: "+TestCardXML.VirtualNumber, " error:", xmlResponse)
		t.Log("System Error:", xmlResponse.SystemError)
		t.Log("Full Error:", xmlResponse)
		t.Fatal()
		return
	}
	jsonReceipt, err := xmlResponse.Receipt.ReceiptToJSON()
	if err != nil {

		panic(err)
	}
	log.Println("Got a refund receipt:", string(jsonReceipt))

}

func CompanyServiceXML_InvalidateAuthorization(t *testing.T) {
	InitTestThings()
	xmlResponse := TCSXML.InvalidateAuthorization(TestCardXML, "ISK", VCAuth.Receipt.TransactionID)
	if xmlResponse.ErrorCode != 0 || xmlResponse.SystemError != nil {
		t.Log("Could not invalidate authorization number: " + VCAuth.Receipt.TransactionID)
		t.Log("System Error:", xmlResponse.SystemError)
		t.Log("Full Error:", xmlResponse)
		t.Fatal()
		return
	}
	jsonReceipt, err := xmlResponse.Receipt.ReceiptToJSON()
	if err != nil {
		panic(err)
	}
	t.Log("Got an invalidation receipt:", string(jsonReceipt))

}

func CompanyServiceXML_UpdateCardExpirationDate(t *testing.T) {
	InitTestThings()
	TestCardXML.ExpMonth = 12
	TestCardXML.ExpYear = 30
	xmlResponse := TCSXML.UpdateCardExpirationDate(TestCardXML)
	if xmlResponse.ErrorCode != 0 || xmlResponse.SystemError != nil {
		t.Log("Could not update card expiration date")
		t.Log("System Error:", xmlResponse.SystemError)
		t.Log("Full Error:", xmlResponse)
		t.Fatal()
		return
	}
	t.Log("Ssuccessfully update the virtual card")
}

func CompanyServiceXML_GetAuthorizationWithoutPayment(t *testing.T) {
	InitTestThings()
	xmlResponse := TCSXML.GetAuthorizationWithoutPayment(TestCardXML, "100", "ISK")
	if xmlResponse.ErrorCode != 0 || xmlResponse.SystemError != nil {
		t.Log("Could not get Authorization (without payment): " + TestCardXML.VirtualNumber)
		t.Log("System Error:", xmlResponse.SystemError)
		t.Log("Full Error:", xmlResponse)
		t.Fatal()
		return
	}
	VCAuthWithoutPayment = xmlResponse
	jsonReceipt, err := xmlResponse.Receipt.ReceiptToJSON()
	if err != nil {
		panic(err)
	}

	t.Log("Got an authorization (without payment) receipt:", string(jsonReceipt))

}

func CompanyServiceXML_UseAuthorization(t *testing.T) {
	InitTestThings()
	xmlResponse := TCSXML.UseAuthorization(TestCardXML, VCAuthWithoutPayment.Receipt.TransactionID)
	if xmlResponse.ErrorCode != 0 || xmlResponse.SystemError != nil {
		t.Log("Could not use Authorization: " + TestCardXML.VirtualNumber)
		t.Log("System Error:", xmlResponse.SystemError)
		t.Log("Full Error:", xmlResponse)
		t.Fatal()
		return
	}

	t.Log("Succsessfully used authorization")

}
