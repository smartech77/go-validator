package xmlcore

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"strconv"
	"strings"
	"sync"

	"github.com/opensourcez/go-valitor/helpers"
)

type Card struct {
	CVC           string
	ExpYear       int
	ExpMonth      int
	Number        string
	VirtualNumber string
	// 3D secure card verification data
	// Specific to the newer JSON API
}

// CompanyService ...
type CompanyService struct {
	Settings *Settings
	Mux      sync.RWMutex
}

// Settings ...
type Settings struct {
	Username               string
	Password               string
	ContractNumber         string
	ContractIdentidyNumber string
	PosID                  string
	URL                    string
}

// VirtualNumber ...
// Documentation: https://specs.valitor.is/CorporatePayments_ISL/Web_Services/#41-fasyndarkortnumer
type FaSyndarkortnumer struct {
	SystemError   error
	ErrorCode     int    `xml:"Body>FaSyndarkortnumerResponse>FaSyndarkortnumerResult>Villunumer"`
	ErrorMessage  string `xml:"Body>FaSyndarkortnumerResponse>FaSyndarkortnumerResult>Villuskilabod"`
	ErrorLogID    string `xml:"Body>FaSyndarkortnumerResponse>FaSyndarkortnumerResult>VilluLogID"`
	VirtualNumber string `xml:"Body>FaSyndarkortnumerResponse>FaSyndarkortnumerResult>Syndarkortnumer"`
}

// GetVirtualNumber ...
// Documentation: https://specs.valitor.is/CorporatePayments_ISL/Web_Services/#41-fasyndarkortnumer
func (cs *CompanyService) FaSyndarkortnumer(card *Card) (response FaSyndarkortnumer) {
	if err := checkCardExpirationDate(card); err != nil {
		response.SystemError = err
		return
	}

	if err := checkCardCVC(card); err != nil {
		response.SystemError = err
		return
	}

	if err := checkCardNumber(card); err != nil {
		response.SystemError = err
		return
	}

	body := `<?xml version="1.0" encoding="utf-8"?> 
		<soap:Envelope xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
		xmlns:xsd="http://www.w3.org/2001/XMLSchema" 
	 	xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
		<soap:Body>
		<FaSyndarkortnumer xmlns="http://api.valitor.is/Fyrirtaekjagreidslur/">
		<Notandanafn>` + cs.Settings.Username + `</Notandanafn>
		<Lykilord>` + cs.Settings.Password + `</Lykilord>
		<Samningsnumer>` + cs.Settings.ContractNumber + `</Samningsnumer>
		<SamningsKennitala>` + cs.Settings.ContractIdentidyNumber + `</SamningsKennitala> 
		<PosiID>` + cs.Settings.PosID + `</PosiID>
		<Kortnumer>` + card.Number + `</Kortnumer>
		<Gildistimi>` + strconv.Itoa(card.ExpMonth) + strconv.Itoa(card.ExpYear) + `</Gildistimi>
		<Oryggisnumer>` + card.CVC + `</Oryggisnumer>
		<Stillingar></Stillingar>
		</FaSyndarkortnumer>
		</soap:Body> </soap:Envelope>`

	resp, err := helpers.Send(cs.Settings.URL, "POST", body)
	if err != nil {
		response.SystemError = err
		return
	}
	if err := xml.Unmarshal(resp, &response); err != nil {
		response.SystemError = err
		return
	}

	return
}

// VirtualCardAuthorization ...
// Documentation: https://specs.valitor.is/CorporatePayments_ISL/Web_Services/#42-faheimild
type FaHeimild struct {
	SystemError  error
	ErrorCode    int     `xml:"Body>FaHeimildResponse>FaHeimildResult>Villunumer"`
	ErrorMessage string  `xml:"Body>FaHeimildResponse>FaHeimildResult>Villuskilabod"`
	ErrorLogID   string  `xml:"Body>FaHeimildResponse>FaHeimildResult>VilluLogID"`
	Receipt      Receipt `xml:"Body>FaHeimildResponse>FaHeimildResult>Kvittun"`
}

