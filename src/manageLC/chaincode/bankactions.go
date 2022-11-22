package chaincode

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// -------------------------------------------------------------------------------------------------------------------------------------
// GetLoCIssuanceRequestsForBank returns LC Issuance Requests for given bank
func (c *LocContract) GetLoCIssuanceRequestsForBank(ctx contractapi.TransactionContextInterface) ([]*LoCIssuanceRequest, error) {
	bank, _ := getBankName(ctx)
	// queryString := fmt.Sprintf(`{"selector":{"docType":"asset","owner":"%s"}}`, owner)
	queryString := fmt.Sprintf(`{"selector":{"doc_type":"LoCIssuanceRequest","applicant_bank":"%s"}}`, bank)
	log.Println("queryString", queryString)
	// Get result iterator
	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		log.Println("error -> ctx.GetStub.GetQueryResult -> GetLoCIssuanceRequestsForBank\n", err)
		return nil, fmt.Errorf("failed to get query result: %v", err)
	}
	defer resultsIterator.Close()
	var locIssuanceRequests []*LoCIssuanceRequest
	// Iterate the results, unmarshal & append to locIssuanceRequests & return
	for resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		if err != nil {
			log.Println("error -> resultsIterator.Next -> GetLoCIssuanceRequestsForBank\n", err)
			return nil, fmt.Errorf("failed to read from result iterator: %v", err)
		}
		var locIssuanceRequest LoCIssuanceRequest
		err = json.Unmarshal(queryResult.Value, &locIssuanceRequest)
		if err != nil {
			log.Println("error -> json.Unmarshal -> GetLoCIssuanceRequestsForBank\n", err)
			return nil, fmt.Errorf("failed to unmarshal query result: %v", err)
		}
		locIssuanceRequests = append(locIssuanceRequests, &locIssuanceRequest)
	}
	return locIssuanceRequests, nil
}

// -------------------------------------------------------------------------------------------------------------------------------------
// IssueLoC issues a new LoC and puts on the ledger
func (c *LocContract) IssueLoC(ctx contractapi.TransactionContextInterface, jsonLoC string) (*RespWrapperLoC, error) {
	// only applicant bank can do it- check
	// Un-Marshal jsonLoC to loc
	var loc LoC
	json.Unmarshal([]byte(jsonLoC), &loc)

	// check for corresponding issue req Id if non-empty in args, and approve the same before issuing LoC
	bank, _ := getBankName(ctx)
	if loc.LoCIssuanceRequestID != "DIRECT_ISSUANCE_BY_BANK" && bank == loc.ApplicantBank {
		locIssuanceRequest, err := c.GetLoCIssuanceRequestById(ctx, loc.LoCIssuanceRequestID)
		if err != nil {
			log.Println("error -> c.GetLoCIssuanceRequestById -> IssueLoC\n", err)
			return nil, fmt.Errorf("LoC Issuance Request with Id@%s does not exist", loc.LoCIssuanceRequestID)
		}
		// approve its status on ledger
		locIssuanceRequest.CurrentStatus = "APPROVED"
		current_time := GetTodaysDateTimeFormatted()
		status := fmt.Sprintf("LoC Issuance Request Approved by %s on %s", loc.ApplicantBank, current_time)
		locIssuanceRequest.StatusLog = append(locIssuanceRequest.StatusLog, status)
		locIssuanceRequestJSON, err := json.Marshal(locIssuanceRequest)
		if err != nil {
			log.Println("error -> json.Marshal -> IssueLoC\n", err)
			return nil, fmt.Errorf("failed to marshal into Json: %v", err)
		}
		// Put on ledger
		err = ctx.GetStub().PutState(locIssuanceRequest.ID, locIssuanceRequestJSON)
		if err != nil {
			log.Println("error -> ctx.GetStub.PutState -> IssueLoC\n", err)
			return nil, fmt.Errorf("failed to put on ledger: %v", err)
		}
	}

	// current status
	loc.CurrentStatus = "ISSUED_BY_APPLICANT_BANK"
	// status log
	current_time := GetTodaysDateTimeFormatted()
	status := fmt.Sprintf("LoC issued by %s on %s", loc.ApplicantBank, current_time)
	loc.StatusLog = append(loc.StatusLog, status)
	// is_Active
	loc.IsActive = true
	// doc_urls - empty array of strings initialised on its own
	loc.DocsUrls = make([]string, 0)
	// timestamp unix
	ts, err := ctx.GetStub().GetTxTimestamp()
	if err != nil {
		log.Println("error -> ctx.GetStub().GetTxTimestamp() -> IssueLoC\n", err)
		return nil, fmt.Errorf("failed to get timestamp from stub: %v", err)
	}
	loc.Timestamp = ts.Seconds
	// Marshal loc
	locJSON, err := json.Marshal(loc)
	if err != nil {
		log.Println("error -> json.Marshal -> IssueLoC\n", err)
		return nil, fmt.Errorf("failed to marshal into Json: %v", err)
	}
	// Put on ledger
	err = ctx.GetStub().PutState(loc.ID, locJSON)
	if err != nil {
		log.Println("error -> ctx.GetStub.PutState -> IssueLoC\n", err)
		return nil, fmt.Errorf("failed to put on ledger: %v", err)
	}
	// Emit the LoCIssued event
	err = ctx.GetStub().SetEvent("LoCIssued", locJSON)
	if err != nil {
		log.Println("error -> ctx.GetStub.SetEvent -> IssueLoC\n", err)
		return nil, fmt.Errorf("failed to set event: %v", err)
	}
	return &RespWrapperLoC{Loc: &loc, TxId: ctx.GetStub().GetTxID()}, nil
}

