package chaincode

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// -------------------------------------------------------------------------------------------------------------------------------------
// GetLoCById returns the LoC stored in the channel with given {id}
func (c *LocContract) GetLoCIssuanceRequestById(ctx contractapi.TransactionContextInterface, id string) (*LoCIssuanceRequest, error) {
	// Get LoC Json if exists
	locIssuanceRequestJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		log.Println("error -> ctx.GetStub.GetState -> GetLoCIssuanceRequestById\n", err)
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if locIssuanceRequestJSON == nil {
		log.Printf("LoC Issuance Request with Id@%s does not exist!\n", id)
		return nil, fmt.Errorf("LoC Issuance Request with Id@%s does not exist", id)
	}
	// Unmarshal and return locIssuanceRequest
	var locIssuanceRequest LoCIssuanceRequest
	err = json.Unmarshal(locIssuanceRequestJSON, &locIssuanceRequest)
	if err != nil {
		log.Println("error -> json.Unmarshal -> GetLoCIssuanceRequestById\n", err)
		return nil, fmt.Errorf("failed to unmarshal from Json: %v", err)
	}
	return &locIssuanceRequest, nil
}

// -------------------------------------------------------------------------------------------------------------------------------------
// GetLoCIssuanceRequests returns LC Issuance Requests for given corporate identity (invoking client)
func (c *LocContract) GetLoCIssuanceRequests(ctx contractapi.TransactionContextInterface) ([]*LoCIssuanceRequest, error) {
	bank, _ := getBankName(ctx)
	corporate, _ := getCorporateName(ctx)
	// queryString := fmt.Sprintf(`{"selector":{"docType":"asset","owner":"%s"}}`, owner)
	queryString := fmt.Sprintf(`{"selector":{"doc_type":"LoCIssuanceRequest","applicant_bank":"%s","applicant":"%s"}}`, bank, corporate)
	log.Println("queryString", queryString)
	// Get result iterator
	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		log.Println("error -> ctx.GetStub.GetQueryResult -> GetLoCIssuanceRequests\n", err)
		return nil, fmt.Errorf("failed to get query result: %v", err)
	}
	defer resultsIterator.Close()
	var locIssuanceRequests []*LoCIssuanceRequest
	// Iterate the results, unmarshal & append to locIssuanceRequests & return
	for resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		if err != nil {
			log.Println("error -> resultsIterator.Next -> GetLoCIssuanceRequests\n", err)
			return nil, fmt.Errorf("failed to read from result iterator: %v", err)
		}
		var locIssuanceRequest LoCIssuanceRequest
		err = json.Unmarshal(queryResult.Value, &locIssuanceRequest)
		if err != nil {
			log.Println("error -> json.Unmarshal -> GetLoCIssuanceRequests\n", err)
			return nil, fmt.Errorf("failed to unmarshal query result: %v", err)
		}
		locIssuanceRequests = append(locIssuanceRequests, &locIssuanceRequest)
	}
	return locIssuanceRequests, nil
}

// -------------------------------------------------------------------------------------------------------------------------------------
// GetApplicantLoCs
func (c *LocContract) GetApplicantLoCs(ctx contractapi.TransactionContextInterface) ([]*LoC, error) {
	corporate, _ := getCorporateName(ctx)
	// queryString := fmt.Sprintf(`{"selector":{"docType":"asset","owner":"%s"}}`, owner)
	queryString := fmt.Sprintf(`{"selector":{"doc_type":"LoC","applicant":"%s"}}`, corporate)
	log.Println("queryString", queryString)
	// Get result iterator
	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		log.Println("error -> ctx.GetStub.GetQueryResult -> GetApplicantLoCs\n", err)
		return nil, fmt.Errorf("failed to get query result: %v", err)
	}
	defer resultsIterator.Close()
	var locs []*LoC
	// Iterate the results, unmarshal & append to locs & return
	for resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		if err != nil {
			log.Println("error -> resultsIterator.Next -> GetApplicantLoCs\n", err)
			return nil, fmt.Errorf("failed to read from result iterator: %v", err)
		}
		var loc LoC
		err = json.Unmarshal(queryResult.Value, &loc)
		if err != nil {
			log.Println("error -> json.Unmarshal -> GetApplicantLoCs\n", err)
			return nil, fmt.Errorf("failed to unmarshal query result: %v", err)
		}
		locs = append(locs, &loc)
	}
	return locs, nil
}

// -------------------------------------------------------------------------------------------------------------------------------------
// GetBeneficiaryLoCs
func (c *LocContract) GetBeneficiaryLoCs(ctx contractapi.TransactionContextInterface) ([]*LoC, error) {
	corporate, _ := getCorporateName(ctx)
	// queryString := fmt.Sprintf(`{"selector":{"docType":"asset","owner":"%s"}}`, owner)
	queryString := fmt.Sprintf(`{"selector":{"doc_type":"LoC","beneficiary":"%s"}}`, corporate)
	log.Println("queryString", queryString)
	// Get result iterator
	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		log.Println("error -> ctx.GetStub.GetQueryResult -> GetBeneficiaryLoCs\n", err)
		return nil, fmt.Errorf("failed to get query result: %v", err)
	}
	defer resultsIterator.Close()
	var locs []*LoC
	// Iterate the results, unmarshal & append to locs & return
	for resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		if err != nil {
			log.Println("error -> resultsIterator.Next -> GetBeneficiaryLoCs\n", err)
			return nil, fmt.Errorf("failed to read from result iterator: %v", err)
		}
		var loc LoC
		err = json.Unmarshal(queryResult.Value, &loc)
		if err != nil {
			log.Println("error -> json.Unmarshal -> GetBeneficiaryLoCs\n", err)
			return nil, fmt.Errorf("failed to unmarshal query result: %v", err)
		}
		locs = append(locs, &loc)
	}
	return locs, nil
}