// GetAuthorization ...
// Documentation: https://specs.valitor.is/CorporatePayments_ISL/Web_Services/#42-faheimild
func (cs *CompanyService) FaHeimild(card *Card, amount string, currency string) (response FaHeimild) {
	if err := checkCardForVirtualNumber(card); err != nil {
		response.SystemError = err
		return
	}

	if currency == "" {
		response.SystemError = errors.New("Currency missing")
		return
	}
	if amount == "" {
		response.SystemError = errors.New("Amount missing")
		return
	}

	body := `<?xml version="1.0" encoding="utf-8"?> 
		<soap:Envelope xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
		xmlns:xsd="http://www.w3.org/2001/XMLSchema" 
	 	xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
		<soap:Body>
		<FaHeimild xmlns="http://api.valitor.is/Fyrirtaekjagreidslur/">
		<Notandanafn>` + cs.Settings.Username + `</Notandanafn>
		<Lykilord>` + cs.Settings.Password + `</Lykilord>
		<Samningsnumer>` + cs.Settings.ContractNumber + `</Samningsnumer>
		<SamningsKennitala>` + cs.Settings.ContractIdentidyNumber + `</SamningsKennitala> 
		<PosiID>` + cs.Settings.PosID + `</PosiID>
		<Syndarkortnumer>` + card.VirtualNumber + `</Syndarkortnumer>
		<Upphaed>` + amount + `</Upphaed>
		<Gjaldmidill>` + strings.ToUpper(currency) + `</Gjaldmidill>
		<Stillingar></Stillingar>
		</FaHeimild>
		</soap:Body> </soap:Envelope>`

	resp, err := helpers.Send(cs.Settings.URL, "POST", body)
	if err != nil {
		response.SystemError = err
		return
	}
	if err := xml.Unmarshal(resp, &response); err != nil {
		response.SystemError = err
		return
	}

	return
}

type FaAdeinsHeimild struct {
	SystemError  error
	ErrorCode    int     `xml:"Body>FaAdeinsheimildResponse>FaAdeinsheimildResult>Villunumer"`
	ErrorMessage string  `xml:"Body>FaAdeinsheimildResponse>FaAdeinsheimildResult>Villuskilabod"`
	ErrorLogID   string  `xml:"Body>FaAdeinsheimildResponse>FaAdeinsheimildResult>VilluLogID"`
	Receipt      Receipt `xml:"Body>FaAdeinsheimildResponse>FaAdeinsheimildResult>Kvittun"`
}

func (cs *CompanyService) FaAdeinsHeimild(card *Card, amount string, currency string) (response FaAdeinsHeimild) {

	if err := checkCardForVirtualNumber(card); err != nil {
		response.SystemError = err
		return
	}

	if err := checkCardCVC(card); err != nil {
		response.SystemError = err
		return
	}
	if currency == "" {
		response.SystemError = errors.New("Currency missing")
		return
	}
	if amount == "" {
		response.SystemError = errors.New("Amount missing")
		return
	}

	body := `<?xml version="1.0" encoding="utf-8"?> 
		<soap:Envelope xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
		xmlns:xsd="http://www.w3.org/2001/XMLSchema" 
	 	xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
		<soap:Body>
		<FaAdeinsheimild xmlns="http://api.valitor.is/Fyrirtaekjagreidslur/">
		<Notandanafn>` + cs.Settings.Username + `</Notandanafn>
		<Lykilord>` + cs.Settings.Password + `</Lykilord>
		<Samningsnumer>` + cs.Settings.ContractNumber + `</Samningsnumer>
		<SamningsKennitala>` + cs.Settings.ContractIdentidyNumber + `</SamningsKennitala> 
		<PosiID>` + cs.Settings.PosID + `</PosiID>
		<Syndarkortnumer>` + card.VirtualNumber + `</Syndarkortnumer>
		<Upphaed>` + amount + `</Upphaed>
		<Gjaldmidill>` + strings.ToUpper(currency) + `</Gjaldmidill>
		<Oryggisnumer>` + card.CVC + `</Oryggisnumer>
		<Stillingar></Stillingar>
		</FaAdeinsheimild>
		</soap:Body> </soap:Envelope>`

	resp, err := helpers.Send(cs.Settings.URL, "POST", body)
	if err != nil {
		response.SystemError = err
		return
	}
	if err := xml.Unmarshal(resp, &response); err != nil {
		response.SystemError = err
		return
	}

	return
}

