package chaincode

import "github.com/hyperledger/fabric-contract-api-go/contractapi"

// -------------------------------------------------------------------------------------------------------------------------------------
// InitLedger
func (c *LocContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	// creating hard-coded first LC- test
	// locs := []*LoC{{ID: "INLCU0100220001", DocType: "LoC", DocumentaryCreditNumber: "INLCU0100220001", FormOfDocumentaryCredit: "IRREVOCABLE", DateOfIssue: "20220105", DateOfExpiry: "20220221", PlaceOfExpiry: "NEGOTIATION BANK COUNTER", ApplicantBank: "sbi", Applicant: "AMBER ENTERPRISES INDIA LTD, C-3, SITE-IV, UPSIDC IND. AREA, KASNA ROAD, GREATER NOIDA-201305, U.P, INDIA", Beneficiary: "POSCO INDIA PROCESSING CENTER PVT", CurrencyCode: "INR", Amount: 11436300, AvailableWithBy: "ANY BANK IN INDIA BY NEGOTIATION", DraftsAt: "90 DAYS FROM THE DATE OF BILL OF EXCHANGE", LoadingFrom: "ANYWHERE IN INDIA", TransportationTo: "ANYWHERE IN INDIA", DescriptionOfGoodsAndServices: "100 MT OF GI SHEET AS PER PI NO. POSCO-IHPL/PI/AEPL/JAN2022/01 DTD 04.01.2022, HS CODE:72104900, CIP, ANY WHERE IN INDIA, INCOTERMS 2020", DocumentsRequired: "1: BILL OF EXCHANGE WILL BE PRESENTED AFTER DEDUCTION OF TDS AT 0.1 PCT ON BASIC VALUE OF THE INVOICE. 2: TAX INVOICE IN ONE ORIGINAL. 3: ORIGINAL LORRY RECEIPT ISSUED BY NON IBA APPROVED TRANSPORTER CONSIGNED TO RBL BANK LTD NOTIFY APPLICANT AND MARKED FREIGHT PREPAID. 4.INSURANCE POLICY/CERTIFICATE IN THE CURRENCY OF THE CREDIT AND BLANK ENDORSED FOR CIP VALUE OF GOODS PLUS 10 PCT SHOWING CLAIMS PAYABLE IN INDIA IRRESPECTIVE OF PERCENTAGE. 5: INSURANCE TO COVER ALL RISKS FROM SUPPLIER WAREHOUSE TO APPLICANT WAREHOUSE.", Charges: "APPLICANT BANK CHARGES TO APPLICANT ACCOUNT AND BENEFICIARY ACCOUNT INCLUDING DISCREPANCY CHARGES TO BENEFICIARY ACCOUNT", PeriodForPresentation: "WITHIN 21 DAYS FROM THE DATE OF SHIPMENT BUT WITHIN THE VALIDITY OF THE LC.", ReimbursingBank: "sbi", InstructionsToThePayingOrAcceptingOrNegotiatingBank: "UPON SUBMISSION OF CREDIT COMPLIANT DOCUMENTS, WE WILL REIMBURSE YOU ON DUE DATE AS PER YOUR INSTRUCTIONS", AdviseThroughBank: "icici", NegotiatingBank: "yesbank", IsActive: true, CurrentStatus: "ISSUED_BY_APPLICANT_BANK", StatusLog: []string{"LoC issued by OrgA on Apr 11, 2022 at 11:46 AM"}, DocsUrls: []string{"https://bafybeidbwaneqilaaytdvwspd6f4mvashv6wbguqxsbawbp23sbz4ypjcy.ipfs.infura-ipfs.io"}}}
	// for _, loc := range locs {
	// 	locJSON, err := json.Marshal(loc)
	// 	if err != nil {
	// 		log.Println("error -> json.Marshal -> InitLedger\n", err)
	// 	}
	// 	err = ctx.GetStub().PutState(loc.ID, locJSON)
	// 	if err != nil {
	// 		return fmt.Errorf("failed to put to world state. %v", err)
	// 	}
	// }
	return nil
}
