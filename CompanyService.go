package valitor

import (
	"encoding/json"
	"encoding/xml"
	"strings"
	"sync"
)

func NewCompanyService(
	username string,
	password string,
	contractNumber string,
	contractIdentidyNumber string,
	posID string,
	url string,
) *CompanyService {
	return &CompanyService{
		Settings: &Settings{
			Username:               username,
			Password:               password,
			ContractNumber:         contractNumber,
			ContractIdentidyNumber: contractIdentidyNumber,
			PosID:                  posID,
			URL:                    url,
		},
	}
}

type CompanyService struct {
	Settings *Settings
	Mux      sync.RWMutex
}
type CompanyServiceError struct {
	OriginalError error
	Number        string
	Message       string
	LogID         string
}

func (cs *CompanyService) GetVirtualNumber(card *Card) (response VirtualNumber) {

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
		<Gildistimi>` + card.ExpMonth + card.ExpYear + `</Gildistimi>
		<Oryggisnumer>` + card.CVC + `</Oryggisnumer>
		<Stillingar></Stillingar>
		</FaSyndarkortnumer>
		</soap:Body> </soap:Envelope>`

	resp, err := send(cs.Settings.URL, "POST", body)
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

type VirtualNumber struct {
	SystemError   error
	ErrorCode     int    `xml:"Body>FaSyndarkortnumerResponse>FaSyndarkortnumerResult>Villunumer"`
	ErrorMessage  string `xml:"Body>FaSyndarkortnumerResponse>FaSyndarkortnumerResult>Villuskilabod"`
	ErrorLogID    string `xml:"Body>FaSyndarkortnumerResponse>FaSyndarkortnumerResult>VilluLogID"`
	VirtualNumber string `xml:"Body>FaSyndarkortnumerResponse>FaSyndarkortnumerResult>Syndarkortnumer"`
}

func (cs *CompanyService) GetAuthorizationUsingAVirtualCard(card *Card, amount string, currency string) (response VirtualCardAuthorization) {

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

	resp, err := send(cs.Settings.URL, "POST", body)
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

type VirtualCardAuthorization struct {
	SystemError  error
	ErrorCode    int     `xml:"Body>FaHeimildResponse>FaHeimildResult>Villunumer"`
	ErrorMessage string  `xml:"Body>FaHeimildResponse>FaHeimildResult>Villuskilabod"`
	ErrorLogID   string  `xml:"Body>FaHeimildResponse>FaHeimildResult>VilluLogID"`
	Receipt      Receipt `xml:"Body>FaHeimildResponse>FaHeimildResult>Kvittun"`
}

func (vca *VirtualCardAuthorization) ReceiptToJSON() (jsonAuth []byte, err error) {
	return json.Marshal(vca.Receipt)
}

func (cs *CompanyService) Refund(card *Card, amount string, currency string) (response Refund) {

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

	resp, err := send(cs.Settings.URL, "POST", body)
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

type Refund struct {
	SystemError  error
	ErrorCode    int     `xml:"Body>FaEndurgreittResponse>FaEndurgreittResult>Villunumer"`
	ErrorMessage string  `xml:"Body>FaEndurgreittResponse>FaEndurgreittResult>Villuskilabod"`
	ErrorLogID   string  `xml:"Body>FaEndurgreittResponse>FaEndurgreittResult>VilluLogID"`
	Receipt      Receipt `xml:"Body>FaEndurgreittResponse>FaEndurgreittResult>Kvittun"`
}

func (vca *Refund) ReceiptToJSON() (jsonAuth []byte, err error) {
	return json.Marshal(vca.Receipt)
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