type NotaAdeinsheimild struct {
	SystemError  error
	ErrorCode    int    `xml:"Body>NotaAdeinsheimildResponse>NotaAdeinsheimildResult>Villunumer"`
	ErrorMessage string `xml:"Body>NotaAdeinsheimildResponse>NotaAdeinsheimildResult>Villuskilabod"`
	ErrorLogID   string `xml:"Body>NotaAdeinsheimildResponse>NotaAdeinsheimildResult>VilluLogID"`
	// Receipt      Receipt `xml:"Body>NotaAdeinsheimildResponse>NotaAdeinsheimildResult>Kvittun"`
}

func (cs *CompanyService) NotaAdeinsheimild(card *Card, authorizationNumber string) (response NotaAdeinsheimild) {

	if err := checkCardForVirtualNumber(card); err != nil {
		response.SystemError = err
		return
	}
	if err := checkCardCVC(card); err != nil {
		response.SystemError = err
		return
	}
	if authorizationNumber == "" {
		response.SystemError = errors.New("Authorization number missing")
		return
	}

	body := `<?xml version="1.0" encoding="utf-8"?> 
		<soap:Envelope xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
		xmlns:xsd="http://www.w3.org/2001/XMLSchema" 
	 	xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
		<soap:Body>
		<NotaAdeinsheimild xmlns="http://api.valitor.is/Fyrirtaekjagreidslur/">
		<Notandanafn>` + cs.Settings.Username + `</Notandanafn>
		<Lykilord>` + cs.Settings.Password + `</Lykilord>
		<Samningsnumer>` + cs.Settings.ContractNumber + `</Samningsnumer>
		<SamningsKennitala>` + cs.Settings.ContractIdentidyNumber + `</SamningsKennitala> 
		<PosiID>` + cs.Settings.PosID + `</PosiID>
		<Syndarkortnumer>` + card.VirtualNumber + `</Syndarkortnumer>
		<Oryggisnumer>` + card.CVC + `</Oryggisnumer>
		<Faerslunumer>` + authorizationNumber + `</Faerslunumer>
		<Stillingar></Stillingar>
		</NotaAdeinsheimild>
		</soap:Body> </soap:Envelope>`

	resp, err := helpers.Send(cs.Settings.URL, "POST", body)
	if err != nil {
		response.SystemError = err
		return
	}
	if err := xml.Unmarshal(resp, &response); err != nil {
		response.SystemError = err
		return
	}

	return
}

type FaEndurgreitt struct {
	SystemError  error
	ErrorCode    int     `xml:"Body>FaEndurgreittResponse>FaEndurgreittResult>Villunumer"`
	ErrorMessage string  `xml:"Body>FaEndurgreittResponse>FaEndurgreittResult>Villuskilabod"`
	ErrorLogID   string  `xml:"Body>FaEndurgreittResponse>FaEndurgreittResult>VilluLogID"`
	Receipt      Receipt `xml:"Body>FaEndurgreittResponse>FaEndurgreittResult>Kvittun"`
}