// -------------------------------------------------------------------------------------------------------------------------------------
// AcknowledgeLoCIssuance acknowledges issued LoC and updates status
func (c *LocContract) AcknowledgeLoCIssuance(ctx contractapi.TransactionContextInterface, id string) (*RespWrapperLoC, error) {
	// only advising bank can do it- check
	// Get LoC if exists
	loc, err := c.GetLoCById(ctx, id)
	if err != nil {
		log.Println("error -> c.GetLoCById -> AcknowledgeLoCIssuance\n", err)
		return nil, fmt.Errorf("LoC with Id@%s does not exist", id)
	}
	// current status
	loc.CurrentStatus = "ISSUANCE_ACKNOWLEDGED_BY_ADVISING_BANK"
	// status log
	current_time := GetTodaysDateTimeFormatted()
	status := fmt.Sprintf("LoC issuance acknowledged by %s on %s", loc.AdviseThroughBank, current_time)
	loc.StatusLog = append(loc.StatusLog, status)
	// Marshal loc
	locJSON, err := json.Marshal(loc)
	if err != nil {
		log.Println("error -> json.Marshal -> AcknowledgeLoCIssuance\n", err)
		return nil, fmt.Errorf("failed to marshal into Json: %v", err)
	}
	// Put on ledger
	err = ctx.GetStub().PutState(loc.ID, locJSON)
	if err != nil {
		log.Println("error -> ctx.GetStub.PutState -> AcknowledgeLoCIssuance\n", err)
		return nil, fmt.Errorf("failed to put on ledger: %v", err)
	}
	// Emit the LoCIssuanceAcknowledged event
	err = ctx.GetStub().SetEvent("LoCIssuanceAcknowledged", locJSON)
	if err != nil {
		log.Println("error -> ctx.GetStub.SetEvent -> AcknowledgeLoCIssuance\n", err)
		return nil, fmt.Errorf("failed to set event: %v", err)
	}
	return &RespWrapperLoC{Loc: loc, TxId: ctx.GetStub().GetTxID()}, nil
}

// -------------------------------------------------------------------------------------------------------------------------------------
// AmendLoCAmount amends the LoC amount for LoC with given {id} and amount
func (c *LocContract) AmendLoCAmount(ctx contractapi.TransactionContextInterface, id string, amount int64) (*RespWrapperLoC, error) {
	// only applicant bank can do it- check
	// Get LoC if exists
	loc, err := c.GetLoCById(ctx, id)
	if err != nil {
		log.Println("error -> c.GetLoCById -> AmendLoCAmount\n", err)
		return nil, fmt.Errorf("LoC with Id@%s does not exist", id)
	}
	// Update LoC amount
	oldAmount := loc.Amount
	loc.Amount = amount
	// current status
	loc.CurrentStatus = "AMENDED_BY_APPLICANT_BANK"
	// status log
	current_time := GetTodaysDateTimeFormatted()
	status := fmt.Sprintf("LoC amount amednded by %s from %d to %d on %s", loc.ApplicantBank, oldAmount, amount, current_time)
	loc.StatusLog = append(loc.StatusLog, status)
	// Marshal loc
	locJSON, err := json.Marshal(loc)
	if err != nil {
		log.Println("error -> json.Marshal -> AmendLoCAmount\n", err)
		return nil, fmt.Errorf("failed to marshal into Json: %v", err)
	}
	// Put on ledger
	err = ctx.GetStub().PutState(loc.ID, locJSON)
	if err != nil {
		log.Println("error -> ctx.GetStub.PutState -> AmendLoCAmount\n", err)
		return nil, fmt.Errorf("failed to put on ledger: %v", err)
	}
	// Emit the LoCAmountAmended event
	err = ctx.GetStub().SetEvent("LoCAmountAmended", locJSON)
	if err != nil {
		log.Println("error -> ctx.GetStub.SetEvent -> AmendLoCAmount\n", err)
		return nil, fmt.Errorf("failed to set event: %v", err)
	}
	return &RespWrapperLoC{Loc: loc, TxId: ctx.GetStub().GetTxID()}, nil
}

