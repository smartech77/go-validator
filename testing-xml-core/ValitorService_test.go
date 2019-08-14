package test

import (
	"log"
	"testing"

	valitor "github.com/zkynetio/go-valitor"
	helpers "github.com/zkynetio/go-valitor/helpers"
	xmlcore "github.com/zkynetio/go-valitor/xmlcore"
)

var TestCard = &xmlcore.Card{
	Number:   "5304259906522887",
	ExpYear:  22,
	ExpMonth: 11,
	CVC:      "749",
	// 5999993615731195
	VirtualNumber: "",
}
var CS = valitor.NewValitorService(
	"Valitortestfyrirtgr",
	"testadgfyrirgr2010",
	"053128",
	"5006830589",
	"225",
	"https://api.processing.uat.valitor.com/Fyrirtaekjagreidslur/Fyrirtaekjagreidslur.asmx",
)
var VCAuth xmlcore.VirtualCardAuthorization
var VCAuthWithoutPayment xmlcore.VirtualCardAuthorizationWithoutPayment
var hasBeenInitialized bool

func InitTestThings() {
	if !hasBeenInitialized {
		helpers.DebugMode = true
		hasBeenInitialized = true
	}
}
func Test_CompanyService_GetVirtualNumber(t *testing.T) {
	InitTestThings()
	xmlResponse := CS.GetVirtualNumber(TestCard)
	if xmlResponse.ErrorCode != 0 || xmlResponse.SystemError != nil {
		t.Log("Could not get a virtual number for card: "+TestCard.Number, " .. not running more tests ...")
		t.Log("System Error:", xmlResponse.SystemError)
		t.Log("Full Error:", xmlResponse)
		t.Fatal()
		return
	}

	FaultyCard := &xmlcore.Card{
		Number:   "5304259906522887",
		ExpYear:  22,
		ExpMonth: 11,
		CVC:      "749",
		// 5999993615731195
		VirtualNumber: "",
	}

	// Break card number
	FaultyCard.Number = ""
	ExpectedErrorMessage := "Card Number missing"
	errorResponse := CS.GetVirtualNumber(FaultyCard)
	if errorResponse.SystemError.Error() != ExpectedErrorMessage {
		t.Log("Expected:", ExpectedErrorMessage, " got ", errorResponse.SystemError)
		t.Fatal()
		return
	}

	FaultyCard.Number = "5304259906522887"

	// break expiration
	FaultyCard.ExpYear = 0
	ExpectedErrorMessage = "Expiration Year missing"
	errorResponse = CS.GetVirtualNumber(FaultyCard)
	if errorResponse.SystemError.Error() != ExpectedErrorMessage {
		t.Log("Expected:", ExpectedErrorMessage, " got ", errorResponse.SystemError)
		t.Fatal()
		return
	}

	FaultyCard.ExpYear = 22

	FaultyCard.ExpMonth = 0
	ExpectedErrorMessage = "Expiration Month missing"
	errorResponse = CS.GetVirtualNumber(FaultyCard)
	if errorResponse.SystemError.Error() != ExpectedErrorMessage {
		t.Log("Expected:", ExpectedErrorMessage, " got ", errorResponse.SystemError)
		t.Fatal()
		return
	}

	FaultyCard.ExpYear = 0
	ExpectedErrorMessage = "Expiration Month and Year missing"
	errorResponse = CS.GetVirtualNumber(FaultyCard)
	if errorResponse.SystemError.Error() != ExpectedErrorMessage {
		t.Log("Expected:", ExpectedErrorMessage, " got ", errorResponse.SystemError)
		t.Fatal()
		return
	}

	FaultyCard.ExpYear = 22
	FaultyCard.ExpMonth = 11

	// break CVC
	FaultyCard.CVC = ""
	ExpectedErrorMessage = "CVC missing"
	errorResponse = CS.GetVirtualNumber(FaultyCard)
	if errorResponse.SystemError.Error() != ExpectedErrorMessage {
		t.Log("Expected:", ExpectedErrorMessage, " got ", errorResponse.SystemError)
		t.Fatal()
		return
	}

	// proceed to verify the rest of our test methods.
	TestCard.VirtualNumber = xmlResponse.VirtualNumber
	t.Log("Got virtual number :", xmlResponse.VirtualNumber, "running more test...")
	t.Run("VirtualCard=GetAuthorization", CompanyService_GetAuthorization)
	t.Run("VirtualCard=Refund", CompanyService_Refund)
	t.Run("VirtualCard=GetAuthorization", CompanyService_GetAuthorization)
	t.Run("VirtualCard=InvalidateAuthorization", CompanyService_InvalidateAuthorization)
	t.Run("VirtualCard=UpdateExpiration", CompanyService_UpdateCardExpirationDate)
	t.Run("VirtualCard=GetAuthorizationWithoutPayment", CompanyService_GetAuthorizationWithoutPayment)
	t.Run("VirtualCard=UseAuthorization", CompanyService_UseAuthorization)
	t.Run("VirtualCard=GetLastFourDigitsFromTheRealCard", CompanyService_GetLastFourDigitsFromTheRealCard)

}

