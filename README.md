
# go-valitor
<a href="https://valitor.is">Valitor</a> payment service. <br>
Feel free to contribute, just poke me: (sveinn at zkynet dot io)

# Notes
1. This module is not thread safe. Make your own locks please <3
2. If you can not open issues, send me an email and I'll fix that.
3. There is a decent amount of stuff going on under the hood, so I recommend panic/defer around this module, just in case.




# In progress
1. Refactor
2. Support for JSON.
3. More documentation




# XML Service
## Supported Methods 
 - https://specs.valitor.is/CorporatePayments_ISL/Web_Services/#41-fasyndarkortnumer
 - https://specs.valitor.is/CorporatePayments_ISL/Web_Services/#42-faheimild
 - https://specs.valitor.is/CorporatePayments_ISL/Web_Services/#43-faadeinsheimild
 -https://specs.valitor.is/CorporatePayments_ISL/Web_Services/#44-notaadeinsheimild
 - https://specs.valitor.is/CorporatePayments_ISL/Web_Services/#45-faendurgreitt
 - https://specs.valitor.is/CorporatePayments_ISL/Web_Services/#46-faogildingu
 - https://specs.valitor.is/CorporatePayments_ISL/Web_Services/#47-uppfaeragildistima
 - https://specs.valitor.is/CorporatePayments_ISL/Web_Services/#48-fasidustufjoraikortnumeriutfrasyndarkortnumeri

## 1. Initialize the service
  - https://specs.valitor.is/CorporatePayments_ISL/
  - https://specs.valitor.is/CorporatePayments_ISL/Test_Access/
```golang
// you can assign this as a global variable or a local one, whichever you desire.
var ValitorService = valitor.NewValitorService(
  // Username
  "Valitortestfyrirtgr",
  // Password
  "testadgfyrirgr2010",
  // Contract number
  "053128",
  // Customer Identidy Number used for this contract 
  // ( Contract SS Number )
  "5006830589",
  // PosID or TerminalID
  "225",
  // The URL, if the url is "" this module will automatically point you towards the testing url
	"https://api.processing.uat.valitor.com/Fyrirtaekjagreidslur/Fyrirtaekjagreidslur.asmx",
)

```

## 2. Get a virtual card form a real one
 - https://specs.valitor.is/CorporatePayments_ISL/Virtual_Card_Numbers/
 - https://specs.valitor.is/CorporatePayments_ISL/Web_Services/#41-fasyndarkortnumer
### !!! You should never save the actualy card information in your system/application. Virtual Cards were created to be stored on your system instead !!! 
```golang
Card := &xmlcore.Card{
	Number:   "5304259906522887",
	ExpYear:  22,
	ExpMonth: 11,
	CVC:      "749",
}

response := ValitorService.GetVirtualNumber(Card)

if response.ErrorCode != 0 || response.SystemError != nil {
 log.Println(response.ErrorCode)
 // handle error here..
}

// Assign the Virtual Card to your Card struct.. or do with it as you please. 
Card.VirtualNumber = response.VirtualNumber

// The final step would be to save the new virtual number in your system to be used later.  
```
<br>
<br>
<br>

# The Responses
## Generic Receipt Response
Example: https://specs.valitor.is/CorporatePayments_ISL/Web_Services/#42-faheimild
Web_Services/
<br/>
The receipt also has a ToJSON() function.
```golang
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

```

## Generic Error Response
Every respone will have the same error fields as shown below. 
 - General Errors: https://specs.valitor.is/CorporatePayments_ISL/Errors_Cardholder/
 - Method Specific Errors: https://specs.valitor.is/CorporatePayments_ISL/Web_Services/
```golang
type [GenericResponseStruct] struct {
  // The system error comes from this module incase something goes wrong internally. 
  SystemError   error
  
  // These three fields are Valitor Specific error fields and can have various combinations of 
  // codes and messages. 
	ErrorCode     int    `xml:"Body>[Method-Name]Response>[Method-Name]Result>Villunumer"`
	ErrorMessage  string `xml:"Body>[Method-Name]Response>[Method-Name]Result>Villuskilabod"`
	ErrorLogID    string `xml:"Body>[Method-Name]Response>[Method-Name]Result>VilluLogID"`

  // Fields specific to each request are below this line ...
}
```


# Random notes
## Testing information for the XML service
- Notandanafn (User name): Valitortestfyrirtgr 
-  Lykilord (Password): testadgfyrirgr2010 
-  Samningsnumer (Contract number): 053128
-  SamningsKennitala (Contract SS number): 5006830589 
-  PosiID (POS identification): 225 
-  Test Card 1: 5304259906522887 2211 749 
-  Test Card 2: 5304259909334470 2211 813
-  Test Card 3: 5304259902386667 2211 376

## Testing information for the JSON service
 - ... in progress