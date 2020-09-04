package test

import (
	"log"
	"testing"

	valitor "github.com/opensourcez/go-valitor"
	helpers "github.com/opensourcez/go-valitor/helpers"
	xmlcore "github.com/opensourcez/go-valitor/xmlcore"
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
var VCAuth xmlcore.FaHeimild
var VCAuthWithoutPayment xmlcore.FaAdeinsHeimild
var hasBeenInitialized bool

func InitTestThings() {
	if !hasBeenInitialized {
		helpers.DebugMode = true
		hasBeenInitialized = true
	}
}
func Test_CompanyService_FaSyndarkortnumer(t *testing.T) {
	InitTestThings()
	xmlResponse := CS.FaSyndarkortnumer(TestCard)
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
	errorResponse := CS.FaSyndarkortnumer(FaultyCard)
	if errorResponse.SystemError.Error() != ExpectedErrorMessage {
		t.Log("Expected:", ExpectedErrorMessage, " got ", errorResponse.SystemError)
		t.Fatal()
		return
	}

	FaultyCard.Number = "5304259906522887"

	// break expiration
	FaultyCard.ExpYear = 0
	ExpectedErrorMessage = "Expiration Year missing"
	errorResponse = CS.FaSyndarkortnumer(FaultyCard)
	if errorResponse.SystemError.Error() != ExpectedErrorMessage {
		t.Log("Expected:", ExpectedErrorMessage, " got ", errorResponse.SystemError)
		t.Fatal()
		return
	}

	FaultyCard.ExpYear = 22

	FaultyCard.ExpMonth = 0
	ExpectedErrorMessage = "Expiration Month missing"
	errorResponse = CS.FaSyndarkortnumer(FaultyCard)
	if errorResponse.SystemError.Error() != ExpectedErrorMessage {
		t.Log("Expected:", ExpectedErrorMessage, " got ", errorResponse.SystemError)
		t.Fatal()
		return
	}

	FaultyCard.ExpYear = 0
	ExpectedErrorMessage = "Expiration Month and Year missing"
	errorResponse = CS.FaSyndarkortnumer(FaultyCard)
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
	errorResponse = CS.FaSyndarkortnumer(FaultyCard)
	if errorResponse.SystemError.Error() != ExpectedErrorMessage {
		t.Log("Expected:", ExpectedErrorMessage, " got ", errorResponse.SystemError)
		t.Fatal()
		return
	}

	// proceed to verify the rest of our test methods.
	TestCard.VirtualNumber = xmlResponse.VirtualNumber
	t.Log("Got virtual number :", xmlResponse.VirtualNumber, "running more test...")
	t.Run("VirtualCard=GetAuthorization", CompanyService_FaHeimild)
	t.Run("VirtualCard=FaEndurgreitt", CompanyService_FaEndurgreitt)
	t.Run("VirtualCard=GetAuthorization", CompanyService_FaHeimild)
	t.Run("VirtualCard=FaOgildingu", CompanyService_FaOgildingu)
	t.Run("VirtualCard=UpdateExpiration", CompanyService_UppfaeraGildistima)
	t.Run("VirtualCard=FaAdeinsHeimild", CompanyService_FaAdeinsHeimild)
	t.Run("VirtualCard=NotaAdeinsheimild", CompanyService_NotaAdeinsheimild)
	t.Run("VirtualCard=FaSidustuFjoraIKortnumeriUtFraSyndarkortnumeri", CompanyService_FaSidustuFjoraIKortnumeriUtFraSyndarkortnumeri)

}

func CompanyService_FaHeimild(t *testing.T) {
	InitTestThings()
	xmlResponse := CS.FaHeimild(TestCard, "100", "ISK")
	if xmlResponse.ErrorCode != 0 || xmlResponse.SystemError != nil {
		t.Log("Could not get Authorization: " + TestCard.VirtualNumber)
		t.Log("System Error:", xmlResponse.SystemError)
		t.Log("Full Error:", xmlResponse)
		t.Fatal()
		return
	}
	VCAuth = xmlResponse
	jsonReceipt, err := xmlResponse.Receipt.ToJSON()
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
	newResponse := CS.FaHeimild(FaultyCard, "100", "ISK")
	if newResponse.SystemError.Error() != ExpectedErrorMessage {
		t.Log("Expected:", ExpectedErrorMessage, " got ", newResponse.SystemError)
		t.Fatal()
		return
	}

	FaultyCard.VirtualNumber = TestCard.VirtualNumber
	ExpectedErrorMessage = "Currency missing"
	newResponse = CS.FaHeimild(FaultyCard, "100", "")
	if newResponse.SystemError.Error() != ExpectedErrorMessage {
		t.Log("Expected:", ExpectedErrorMessage, " got ", newResponse.SystemError)
		t.Fatal()
		return
	}

	ExpectedErrorMessage = "Amount missing"
	newResponse = CS.FaHeimild(FaultyCard, "", "ISK")
	if newResponse.SystemError.Error() != ExpectedErrorMessage {
		t.Log("Expected:", ExpectedErrorMessage, " got ", newResponse.SystemError)
		t.Fatal()
		return
	}

}

func CompanyService_FaEndurgreitt(t *testing.T) {
	InitTestThings()
	xmlResponse := CS.FaEndurgreitt(TestCard, "100", "ISK")
	if xmlResponse.ErrorCode != 0 || xmlResponse.SystemError != nil {
		t.Log("Could not FaEndurgreitt: "+TestCard.VirtualNumber, " error:", xmlResponse)
		t.Log("System Error:", xmlResponse.SystemError)
		t.Log("Full Error:", xmlResponse)
		t.Fatal()
		return
	}
	jsonReceipt, err := xmlResponse.Receipt.ToJSON()
	if err != nil {

		panic(err)
	}
	log.Println("Got a FaEndurgreitt receipt:", string(jsonReceipt))

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
	newResponse := CS.FaEndurgreitt(FaultyCard, "100", "ISK")
	if newResponse.SystemError.Error() != ExpectedErrorMessage {
		t.Log("Expected:", ExpectedErrorMessage, " got ", newResponse.SystemError)
		t.Fatal()
		return
	}

	FaultyCard.VirtualNumber = TestCard.VirtualNumber
	ExpectedErrorMessage = "Currency missing"
	newResponse = CS.FaEndurgreitt(FaultyCard, "100", "")
	if newResponse.SystemError.Error() != ExpectedErrorMessage {
		t.Log("Expected:", ExpectedErrorMessage, " got ", newResponse.SystemError)
		t.Fatal()
		return
	}

	ExpectedErrorMessage = "Amount missing"
	newResponse = CS.FaEndurgreitt(FaultyCard, "", "ISK")
	if newResponse.SystemError.Error() != ExpectedErrorMessage {
		t.Log("Expected:", ExpectedErrorMessage, " got ", newResponse.SystemError)
		t.Fatal()
		return
	}
}

func CompanyService_FaOgildingu(t *testing.T) {
	InitTestThings()
	xmlResponse := CS.FaOgildingu(TestCard, "ISK", VCAuth.Receipt.TransactionID)
	if xmlResponse.ErrorCode != 0 || xmlResponse.SystemError != nil {
		t.Log("Could not invalidate authorization number: " + VCAuth.Receipt.TransactionID)
		t.Log("System Error:", xmlResponse.SystemError)
		t.Log("Full Error:", xmlResponse)
		t.Fatal()
		return
	}
	jsonReceipt, err := xmlResponse.Receipt.ToJSON()
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
	newResponse := CS.FaOgildingu(FaultyCard, "100", "ISK")
	if newResponse.SystemError.Error() != ExpectedErrorMessage {
		t.Log("Expected:", ExpectedErrorMessage, " got ", newResponse.SystemError)
		t.Fatal()
		return
	}

	FaultyCard.VirtualNumber = TestCard.VirtualNumber
	ExpectedErrorMessage = "Currency missing"
	newResponse = CS.FaOgildingu(FaultyCard, "", "randomnumber")
	if newResponse.SystemError.Error() != ExpectedErrorMessage {
		t.Log("Expected:", ExpectedErrorMessage, " got ", newResponse.SystemError)
		t.Fatal()
		return
	}

	ExpectedErrorMessage = "Authorization number missing"
	newResponse = CS.FaOgildingu(FaultyCard, "ISK", "")
	if newResponse.SystemError.Error() != ExpectedErrorMessage {
		t.Log("Expected:", ExpectedErrorMessage, " got ", newResponse.SystemError)
		t.Fatal()
		return
	}
}

func CompanyService_UppfaeraGildistima(t *testing.T) {
	InitTestThings()
	TestCard.ExpMonth = 12
	TestCard.ExpYear = 30
	xmlResponse := CS.UppfaeraGildistima(TestCard)
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
	newResponse := CS.UppfaeraGildistima(FaultyCard)
	if newResponse.SystemError.Error() != ExpectedErrorMessage {
		t.Log("Expected:", ExpectedErrorMessage, " got ", newResponse.SystemError)
		t.Fatal()
		return
	}

	FaultyCard.VirtualNumber = TestCard.VirtualNumber
	// break expiration
	FaultyCard.ExpYear = 0
	ExpectedErrorMessage = "Expiration Year missing"
	newResponse = CS.UppfaeraGildistima(FaultyCard)
	if newResponse.SystemError.Error() != ExpectedErrorMessage {
		t.Log("Expected:", ExpectedErrorMessage, " got ", newResponse.SystemError)
		t.Fatal()
		return
	}

	FaultyCard.ExpYear = 22

	FaultyCard.ExpMonth = 0
	ExpectedErrorMessage = "Expiration Month missing"
	newResponse = CS.UppfaeraGildistima(FaultyCard)
	if newResponse.SystemError.Error() != ExpectedErrorMessage {
		t.Log("Expected:", ExpectedErrorMessage, " got ", newResponse.SystemError)
		t.Fatal()
		return
	}

	FaultyCard.ExpYear = 0
	ExpectedErrorMessage = "Expiration Month and Year missing"
	newResponse = CS.UppfaeraGildistima(FaultyCard)
	if newResponse.SystemError.Error() != ExpectedErrorMessage {
		t.Log("Expected:", ExpectedErrorMessage, " got ", newResponse.SystemError)
		t.Fatal()
		return
	}

}

func CompanyService_FaAdeinsHeimild(t *testing.T) {
	InitTestThings()
	xmlResponse := CS.FaAdeinsHeimild(TestCard, "100", "ISK")
	if xmlResponse.ErrorCode != 0 || xmlResponse.SystemError != nil {
		t.Log("Could not get Authorization (without payment): " + TestCard.VirtualNumber)
		t.Log("System Error:", xmlResponse.SystemError)
		t.Log("Full Error:", xmlResponse)
		t.Fatal()
		return
	}
	VCAuthWithoutPayment = xmlResponse
	jsonReceipt, err := xmlResponse.Receipt.ToJSON()
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
	newResponse := CS.FaAdeinsHeimild(FaultyCard, "100", "ISK")
	if newResponse.SystemError.Error() != ExpectedErrorMessage {
		t.Log("Expected:", ExpectedErrorMessage, " got ", newResponse.SystemError)
		t.Fatal()
		return
	}

	FaultyCard.VirtualNumber = TestCard.VirtualNumber
	ExpectedErrorMessage = "Currency missing"
	newResponse = CS.FaAdeinsHeimild(FaultyCard, "100", "")
	if newResponse.SystemError.Error() != ExpectedErrorMessage {
		t.Log("Expected:", ExpectedErrorMessage, " got ", newResponse.SystemError)
		t.Fatal()
		return
	}

	ExpectedErrorMessage = "Amount missing"
	newResponse = CS.FaAdeinsHeimild(FaultyCard, "", "ISK")
	if newResponse.SystemError.Error() != ExpectedErrorMessage {
		t.Log("Expected:", ExpectedErrorMessage, " got ", newResponse.SystemError)
		t.Fatal()
		return
	}

	FaultyCard.CVC = ""
	ExpectedErrorMessage = "CVC missing"
	newResponse = CS.FaAdeinsHeimild(FaultyCard, "100", "ISK")
	if newResponse.SystemError.Error() != ExpectedErrorMessage {
		t.Log("Expected:", ExpectedErrorMessage, " got ", newResponse.SystemError)
		t.Fatal()
		return
	}
}

func CompanyService_NotaAdeinsheimild(t *testing.T) {
	InitTestThings()
	xmlResponse := CS.NotaAdeinsheimild(TestCard, VCAuthWithoutPayment.Receipt.TransactionID)
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
	newResponse := CS.NotaAdeinsheimild(FaultyCard, "randomnumber")
	if newResponse.SystemError.Error() != ExpectedErrorMessage {
		t.Log("Expected:", ExpectedErrorMessage, " got ", newResponse.SystemError)
		t.Fatal()
		return
	}

	FaultyCard.VirtualNumber = TestCard.VirtualNumber
	ExpectedErrorMessage = "Authorization number missing"
	newResponse = CS.NotaAdeinsheimild(FaultyCard, "")
	if newResponse.SystemError.Error() != ExpectedErrorMessage {
		t.Log("Expected:", ExpectedErrorMessage, " got ", newResponse.SystemError)
		t.Fatal()
		return
	}

	FaultyCard.CVC = ""
	ExpectedErrorMessage = "CVC missing"
	newResponse = CS.NotaAdeinsheimild(FaultyCard, "randomnumber")
	if newResponse.SystemError.Error() != ExpectedErrorMessage {
		t.Log("Expected:", ExpectedErrorMessage, " got ", newResponse.SystemError)
		t.Fatal()
		return
	}

}

func CompanyService_FaSidustuFjoraIKortnumeriUtFraSyndarkortnumeri(t *testing.T) {
	InitTestThings()
	xmlResponse := CS.FaSidustuFjoraIKortnumeriUtFraSyndarkortnumeri(TestCard)
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
		newResponse := CS.NotaAdeinsheimild(FaultyCard, "randomnumber")
		if newResponse.SystemError.Error() != ExpectedErrorMessage {
			t.Log("Expected:", ExpectedErrorMessage, " got ", newResponse.SystemError)
			t.Fatal()
			return
		}

	}

}
