package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type SmartContract struct {
	contractapi.Contract
}

type BlockerCancleContract struct {
	Hash        string `json:"hash"`
	Canclehash  string `json:"canclehash"`
	Contractors string `json:"contractors"`
	Date        string `json:"date"`
}

func (s *SmartContract) Init(ctx contractapi.TransactionContextInterface) error {
	b_contract := BlockerCancleContract{
		Hash:        "Genesis",
		Canclehash:  "Genesis",
		Contractors: "Genesis",
		Date:        "Genesis",
	}
	ctrAsByte, err := json.Marshal(b_contract)
	if err != nil {
		return err
	}

	err = ctx.GetStub().PutState("Genesis", ctrAsByte)
	if err != nil {
		return fmt.Errorf("failed to put to world state. %v", err)
	}
	return nil
}

func (s *SmartContract) Getlastkey(ctx contractapi.TransactionContextInterface) (*BlockerCancleContract, error) {
	bctJSON, err := ctx.GetStub().GetState("Genesis")
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}

	if bctJSON == nil {
		return nil, fmt.Errorf("the asset genesis does not exist")
	}

	var btc BlockerCancleContract
	err = json.Unmarshal(bctJSON, &btc)
	if err != nil {
		return nil, err
	}

	return &btc, nil
}

func (s *SmartContract) Getcontract(ctx contractapi.TransactionContextInterface, keyvalue string) (*BlockerCancleContract, error) {
	bctJSON, err := ctx.GetStub().GetState(keyvalue)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}

	if bctJSON == nil {
		return nil, fmt.Errorf("the asset genesis does not exist")
	}

	var btc BlockerCancleContract
	err = json.Unmarshal(bctJSON, &btc)
	if err != nil {
		return nil, err
	}

	return &btc, nil
}

func (s *SmartContract) Break(ctx contractapi.TransactionContextInterface, keyvalue string, input_hash string, input_canclehash string, input_contractors string, input_date string) error {
	b_contract := BlockerCancleContract{
		Hash:        input_hash,
		Canclehash:  input_canclehash,
		Contractors: input_contractors,
		Date:        input_date,
	}

	ctrAsByte, err := json.Marshal(b_contract)
	if err != nil {
		return err
	}

	genesis_contract := BlockerCancleContract{
		Hash:        keyvalue,
		Canclehash:  "Genesis",
		Contractors: "Genesis",
		Date:        "Genesis",
	}

	gnsAsByte, err := json.Marshal(genesis_contract)
	if err != nil {
		return err
	}

	err = ctx.GetStub().PutState("Genesis", gnsAsByte)
	if err != nil {
		return fmt.Errorf("failed to put to world state. %v", err)
	}

	return ctx.GetStub().PutState(keyvalue, ctrAsByte)
}

func (s *SmartContract) Verification(ctx contractapi.TransactionContextInterface) ([]*BlockerCancleContract, error) {
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var bcts []*BlockerCancleContract
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var bct BlockerCancleContract
		err = json.Unmarshal(queryResponse.Value, &bct)
		if err != nil {
			return nil, err
		}
		bcts = append(bcts, &bct)
	}

	return bcts, nil
}

func main() {
	assetChaincode, err := contractapi.NewChaincode(&SmartContract{})
	if err != nil {
		log.Panicf("Error creating asset-transfer-basic chaincode: %v", err)
	}

	if err := assetChaincode.Start(); err != nil {
		log.Panicf("Error starting asset-transfer-basic chaincode: %v", err)
	}
}