func (cs *CompanyService) FaEndurgreitt(card *Card, amount string, currency string) (response FaEndurgreitt) {

	if err := checkCardForVirtualNumber(card); err != nil {
		response.SystemError = err
		return
	}

	if currency == "" {
		response.SystemError = errors.New("Currency missing")
		return
	}

	if amount == "" {
		response.SystemError = errors.New("Amount missing")
		return
	}
	body := `<?xml version="1.0" encoding="utf-8"?> 
		<soap:Envelope xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
		xmlns:xsd="http://www.w3.org/2001/XMLSchema" 
	 	xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
		<soap:Body>
		<FaEndurgreitt xmlns="http://api.valitor.is/Fyrirtaekjagreidslur/">
		<Notandanafn>` + cs.Settings.Username + `</Notandanafn>
		<Lykilord>` + cs.Settings.Password + `</Lykilord>
		<Samningsnumer>` + cs.Settings.ContractNumber + `</Samningsnumer>
		<SamningsKennitala>` + cs.Settings.ContractIdentidyNumber + `</SamningsKennitala> 
		<PosiID>` + cs.Settings.PosID + `</PosiID>
		<Syndarkortnumer>` + card.VirtualNumber + `</Syndarkortnumer>
		<Upphaed>` + amount + `</Upphaed>
		<Gjaldmidill>` + strings.ToUpper(currency) + `</Gjaldmidill>
		<Stillingar></Stillingar>
		</FaEndurgreitt>
		</soap:Body> </soap:Envelope>`

	resp, err := helpers.Send(cs.Settings.URL, "POST", body)
	if err != nil {
		response.SystemError = err
		return
	}
	if err := xml.Unmarshal(resp, &response); err != nil {
		response.SystemError = err
		return
	}

	return
}

type FaOgildingu struct {
	SystemError  error
	ErrorCode    int     `xml:"Body>FaOgildinguResponse>FaOgildinguResult>Villunumer"`
	ErrorMessage string  `xml:"Body>FaOgildinguResponse>FaOgildinguResult>Villuskilabod"`
	ErrorLogID   string  `xml:"Body>FaOgildinguResponse>FaOgildinguResult>VilluLogID"`
	Receipt      Receipt `xml:"Body>FaOgildinguResponse>FaOgildinguResult>Kvittun"`
}

func (cs *CompanyService) FaOgildingu(card *Card, currency string, authorizationNumber string) (response FaOgildingu) {

	if err := checkCardForVirtualNumber(card); err != nil {
		response.SystemError = err
		return
	}
	if authorizationNumber == "" {
		response.SystemError = errors.New("Authorization number missing")
		return
	}
	if currency == "" {
		response.SystemError = errors.New("Currency missing")
		return
	}

	body := `<?xml version="1.0" encoding="utf-8"?> 
		<soap:Envelope xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
		xmlns:xsd="http://www.w3.org/2001/XMLSchema" 
	 	xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
		<soap:Body>
		<FaOgildingu xmlns="http://api.valitor.is/Fyrirtaekjagreidslur/">
		<Notandanafn>` + cs.Settings.Username + `</Notandanafn>
		<Lykilord>` + cs.Settings.Password + `</Lykilord>
		<Samningsnumer>` + cs.Settings.ContractNumber + `</Samningsnumer>
		<SamningsKennitala>` + cs.Settings.ContractIdentidyNumber + `</SamningsKennitala> 
		<Syndarkortnumer>` + card.VirtualNumber + `</Syndarkortnumer>
		<Faerslunumer>` + authorizationNumber + `</Faerslunumer>
		<PosiID>` + cs.Settings.PosID + `</PosiID>
		<Gjaldmidill>` + strings.ToUpper(currency) + `</Gjaldmidill>
		<Stillingar></Stillingar>
		</FaOgildingu>
		</soap:Body> </soap:Envelope>`

	resp, err := helpers.Send(cs.Settings.URL, "POST", body)
	if err != nil {
		response.SystemError = err
		return
	}
	if err := xml.Unmarshal(resp, &response); err != nil {
		response.SystemError = err
		return
	}

	return
}

type UppfaeraGildistima struct {
	SystemError  error
	ErrorCode    int    `xml:"Body>UppfaeraGildistimaResponse>UppfaeraGildistimaResult>Villunumer"`
	ErrorMessage string `xml:"Body>UppfaeraGildistimaResponse>UppfaeraGildistimaResult>Villuskilabod"`
	ErrorLogID   string `xml:"Body>UppfaeraGildistimaResponse>UppfaeraGildistimaResult>VilluLogID"`
}

