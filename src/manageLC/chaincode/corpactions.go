package chaincode

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// -------------------------------------------------------------------------------------------------------------------------------------
// RequestLoCIssuance
func (c *LocContract) RequestLoCIssuance(ctx contractapi.TransactionContextInterface, jsonIssuanceRequest string) (*RespWrapperLoCIssuanceRequest, error) {
	// only corporate can do it ABAC- check
	// Un-Marshal jsonIssuanceRequest to LoCIssuanceRequest
	var issuanceRequest LoCIssuanceRequest
	json.Unmarshal([]byte(jsonIssuanceRequest), &issuanceRequest)
	// CurrentStatus
	issuanceRequest.CurrentStatus = "PENDING_APPROVAL"
	// StatusLog
	current_time := GetTodaysDateTimeFormatted()
	status := fmt.Sprintf("LoC Issuance Requested by %s on %s", issuanceRequest.Applicant, current_time)
	issuanceRequest.StatusLog = append(issuanceRequest.StatusLog, status)
	// timestamp unix
	ts, err := ctx.GetStub().GetTxTimestamp()
	if err != nil {
		log.Println("error -> ctx.GetStub().GetTxTimestamp() -> IssueLoC\n", err)
		return nil, fmt.Errorf("failed to get timestamp from stub: %v", err)
	}
	issuanceRequest.Timestamp = ts.Seconds
	// Marshal issuanceRequest
	issuanceRequestJSON, err := json.Marshal(issuanceRequest)
	if err != nil {
		log.Println("error -> json.Marshal -> RequestLoCIssuance\n", err)
		return nil, fmt.Errorf("failed to marshal into Json: %v", err)
	}
	// Put on ledger
	err = ctx.GetStub().PutState(issuanceRequest.ID, issuanceRequestJSON)
	if err != nil {
		log.Println("error -> ctx.GetStub.PutState -> RequestLoCIssuance\n", err)
		return nil, fmt.Errorf("failed to put on ledger: %v", err)
	}
	// Emit the LoCIssued event
	err = ctx.GetStub().SetEvent("LoCIssuanceRequested", issuanceRequestJSON)
	if err != nil {
		log.Println("error -> ctx.GetStub.SetEvent -> RequestLoCIssuance\n", err)
		return nil, fmt.Errorf("failed to set event: %v", err)
	}
	return &RespWrapperLoCIssuanceRequest{LocIssuanceRequest: &issuanceRequest, TxId: ctx.GetStub().GetTxID()}, nil
}