// -------------------------------------------------------------------------------------------------------------------------------------
// AcknowledgeLoCAmendment acknowledges amended LoC and updates status
func (c *LocContract) AcknowledgeLoCAmendment(ctx contractapi.TransactionContextInterface, id string) (*RespWrapperLoC, error) {
	// only advising bank can do it- check
	// Get LoC if exists
	loc, err := c.GetLoCById(ctx, id)
	if err != nil {
		log.Println("error -> c.GetLoCById -> AcknowledgeLoCAmendment\n", err)
		return nil, fmt.Errorf("LoC with Id@%s does not exist", id)
	}
	// current status
	loc.CurrentStatus = "AMENDMENT_ACKNOWLEDGED_BY_ADVISING_BANK"
	// status log
	current_time := GetTodaysDateTimeFormatted()
	status := fmt.Sprintf("LoC amendment acknowledged by %s on %s", loc.AdviseThroughBank, current_time)
	loc.StatusLog = append(loc.StatusLog, status)
	// current status
	loc.CurrentStatus = "AWAITING_DOCUMENTS"
	// status log
	status = fmt.Sprintf("%s awaiting documents from %s", loc.ApplicantBank, loc.NegotiatingBank)
	loc.StatusLog = append(loc.StatusLog, status)
	// Marshal loc
	locJSON, err := json.Marshal(loc)
	if err != nil {
		log.Println("error -> json.Marshal -> AcknowledgeLoCAmendment\n", err)
		return nil, fmt.Errorf("failed to marshal into Json: %v", err)
	}
	// Put on ledger
	err = ctx.GetStub().PutState(loc.ID, locJSON)
	if err != nil {
		log.Println("error -> ctx.GetStub.PutState -> AcknowledgeLoCAmendment\n", err)
		return nil, fmt.Errorf("failed to put on ledger: %v", err)
	}
	// Emit the LoCAmendmentAcknowledged event
	err = ctx.GetStub().SetEvent("LoCAmendmentAcknowledged", locJSON)
	if err != nil {
		log.Println("error -> ctx.GetStub.SetEvent -> AcknowledgeLoCAmendment\n", err)
		return nil, fmt.Errorf("failed to set event: %v", err)
	}
	return &RespWrapperLoC{Loc: loc, TxId: ctx.GetStub().GetTxID()}, nil
}

// -------------------------------------------------------------------------------------------------------------------------------------
// SubmitDocuments just updates in ledger that docs for given LoC {id} has been submitted, TODO S3
func (c *LocContract) SubmitDocuments(ctx contractapi.TransactionContextInterface, id string, jsonDocsUrls string) (*RespWrapperLoC, error) {
	// only negotiating bank can do it- check
	docsUrls := []string{}
	json.Unmarshal([]byte(jsonDocsUrls), &docsUrls)
	log.Println("docsUrls >>>>>> \n", docsUrls)
	// Get LoC if exists
	loc, err := c.GetLoCById(ctx, id)
	if err != nil {
		log.Println("error -> c.GetLoCById -> SubmitDocuments\n", err)
		return nil, fmt.Errorf("LoC with Id@%s does not exist", id)
	}
	// append to docs Urls array
	loc.DocsUrls = docsUrls
	// current status
	loc.CurrentStatus = "DOCUMENTS_SUBMITTED_BY_NEGOTIATING_BANK"
	// status log
	current_time := GetTodaysDateTimeFormatted()
	status := fmt.Sprintf("Document(s) submitted by %s to %s on %s", loc.NegotiatingBank, loc.ApplicantBank, current_time)
	loc.StatusLog = append(loc.StatusLog, status)
	// Marshal loc
	locJSON, err := json.Marshal(loc)
	if err != nil {
		log.Println("error -> json.Marshal -> SubmitDocuments\n", err)
		return nil, fmt.Errorf("failed to marshal into Json: %v", err)
	}
	// Put on ledger
	err = ctx.GetStub().PutState(loc.ID, locJSON)
	if err != nil {
		log.Println("error -> ctx.GetStub.PutState -> SubmitDocuments\n", err)
		return nil, fmt.Errorf("failed to put on ledger: %v", err)
	}
	// Emit the DocumentsSubmitted event
	err = ctx.GetStub().SetEvent("DocumentsSubmitted", locJSON)
	if err != nil {
		log.Println("error -> ctx.GetStub.SetEvent -> SubmitDocuments\n", err)
		return nil, fmt.Errorf("failed to set event: %v", err)
	}
	return &RespWrapperLoC{Loc: loc, TxId: ctx.GetStub().GetTxID()}, nil
}