func CompanyService_GetAuthorization(t *testing.T) {
	InitTestThings()
	xmlResponse := CS.GetAuthorization(TestCard, "100", "ISK")
	if xmlResponse.ErrorCode != 0 || xmlResponse.SystemError != nil {
		t.Log("Could not get Authorization: " + TestCard.VirtualNumber)
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

	FaultyCard := &xmlcore.Card{
		Number:   "5304259906522887",
		ExpYear:  22,
		ExpMonth: 11,
		CVC:      "749",
		// 5999993615731195
		VirtualNumber: "",
	}
	FaultyCard.VirtualNumber = ""
	ExpectedErrorMessage := "Virtual Number missing"
	newResponse := CS.GetAuthorization(FaultyCard, "100", "ISK")
	if newResponse.SystemError.Error() != ExpectedErrorMessage {
		t.Log("Expected:", ExpectedErrorMessage, " got ", newResponse.SystemError)
		t.Fatal()
		return
	}

	FaultyCard.VirtualNumber = TestCard.VirtualNumber
	ExpectedErrorMessage = "Currency missing"
	newResponse = CS.GetAuthorization(FaultyCard, "100", "")
	if newResponse.SystemError.Error() != ExpectedErrorMessage {
		t.Log("Expected:", ExpectedErrorMessage, " got ", newResponse.SystemError)
		t.Fatal()
		return
	}

	ExpectedErrorMessage = "Amount missing"
	newResponse = CS.GetAuthorization(FaultyCard, "", "ISK")
	if newResponse.SystemError.Error() != ExpectedErrorMessage {
		t.Log("Expected:", ExpectedErrorMessage, " got ", newResponse.SystemError)
		t.Fatal()
		return
	}

}

func CompanyService_Refund(t *testing.T) {
	InitTestThings()
	xmlResponse := CS.Refund(TestCard, "100", "ISK")
	if xmlResponse.ErrorCode != 0 || xmlResponse.SystemError != nil {
		t.Log("Could not refund: "+TestCard.VirtualNumber, " error:", xmlResponse)
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

	FaultyCard := &xmlcore.Card{
		Number:   "5304259906522887",
		ExpYear:  22,
		ExpMonth: 11,
		CVC:      "749",
		// 5999993615731195
		VirtualNumber: "",
	}
	FaultyCard.VirtualNumber = ""
	ExpectedErrorMessage := "Virtual Number missing"
	newResponse := CS.Refund(FaultyCard, "100", "ISK")
	if newResponse.SystemError.Error() != ExpectedErrorMessage {
		t.Log("Expected:", ExpectedErrorMessage, " got ", newResponse.SystemError)
		t.Fatal()
		return
	}

	FaultyCard.VirtualNumber = TestCard.VirtualNumber
	ExpectedErrorMessage = "Currency missing"
	newResponse = CS.Refund(FaultyCard, "100", "")
	if newResponse.SystemError.Error() != ExpectedErrorMessage {
		t.Log("Expected:", ExpectedErrorMessage, " got ", newResponse.SystemError)
		t.Fatal()
		return
	}

	ExpectedErrorMessage = "Amount missing"
	newResponse = CS.Refund(FaultyCard, "", "ISK")
	if newResponse.SystemError.Error() != ExpectedErrorMessage {
		t.Log("Expected:", ExpectedErrorMessage, " got ", newResponse.SystemError)
		t.Fatal()
		return
	}
}

func CompanyService_InvalidateAuthorization(t *testing.T) {
	InitTestThings()
	xmlResponse := CS.InvalidateAuthorization(TestCard, "ISK", VCAuth.Receipt.TransactionID)
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

	FaultyCard := &xmlcore.Card{
		Number:   "5304259906522887",
		ExpYear:  22,
		ExpMonth: 11,
		CVC:      "749",
		// 5999993615731195
		VirtualNumber: "",
	}
	FaultyCard.VirtualNumber = ""
	ExpectedErrorMessage := "Virtual Number missing"
	newResponse := CS.InvalidateAuthorization(FaultyCard, "100", "ISK")
	if newResponse.SystemError.Error() != ExpectedErrorMessage {
		t.Log("Expected:", ExpectedErrorMessage, " got ", newResponse.SystemError)
		t.Fatal()
		return
	}

	FaultyCard.VirtualNumber = TestCard.VirtualNumber
	ExpectedErrorMessage = "Currency missing"
	newResponse = CS.InvalidateAuthorization(FaultyCard, "", "randomnumber")
	if newResponse.SystemError.Error() != ExpectedErrorMessage {
		t.Log("Expected:", ExpectedErrorMessage, " got ", newResponse.SystemError)
		t.Fatal()
		return
	}

	ExpectedErrorMessage = "Authorization number missing"
	newResponse = CS.InvalidateAuthorization(FaultyCard, "ISK", "")
	if newResponse.SystemError.Error() != ExpectedErrorMessage {
		t.Log("Expected:", ExpectedErrorMessage, " got ", newResponse.SystemError)
		t.Fatal()
		return
	}
}

func CompanyService_UpdateCardExpirationDate(t *testing.T) {
	InitTestThings()
	TestCard.ExpMonth = 12
	TestCard.ExpYear = 30
	xmlResponse := CS.UpdateCardExpirationDate(TestCard)
	if xmlResponse.ErrorCode != 0 || xmlResponse.SystemError != nil {
		t.Log("Could not update card expiration date")
		t.Log("System Error:", xmlResponse.SystemError)
		t.Log("Full Error:", xmlResponse)
		t.Fatal()
		return
	}
	t.Log("Ssuccessfully update the virtual card")

	FaultyCard := &xmlcore.Card{
		Number:   "5304259906522887",
		ExpYear:  22,
		ExpMonth: 11,
		CVC:      "749",
		// 5999993615731195
		VirtualNumber: "",
	}
	FaultyCard.VirtualNumber = ""
	ExpectedErrorMessage := "Virtual Number missing"
	newResponse := CS.UpdateCardExpirationDate(FaultyCard)
	if newResponse.SystemError.Error() != ExpectedErrorMessage {
		t.Log("Expected:", ExpectedErrorMessage, " got ", newResponse.SystemError)
		t.Fatal()
		return
	}

	FaultyCard.VirtualNumber = TestCard.VirtualNumber
	// break expiration
	FaultyCard.ExpYear = 0
	ExpectedErrorMessage = "Expiration Year missing"
	newResponse = CS.UpdateCardExpirationDate(FaultyCard)
	if newResponse.SystemError.Error() != ExpectedErrorMessage {
		t.Log("Expected:", ExpectedErrorMessage, " got ", newResponse.SystemError)
		t.Fatal()
		return
	}

	FaultyCard.ExpYear = 22

	FaultyCard.ExpMonth = 0
	ExpectedErrorMessage = "Expiration Month missing"
	newResponse = CS.UpdateCardExpirationDate(FaultyCard)
	if newResponse.SystemError.Error() != ExpectedErrorMessage {
		t.Log("Expected:", ExpectedErrorMessage, " got ", newResponse.SystemError)
		t.Fatal()
		return
	}

	FaultyCard.ExpYear = 0
	ExpectedErrorMessage = "Expiration Month and Year missing"
	newResponse = CS.UpdateCardExpirationDate(FaultyCard)
	if newResponse.SystemError.Error() != ExpectedErrorMessage {
		t.Log("Expected:", ExpectedErrorMessage, " got ", newResponse.SystemError)
		t.Fatal()
		return
	}

}

func CompanyService_GetAuthorizationWithoutPayment(t *testing.T) {
	InitTestThings()
	xmlResponse := CS.GetAuthorizationWithoutPayment(TestCard, "100", "ISK")
	if xmlResponse.ErrorCode != 0 || xmlResponse.SystemError != nil {
		t.Log("Could not get Authorization (without payment): " + TestCard.VirtualNumber)
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

	FaultyCard := &xmlcore.Card{
		Number:   "5304259906522887",
		ExpYear:  22,
		ExpMonth: 11,
		CVC:      "749",
		// 5999993615731195
		VirtualNumber: "",
	}
	FaultyCard.VirtualNumber = ""
	ExpectedErrorMessage := "Virtual Number missing"
	newResponse := CS.GetAuthorizationWithoutPayment(FaultyCard, "100", "ISK")
	if newResponse.SystemError.Error() != ExpectedErrorMessage {
		t.Log("Expected:", ExpectedErrorMessage, " got ", newResponse.SystemError)
		t.Fatal()
		return
	}

	FaultyCard.VirtualNumber = TestCard.VirtualNumber
	ExpectedErrorMessage = "Currency missing"
	newResponse = CS.GetAuthorizationWithoutPayment(FaultyCard, "100", "")
	if newResponse.SystemError.Error() != ExpectedErrorMessage {
		t.Log("Expected:", ExpectedErrorMessage, " got ", newResponse.SystemError)
		t.Fatal()
		return
	}

	ExpectedErrorMessage = "Amount missing"
	newResponse = CS.GetAuthorizationWithoutPayment(FaultyCard, "", "ISK")
	if newResponse.SystemError.Error() != ExpectedErrorMessage {
		t.Log("Expected:", ExpectedErrorMessage, " got ", newResponse.SystemError)
		t.Fatal()
		return
	}

	FaultyCard.CVC = ""
	ExpectedErrorMessage = "CVC missing"
	newResponse = CS.GetAuthorizationWithoutPayment(FaultyCard, "100", "ISK")
	if newResponse.SystemError.Error() != ExpectedErrorMessage {
		t.Log("Expected:", ExpectedErrorMessage, " got ", newResponse.SystemError)
		t.Fatal()
		return
	}
}

func CompanyService_UseAuthorization(t *testing.T) {
	InitTestThings()
	xmlResponse := CS.UseAuthorization(TestCard, VCAuthWithoutPayment.Receipt.TransactionID)
	if xmlResponse.ErrorCode != 0 || xmlResponse.SystemError != nil {
		t.Log("Could not use Authorization: " + TestCard.VirtualNumber)
		t.Log("System Error:", xmlResponse.SystemError)
		t.Log("Full Error:", xmlResponse)
		t.Fatal()
		return
	}

	t.Log("Succsessfully used authorization")

	FaultyCard := &xmlcore.Card{
		Number:   "5304259906522887",
		ExpYear:  22,
		ExpMonth: 11,
		CVC:      "749",
		// 5999993615731195
		VirtualNumber: "",
	}
	FaultyCard.VirtualNumber = ""
	ExpectedErrorMessage := "Virtual Number missing"
	newResponse := CS.UseAuthorization(FaultyCard, "randomnumber")
	if newResponse.SystemError.Error() != ExpectedErrorMessage {
		t.Log("Expected:", ExpectedErrorMessage, " got ", newResponse.SystemError)
		t.Fatal()
		return
	}

	FaultyCard.VirtualNumber = TestCard.VirtualNumber
	ExpectedErrorMessage = "Authorization number missing"
	newResponse = CS.UseAuthorization(FaultyCard, "")
	if newResponse.SystemError.Error() != ExpectedErrorMessage {
		t.Log("Expected:", ExpectedErrorMessage, " got ", newResponse.SystemError)
		t.Fatal()
		return
	}

	FaultyCard.CVC = ""
	ExpectedErrorMessage = "CVC missing"
	newResponse = CS.UseAuthorization(FaultyCard, "randomnumber")
	if newResponse.SystemError.Error() != ExpectedErrorMessage {
		t.Log("Expected:", ExpectedErrorMessage, " got ", newResponse.SystemError)
		t.Fatal()
		return
	}

}

func CompanyService_GetLastFourDigitsFromTheRealCard(t *testing.T) {
	InitTestThings()
	xmlResponse := CS.GetLastFourDigitsFromTheRealCard(TestCard)
	if xmlResponse.ErrorCode != 0 || xmlResponse.SystemError != nil {
		t.Log("Could not get last four digits: " + TestCard.VirtualNumber)
		t.Log("System Error:", xmlResponse.SystemError)
		t.Log("Full Error:", xmlResponse)
		t.Fatal()
		return
	}

	if xmlResponse.Kortnumer != TestCard.Number[len(TestCard.Number)-4:] {
		t.Log("If you are seeing this message we recommend investigating it in detail. Not having the last four digits match could indicate a system error in Valitors system.")
		t.Log("Last four from valitor:", xmlResponse.Kortnumer)
		t.Log("Last four from the real card used:", TestCard.Number[len(TestCard.Number)-4:])
		t.Fatal()
	} else {
		t.Log("Got the last four digits:", xmlResponse.Kortnumer)

		FaultyCard := &xmlcore.Card{
			Number:   "5304259906522887",
			ExpYear:  22,
			ExpMonth: 11,
			CVC:      "749",
			// 5999993615731195
			VirtualNumber: "",
		}
		FaultyCard.VirtualNumber = ""
		ExpectedErrorMessage := "Virtual Number missing"
		newResponse := CS.UseAuthorization(FaultyCard, "randomnumber")
		if newResponse.SystemError.Error() != ExpectedErrorMessage {
			t.Log("Expected:", ExpectedErrorMessage, " got ", newResponse.SystemError)
			t.Fatal()
			return
		}

	}

}
