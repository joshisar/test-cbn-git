package chaincode

import (
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// C H A I N C O D E
// TODO- check in corporate action, correct access via ABAC/attributes
// TODO- ensure the correct in order movement of states- chaincode, AND accordingly, only correct actions can-be-taken
// TODO- check in issuance , with same id, present or not, before issuing
// TODO- only advising-org can call its flow
// TODO- only issuing-org can call its flow
// TODO- only negotiating-org can call its flow
// TODO- only parties involved can access read locs from ledger(restrict)

// SmartContract provides functions for creating & managing our LoC
type LocContract struct {
	contractapi.Contract
}

// LoC describes basic details of what makes up a letter of credit
type LoC struct {
	ID                                                  string   `json:"ID"` // serial number which uniquely identifies the LoC
	LoCIssuanceRequestID                                string   `json:"loc_issuance_request_id"`
	DocType                                             string   `json:"doc_type"` // doc_type is used to distinguish the various types of objects in state database
	DocumentaryCreditNumber                             string   `json:"documentary_credit_number"`
	FormOfDocumentaryCredit                             string   `json:"form_of_documentary_credit"`
	DateOfIssue                                         string   `json:"date_of_issue"`
	DateOfExpiry                                        string   `json:"date_of_expiry"`
	PlaceOfExpiry                                       string   `json:"place_of_expiry"`
	ApplicantBank                                       string   `json:"applicant_bank"`
	Applicant                                           string   `json:"applicant"`
	Beneficiary                                         string   `json:"beneficiary"`
	CurrencyCode                                        string   `json:"currency_code"`
	Amount                                              int64    `json:"amount"`
	AvailableWithBy                                     string   `json:"available_with_by"`
	DraftsAt                                            string   `json:"drafts_at"`
	LoadingFrom                                         string   `json:"loading_from"`
	TransportationTo                                    string   `json:"transportation_to"`
	DescriptionOfGoodsAndServices                       string   `json:"description_of_goods_and_services"`
	DocumentsRequired                                   string   `json:"documents_required"`
	Charges                                             string   `json:"charges"`
	PeriodForPresentation                               string   `json:"period_for_presentation"`
	ReimbursingBank                                     string   `json:"reimbursing_bank"`
	InstructionsToThePayingOrAcceptingOrNegotiatingBank string   `json:"instructions_to_the_paying_or_accepting_or_negotiating_bank"`
	AdviseThroughBank                                   string   `json:"advise_through_bank"`
	NegotiatingBank                                     string   `json:"negotiating_bank"`
	IsActive                                            bool     `json:"is_active"` // Is LoC active/expired
	CurrentStatus                                       string   `json:"current_status"`
	StatusLog                                           []string `json:"status_log"`
	DocsUrls                                            []string `json:"docs_urls"`
	Timestamp                                           int64    `json:"timestamp"`
}

// RespWrapperLoC
type RespWrapperLoC struct {
	Loc  *LoC   `json:"loc"`
	TxId string `json:"txid"`
}

// LoCIssuanceRequest
type LoCIssuanceRequest struct {
	ID                                                  string   `json:"ID"`       // serial number which uniquely identifies the LoC Issuance Request
	DocType                                             string   `json:"doc_type"` // doc_type is used to distinguish the various types of objects in state database
	FormOfDocumentaryCredit                             string   `json:"form_of_documentary_credit"`
	PlaceOfExpiry                                       string   `json:"place_of_expiry"`
	ApplicantBank                                       string   `json:"applicant_bank"`
	Applicant                                           string   `json:"applicant"`
	Beneficiary                                         string   `json:"beneficiary"`
	CurrencyCode                                        string   `json:"currency_code"`
	Amount                                              int64    `json:"amount"`
	AvailableWithBy                                     string   `json:"available_with_by"`
	DraftsAt                                            string   `json:"drafts_at"`
	LoadingFrom                                         string   `json:"loading_from"`
	TransportationTo                                    string   `json:"transportation_to"`
	DescriptionOfGoodsAndServices                       string   `json:"description_of_goods_and_services"`
	DocumentsRequired                                   string   `json:"documents_required"`
	PeriodForPresentation                               string   `json:"period_for_presentation"`
	ReimbursingBank                                     string   `json:"reimbursing_bank"`
	InstructionsToThePayingOrAcceptingOrNegotiatingBank string   `json:"instructions_to_the_paying_or_accepting_or_negotiating_bank"`
	AdviseThroughBank                                   string   `json:"advise_through_bank"`
	NegotiatingBank                                     string   `json:"negotiating_bank"`
	CurrentStatus                                       string   `json:"current_status"`
	StatusLog                                           []string `json:"status_log"`
	Timestamp                                           int64    `json:"timestamp"`
}

// RespWrapperLoCIssuanceRequest
type RespWrapperLoCIssuanceRequest struct {
	LocIssuanceRequest *LoCIssuanceRequest `json:"loc_issuance_request"`
	TxId               string              `json:"txid"`
}

// ********************** HAPPY FLOW START **********************
// ----------------------------------------------------------------
// xxx ISSUANCE_REQUESTED_BY_APPLICANT
// ----------------------------------------------------------------
// vvv ISSUED_BY_APPLICANT_BANK
// vvv ISSUANCE_ACKNOWLEDGED_BY_ADVISING_BANK
// ----------------------------------------------------------------
// vvv AMENDED_BY_APPLICANT_BANK
// vvv AMENDMENT_ACKNOWLEDGED_BY_ADVISING_BANK
// ----------------------------------------------------------------
// xxx AWAITING_DOCUMENTS
// ----------------------------------------------------------------
// xxx DOCUMENTS_SUBMITTED_BY_BENEFICIARY
// xxx PAYMENT_MADE_TO_BENEFICIARY
// ----------------------------------------------------------------
// vvv DOCUMENTS_SUBMITTED_BY_NEGOTIATING_BANK
// vvv DOCUMENTS_ACCEPTED_BY_APPLICANT_BANK
// ----------------------------------------------------------------
// xxx PAYMENT_MADE_TO_NEGOTIATING_BANK
// ----------------------------------------------------------------
// xxx PAYMENT_ACKNOWLEDGED_BY_NEGOTIATING_BANK
// ----------------------------------------------------------------
// vvv CLOSED_BY_APPLICANT_BANK
// ----------------------------------------------------------------
// ********************** HAPPY FLOW END **********************
