package valitor

import (
	jsoncore "github.com/opensourcez/go-valitor/jsoncore"
	xmlcore "github.com/opensourcez/go-valitor/xmlcore"
)

// NewValitorPayService ...
// This payment service will use the new JSON api from Valitor
// Documentation: https://uat.valitorpay.com
func NewValitorPayService(
	agreementNumber string,
	terminalID string,
	url string,
) *jsoncore.CompanyService {

	if url == "" {
		// Setting the default url as the test url
		url = "https://uat.valitorpay.com"
	}
	return &jsoncore.CompanyService{
		Settings: &jsoncore.Settings{
			AgreementNumber: agreementNumber,
			TerminalID:      terminalID,
			URL:             url,
		},
	}
}

// NewValitorService ...
// This payment service will use the old XML endpoint from Valitor
// Documentation: https://specs.valitor.is/CorporatePayments_ISL/
func NewValitorService(
	username string,
	password string,
	contractNumber string,
	contractIdentidyNumber string,
	posID string,
	url string,
) *xmlcore.CompanyService {
	if url == "" {
		// Setting the default url as the test url
		url = "	https://api.processing.uat.valitor.com/Fyrirtaekjagreidslur/Fyrirtaekjagreidslur.asmx"
	}
	return &xmlcore.CompanyService{
		Settings: &xmlcore.Settings{
			Username:               username,
			Password:               password,
			ContractNumber:         contractNumber,
			ContractIdentidyNumber: contractIdentidyNumber,
			PosID:                  posID,
			URL:                    url,
		},
	}
}