func (cs *CompanyService) UppfaeraGildistima(card *Card) (response UppfaeraGildistima) {
	if err := checkCardForVirtualNumber(card); err != nil {
		response.SystemError = err
		return
	}
	if err := checkCardExpirationDate(card); err != nil {
		response.SystemError = err
		return
	}
	body := `<?xml version="1.0" encoding="utf-8"?> 
		<soap:Envelope xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
		xmlns:xsd="http://www.w3.org/2001/XMLSchema" 
	 	xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
		<soap:Body>
		<UppfaeraGildistima xmlns="http://api.valitor.is/Fyrirtaekjagreidslur/">
		<Notandanafn>` + cs.Settings.Username + `</Notandanafn>
		<Lykilord>` + cs.Settings.Password + `</Lykilord>
		<Samningsnumer>` + cs.Settings.ContractNumber + `</Samningsnumer>
		<SamningsKennitala>` + cs.Settings.ContractIdentidyNumber + `</SamningsKennitala> 
		<Syndarkortnumer>` + card.VirtualNumber + `</Syndarkortnumer>
		<NyrGildistimi>` + strconv.Itoa(card.ExpMonth) + strconv.Itoa(card.ExpYear) + `</NyrGildistimi>
		<Stillingar></Stillingar>
		</UppfaeraGildistima>
		</soap:Body> </soap:Envelope>`

	resp, err := helpers.Send(cs.Settings.URL, "POST", body)
	if err != nil {
		response.SystemError = err
		return
	}
	if err := xml.Unmarshal(resp, &response); err != nil {
		response.SystemError = err
		return
	}

	return
}

type FaSidustuFjoraIKortnumeriUtFraSyndarkortnumeri struct {
	SystemError  error
	ErrorCode    int    `xml:"Body>FaSidustuFjoraIKortnumeriUtFraSyndarkortnumeriResponse>FaSidustuFjoraIKortnumeriUtFraSyndarkortnumeriResult>Villunumer"`
	ErrorMessage string `xml:"Body>FaSidustuFjoraIKortnumeriUtFraSyndarkortnumeriResponse>FaSidustuFjoraIKortnumeriUtFraSyndarkortnumeriResult>Villuskilabod"`
	ErrorLogID   string `xml:"Body>FaSidustuFjoraIKortnumeriUtFraSyndarkortnumeriResponse>FaSidustuFjoraIKortnumeriUtFraSyndarkortnumeriResult>VilluLogID"`
	Kortnumer    string `xml:"Body>FaSidustuFjoraIKortnumeriUtFraSyndarkortnumeriResponse>FaSidustuFjoraIKortnumeriUtFraSyndarkortnumeriResult>Kortnumer"`
}

func (cs *CompanyService) FaSidustuFjoraIKortnumeriUtFraSyndarkortnumeri(card *Card) (response FaSidustuFjoraIKortnumeriUtFraSyndarkortnumeri) {
	if err := checkCardForVirtualNumber(card); err != nil {
		response.SystemError = err
		return
	}
	body := `<?xml version="1.0" encoding="utf-8"?> 
		<soap:Envelope xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
		xmlns:xsd="http://www.w3.org/2001/XMLSchema" 
	 	xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
		<soap:Body>
		<FaSidustuFjoraIKortnumeriUtFraSyndarkortnumeri xmlns="http://api.valitor.is/Fyrirtaekjagreidslur/">
		<Notandanafn>` + cs.Settings.Username + `</Notandanafn>
		<Lykilord>` + cs.Settings.Password + `</Lykilord>
		<Samningsnumer>` + cs.Settings.ContractNumber + `</Samningsnumer>
		<SamningsKennitala>` + cs.Settings.ContractIdentidyNumber + `</SamningsKennitala> 
		<Syndarkortnumer>` + card.VirtualNumber + `</Syndarkortnumer>
		<Stillingar></Stillingar>
		</FaSidustuFjoraIKortnumeriUtFraSyndarkortnumeri>
		</soap:Body> </soap:Envelope>`

	resp, err := helpers.Send(cs.Settings.URL, "POST", body)
	if err != nil {
		response.SystemError = err
		return
	}
	if err := xml.Unmarshal(resp, &response); err != nil {
		response.SystemError = err
		return
	}

	return
}