// -------------------------------------------------------------------------------------------------------------------------------------
// AcceptDocuments is done after its accepted by applicant bank for given LoC, it updates status
func (c *LocContract) AcceptDocuments(ctx contractapi.TransactionContextInterface, id string) (*RespWrapperLoC, error) {
	// only applicant bank can do it- check
	// Get LoC if exists
	loc, err := c.GetLoCById(ctx, id)
	if err != nil {
		log.Println("error -> c.GetLoCById -> AcceptDocuments\n", err)
		return nil, fmt.Errorf("LoC with Id@%s does not exist", id)
	}
	// current status
	loc.CurrentStatus = "DOCUMENTS_ACCEPTED_BY_APPLICANT_BANK"
	// status log
	current_time := GetTodaysDateTimeFormatted()
	status := fmt.Sprintf("Documents accepted by %s from %s on %s", loc.ApplicantBank, loc.NegotiatingBank, current_time)
	loc.StatusLog = append(loc.StatusLog, status)
	// Marshal loc
	locJSON, err := json.Marshal(loc)
	if err != nil {
		log.Println("error -> json.Marshal -> AcceptDocuments\n", err)
		return nil, fmt.Errorf("failed to marshal into Json: %v", err)
	}
	// Put on ledger
	err = ctx.GetStub().PutState(loc.ID, locJSON)
	if err != nil {
		log.Println("error -> ctx.GetStub.PutState -> AcceptDocuments\n", err)
		return nil, fmt.Errorf("failed to put on ledger: %v", err)
	}
	// Emit the LoCAmendmentAcknowledged event
	err = ctx.GetStub().SetEvent("DocumentsAccepted", locJSON)
	if err != nil {
		log.Println("error -> ctx.GetStub.SetEvent -> AcceptDocuments\n", err)
		return nil, fmt.Errorf("failed to set event: %v", err)
	}
	return &RespWrapperLoC{Loc: loc, TxId: ctx.GetStub().GetTxID()}, nil
}

// -------------------------------------------------------------------------------------------------------------------------------------
// ConfirmPayment just updates in ledger that payment to negotiating bank for given LoC {id} has been done, TODO S3
// func (c *LocContract) ConfirmPayment(ctx contractapi.TransactionContextInterface, id string) (*LoC, error) {
// 	// only applicant bank can do it- check
// 	// Get LoC if exists
// 	loc, err := c.GetLoCById(ctx, id)
// 	if err != nil {
// 		log.Println("error -> c.GetLoCById -> ConfirmPayment\n", err)
// 		return nil, fmt.Errorf("LoC with Id@%s does not exist", id)
// 	}
// 	// current status
// 	loc.CurrentStatus = "PAYMENT_DONE_FROM_APPLICANT_BANK_TO_NEGOTIATING_BANK"
// 	// status log
// 	current_time := GetTodaysDateTimeFormatted()
// 	status := fmt.Sprintf("Payment confirmed from %s to %s on %s", loc.ApplicantBank, loc.NegotiatingBank, current_time)
// 	loc.StatusLog = append(loc.StatusLog, status)
// 	// Marshal loc
// 	locJSON, err := json.Marshal(loc)
// 	if err != nil {
// 		log.Println("error -> json.Marshal -> ConfirmPayment\n", err)
// 		return nil, fmt.Errorf("failed to marshal into Json: %v", err)
// 	}
// 	// Put on ledger
// 	err = ctx.GetStub().PutState(loc.ID, locJSON)
// 	if err != nil {
// 		log.Println("error -> ctx.GetStub.PutState -> ConfirmPayment\n", err)
// 		return nil, fmt.Errorf("failed to put on ledger: %v", err)
// 	}
// 	// Emit the PaymentConfirmed event
// 	err = ctx.GetStub().SetEvent("PaymentConfirmed", locJSON)
// 	if err != nil {
// 		log.Println("error -> ctx.GetStub.SetEvent -> ConfirmPayment\n", err)
// 		return nil, fmt.Errorf("failed to set event: %v", err)
// 	}
// 	return loc, nil
// }

