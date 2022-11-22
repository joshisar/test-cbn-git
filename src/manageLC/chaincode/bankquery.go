package chaincode

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// -------------------------------------------------------------------------------------------------------------------------------------
// GetLoCById returns the LoC stored in the channel with given {id}
func (c *LocContract) GetLoCById(ctx contractapi.TransactionContextInterface, id string) (*LoC, error) {
	// Get LoC Json if exists
	locJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		log.Println("error -> ctx.GetStub.GetState -> GetLoCById\n", err)
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if locJSON == nil {
		log.Printf("the LoC with Id@%s does not exist\n", id)
		return nil, fmt.Errorf("the LoC with Id@%s does not exist", id)
	}
	// Unmarshal and return LoC
	var loc LoC
	err = json.Unmarshal(locJSON, &loc)
	if err != nil {
		log.Println("error -> json.Unmarshal -> GetLoCById\n", err)
		return nil, fmt.Errorf("failed to unmarshal from Json: %v", err)
	}
	return &loc, nil
}

// -------------------------------------------------------------------------------------------------------------------------------------
// GetIssuedLoCs returns issued LCs for bank of invoking client
func (c *LocContract) GetIssuedLoCs(ctx contractapi.TransactionContextInterface) ([]*LoC, error) {
	bank, _ := getBankName(ctx)
	// log.Println("bank", bank)
	// Query string
	// queryString := fmt.Sprintf(`{"selector":{"docType":"asset","owner":"%s"}}`, owner)
	queryString := fmt.Sprintf(`{"selector":{"doc_type":"LoC","applicant_bank":"%s"}}`, bank)
	log.Println("queryString", queryString)
	// Get result iterator
	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		log.Println("error -> ctx.GetStub.GetQueryResult -> GetIssuedLoCs\n", err)
		return nil, fmt.Errorf("failed to get query result: %v", err)
	}
	defer resultsIterator.Close()
	var locs []*LoC
	// Iterate the results, unmarshal & append to locs & return
	for resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		if err != nil {
			log.Println("error -> resultsIterator.Next -> GetIssuedLoCs\n", err)
			return nil, fmt.Errorf("failed to read from result iterator: %v", err)
		}
		var loc LoC
		err = json.Unmarshal(queryResult.Value, &loc)
		if err != nil {
			log.Println("error -> json.Unmarshal -> GetIssuedLoCs\n", err)
			return nil, fmt.Errorf("failed to unmarshal query result: %v", err)
		}
		locs = append(locs, &loc)
	}
	return locs, nil
}

// -------------------------------------------------------------------------------------------------------------------------------------
// GetAdvisedLoCs returns advising LCs for bank of invoking client
func (c *LocContract) GetAdvisingLoCs(ctx contractapi.TransactionContextInterface) ([]*LoC, error) {
	bank, _ := getBankName(ctx)
	// Query string
	queryString := fmt.Sprintf(`{"selector":{"doc_type":"LoC","advise_through_bank":"%s"}}`, bank)
	// Get result iterator
	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		log.Println("error -> ctx.GetStub.GetQueryResult -> GetAdvisingLoCs\n", err)
		return nil, fmt.Errorf("failed to get query result: %v", err)
	}
	defer resultsIterator.Close()
	var locs []*LoC
	// Iterate the results, unmarshal & append to locs & return
	for resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		if err != nil {
			log.Println("error -> resultsIterator.Next -> GetAdvisingLoCs\n", err)
			return nil, fmt.Errorf("failed to read from result iterator: %v", err)
		}
		var loc LoC
		err = json.Unmarshal(queryResult.Value, &loc)
		if err != nil {
			log.Println("error -> json.Unmarshal -> GetAdvisingLoCs\n", err)
			return nil, fmt.Errorf("failed to unmarshal query result: %v", err)
		}
		locs = append(locs, &loc)
	}
	return locs, nil
}

// -------------------------------------------------------------------------------------------------------------------------------------
// GetNegotiatingLoCs returns negotiating LCs for bank of invoking client
func (c *LocContract) GetNegotiatingLoCs(ctx contractapi.TransactionContextInterface) ([]*LoC, error) {
	bank, _ := getBankName(ctx)
	// Query string
	queryString := fmt.Sprintf(`{"selector":{"doc_type":"LoC","negotiating_bank":"%s"}}`, bank)
	// Get result iterator
	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		log.Println("error -> ctx.GetStub.GetQueryResult -> GetNegotiatingLoCs\n", err)
		return nil, fmt.Errorf("failed to get query result: %v", err)
	}
	defer resultsIterator.Close()
	var locs []*LoC
	// Iterate the results, unmarshal & append to locs & return
	for resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		if err != nil {
			log.Println("error -> resultsIterator.Next -> GetNegotiatingLoCs\n", err)
			return nil, fmt.Errorf("failed to read from result iterator: %v", err)
		}
		var loc LoC
		err = json.Unmarshal(queryResult.Value, &loc)
		if err != nil {
			log.Println("error -> json.Unmarshal -> GetNegotiatingLoCs\n", err)
			return nil, fmt.Errorf("failed to unmarshal query result: %v", err)
		}
		locs = append(locs, &loc)
	}
	return locs, nil
}