// This struct is used in many other structs, be carefull when changing it!
type Receipt struct {
	CompanyName           string `json:"VerslunNafn,omitempty" xml:"VerslunNafn,omitempty"`
	CompanyAddress        string `json:",omitempty" xml:"VerslunHeimilisfang,omitempty"`
	CompanyCity           string `json:",omitempty" xml:"VerslunStadur,omitempty"`
	CardTypeName          string `json:",omitempty" xml:"TegundKorts,omitempty"`
	CardTypeCode          string `json:",omitempty" xml:"TegundKortsKodi,omitempty"`
	Date                  string `json:",omitempty" xml:"Dagsetning,omitempty"`
	Time                  string `json:",omitempty" xml:"Timi,omitempty"`
	MaskedPAN             string `json:",omitempty" xml:"Kortnumer,omitempty"`
	Amount                int    `json:",omitempty" xml:"Upphaed,omitempty"`
	TransactionID         string `json:",omitempty" xml:"Faerslunumer,omitempty"`
	ProcessorInfo         string `json:",omitempty" xml:"Faersluhirdir,omitempty"`
	AuthorizationID       string `json:",omitempty" xml:"Heimildarnumer,omitempty"`
	PositionID            string `json:",omitempty" xml:"StadsetningNumer,omitempty"`
	WorkstationID         string `json:",omitempty" xml:"UtstodNumer,omitempty"`
	Invalidated           bool   `json:",omitempty" xml:"BuidAdOgilda,omitempty"`
	BatchNumber           string `json:",omitempty" xml:"Bunkanumer,omitempty"`
	SellerID              string `json:",omitempty" xml:"Soluadilinumer,omitempty"`
	SoftwareID            string `json:",omitempty" xml:"Hugbunadarnumer,omitempty"`
	POSID                 int    `json:",omitempty" xml:"PosiID,omitempty"`
	PINMessage            string `json:",omitempty" xml:"PinSkilabod,omitempty"`
	ReceiptMessage        string `json:",omitempty" xml:"Vidskiptaskilabod,omitempty"`
	F221to4               string `json:",omitempty" xml:"F22_1til4,omitempty"`
	LineC1                string `json:",omitempty" xml:"LinaC1,omitempty"`
	LineC2                string `json:",omitempty" xml:"LinaC2,omitempty"`
	LineC3                string `json:",omitempty" xml:"LinaC3,omitempty"`
	LineC4                string `json:",omitempty" xml:"LinaC4,omitempty"`
	LineD1                string `json:",omitempty" xml:"LinaD1,omitempty"`
	LineD2                string `json:",omitempty" xml:"LinaD2,omitempty"`
	Operation             string `json:",omitempty" xml:"TegundAdgerd,omitempty"`
	OriginalTransactionID string `json:",omitempty" xml:"FaerslunumerUpphafleguFaerslu,omitempty"`
	TerminalID            string `json:",omitempty" xml:"TerminalID,omitempty"`
}

func (r *Receipt) ToJSON() (jsonAuth []byte, err error) {
	return json.Marshal(r)
}

func checkCardCVC(card *Card) error {
	if card.CVC == "" {
		return errors.New("CVC missing")
	}
	return nil
}
func checkCardNumber(card *Card) error {
	if card.Number == "" {
		return errors.New("Card Number missing")
	}

	return nil
}
func checkCardExpirationDate(card *Card) error {
	if card.ExpYear == 0 && card.ExpMonth == 0 {
		return errors.New("Expiration Month and Year missing")
	} else if card.ExpMonth == 0 {
		return errors.New("Expiration Month missing")
	} else if card.ExpYear == 0 {
		return errors.New("Expiration Year missing")
	}
	return nil
}
func checkCardForVirtualNumber(card *Card) error {
	if card.VirtualNumber == "" {
		return errors.New("Virtual Number missing")
	}

	return nil
}