// -------------------------------------------------------------------------------------------------------------------------------------
// AcknowledgePayment is done after payment_receive is checked by negotiating bank for given LoC, it updates status
// func (c *LocContract) AcknowledgePayment(ctx contractapi.TransactionContextInterface, id string) (*LoC, error) {
// 	// only negotiating bank can do it- check
// 	// Get LoC if exists
// 	loc, err := c.GetLoCById(ctx, id)
// 	if err != nil {
// 		log.Println("error -> c.GetLoCById -> AcknowledgePayment\n", err)
// 		return nil, fmt.Errorf("LoC with Id@%s does not exist", id)
// 	}
// 	// current status
// 	loc.CurrentStatus = "PAYMENT_ACKNOWLEDGED_FROM_APPLICANT_BANK_TO_NEGOTIATING_BANK"
// 	// status log
// 	current_time := GetTodaysDateTimeFormatted()
// 	status := fmt.Sprintf("Payment acknowledged from %s to %s on %s", loc.ApplicantBank, loc.NegotiatingBank, current_time)
// 	loc.StatusLog = append(loc.StatusLog, status)
// 	// Marshal loc
// 	locJSON, err := json.Marshal(loc)
// 	if err != nil {
// 		log.Println("error -> json.Marshal -> AcknowledgePayment\n", err)
// 		return nil, fmt.Errorf("failed to marshal into Json: %v", err)
// 	}
// 	// Put on ledger
// 	err = ctx.GetStub().PutState(loc.ID, locJSON)
// 	if err != nil {
// 		log.Println("error -> ctx.GetStub.PutState -> AcknowledgePayment\n", err)
// 		return nil, fmt.Errorf("failed to put on ledger: %v", err)
// 	}
// 	// Emit the PaymentAcknowledged event
// 	err = ctx.GetStub().SetEvent("PaymentAcknowledged", locJSON)
// 	if err != nil {
// 		log.Println("error -> ctx.GetStub.SetEvent -> AcknowledgePayment\n", err)
// 		return nil, fmt.Errorf("failed to set event: %v", err)
// 	}
// 	return loc, nil
// }

// -------------------------------------------------------------------------------------------------------------------------------------
// CloseLoC closes the LoC with given {id}
func (c *LocContract) CloseLoC(ctx contractapi.TransactionContextInterface, id string) (*RespWrapperLoC, error) {
	// only applicant bank can do it- check
	// Get LoC Json if exists
	loc, err := c.GetLoCById(ctx, id)
	if err != nil {
		log.Println("error -> c.GetLoCById -> CloseLoC\n", err)
		return nil, fmt.Errorf("LoC with Id@%s does not exist", id)
	}
	// current status
	loc.CurrentStatus = "CLOSED_BY_APPLICANT_BANK"
	// status log
	current_time := GetTodaysDateTimeFormatted()
	status := fmt.Sprintf("LoC closed by %s on %s", loc.ApplicantBank, current_time)
	loc.StatusLog = append(loc.StatusLog, status)
	// mark as inactive
	loc.IsActive = false
	// Marshal loc
	locJSON, err := json.Marshal(loc)
	if err != nil {
		log.Println("error -> json.Marshal -> CloseLoC\n", err)
		return nil, fmt.Errorf("failed to marshal into Json: %v", err)
	}
	// Put on ledger
	err = ctx.GetStub().PutState(loc.ID, locJSON)
	if err != nil {
		log.Println("error -> ctx.GetStub.PutState -> CloseLoC\n", err)
		return nil, fmt.Errorf("failed to put on ledger: %v", err)
	}
	// Emit the LoCClosed event
	err = ctx.GetStub().SetEvent("LoCClosed", locJSON)
	if err != nil {
		log.Println("error -> ctx.GetStub.SetEvent -> CloseLoC\n", err)
		return nil, fmt.Errorf("failed to set event: %v", err)
	}
	return &RespWrapperLoC{Loc: loc, TxId: ctx.GetStub().GetTxID()}, nil
}
