package chaincode

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// util function for creating sha256 hash of any given generic type -> Returns Bytes
func GetSHA256Hash(data interface{}) []byte {
	bytes := []byte(fmt.Sprintf("%v", data)) // can convert struct to string via JSON, GOB too. Using fmt here
	hasher := sha256.New()
	hasher.Write(bytes)
	return hasher.Sum(nil)
}

// wrapper function to return hex of hash generated from GetSHA256Hash
func GetSHA256HashHexString(data interface{}) string {
	dataHash := GetSHA256Hash(data)
	return hex.EncodeToString(dataHash[:])
}

// wrapper function to return base64 encoded string of hash generated from GetSHA256Hash
func GetSHA256HashBase64String(data interface{}) string {
	dataHash := GetSHA256Hash(data)
	return base64.URLEncoding.EncodeToString(dataHash[:])
}

// Get current TimeStamp -> local
func GetTimeStamp() string {
	loc, _ := time.LoadLocation("Asia/Kolkata")
	current_time := time.Now().In(loc)
	return current_time.Format("20060102150405")
}

// Get Today's Date -> local
func GetTodaysDate() string {
	loc, _ := time.LoadLocation("Asia/Kolkata")
	current_time := time.Now().In(loc)
	return current_time.Format("2006-01-02")
}

// Get Today's Date & Time -> local
func GetTodaysDateTime() string {
	loc, _ := time.LoadLocation("Asia/Kolkata")
	current_time := time.Now().In(loc)
	return current_time.Format("2006-01-02 15:04:05")
}

// Get Today's Date & Time Formatted -> local
func GetTodaysDateTimeFormatted() string {
	loc, _ := time.LoadLocation("Asia/Kolkata")
	current_time := time.Now().In(loc)
	return current_time.Format("Jan 2, 2006 at 3:04 PM")
}

// getBankName is an internal helper function to get bank/org name from submitting client identity.
func getBankName(ctx contractapi.TransactionContextInterface) (string, error) {
	// Get the MSP ID of submitting client identity
	clientMSPID, err := ctx.GetClientIdentity().GetMSPID()
	log.Println("clientMSPID", clientMSPID)
	if err != nil {
		return "", fmt.Errorf("failed to get verified MSPID: %v", err)
	}
	// Get the org name from "peer-sbi-network-1-03b9"
	slice := strings.Split(clientMSPID, "-")
	org := slice[1]
	return org, nil
}

// getCorporateName is an internal helper function to get corporate name from submitting client identity cert
func getCorporateName(ctx contractapi.TransactionContextInterface) (string, error) {
	// GetAttributeValue corporateName
	val, found, err := ctx.GetClientIdentity().GetAttributeValue("corporateName")
	if err != nil {
		log.Printf("failed to get attribute@corporateName from client cert: %v\n", err)
		return "", fmt.Errorf("failed to get attribute@corporateName from client cert, err : %v", err)
	}
	if !found {
		log.Println("attribute@corporateName not found")
		return "", fmt.Errorf("attribute@corporateName not found")
	}
	log.Printf("attr{corporateName} : %v", val)
	return val, nil
}
