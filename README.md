<img src="gopher.png" alt="drawing" width="250"/>

# go-valitor
<a href="https://valitor.is">Valitor</a> payment service. <br>
Feel free to contribute, just poke me: (sveinn at zkynet dot io)

# Notes
1. This module is not thread safe. Make your own locks please <3
2. If you can not open issues, send me an email and I'll fix that.
3. There is a decent amount of stuff going on under the hood, so I recommend panic/defere around this module, just in case.




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
// this card is populated with example values, can't get rich of this readme, sry not sry <3
Card := &xmlcore.Card{
	Number:   "5304259906522887",
	ExpYear:  22,
	ExpMonth: 11,
	CVC:      "749",
	VirtualNumber: "",
}

response := ValitorService.GetVirtualNumber(Card)

if response.ErrorCode != 0 || response.SystemError != nil {
  Log.Println(response)
  panic("Whoopsie!")
}

Card.VirtualNumber = response.VirtualNumber

// The final step would be to save the new virtual number in your system to be used later.  
```
<br>
<br>
<br>

# The Responses
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
## Some Keywords regarding testing
- Notandanafn (User name): Valitortestfyrirtgr 
-  Lykilord (Password): testadgfyrirgr2010 
-  Samningsnumer (Contract number): 053128
-  SamningsKennitala (Contract SS number): 5006830589 
-  PosiID (POS identification): 225 


## Usable Test Cards
1. 5304259906522887 2211 749
2. 5304259909334470 2211 813
3. 5304259902386667 2211 376
