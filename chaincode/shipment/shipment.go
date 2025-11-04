package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing Shipments
type SmartContract struct {
	contractapi.Contract
}

// Shipment describes basic details of what makes up a shipment
type Shipment struct {
	ID          string `json:"id"`
	ProductName string `json:"productName"`
	Owner       string `json:"owner"`
	Status      string `json:"status"`
	Timestamp   string `json:"timestamp,omitempty"`
}

// CreateShipment adds a new shipment to the world state
func (s *SmartContract) CreateShipment(ctx contractapi.TransactionContextInterface, id, productName, owner, timestamp string) error {
	exists, err := s.ShipmentExists(ctx, id)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("shipment %s already exists", id)
	}
	shipment := Shipment{
		ID:          id,
		ProductName: productName,
		Owner:       owner,
		Status:      "Created",
		Timestamp:   timestamp,
	}
	data, err := json.Marshal(shipment)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(id, data)
}

// TransferShipment changes owner and sets status to In Transit
func (s *SmartContract) TransferShipment(ctx contractapi.TransactionContextInterface, id, newOwner, timestamp string) error {
	data, err := ctx.GetStub().GetState(id)
	if err != nil {
		return err
	}
	if data == nil {
		return fmt.Errorf("shipment %s does not exist", id)
	}
	var shipment Shipment
	if err := json.Unmarshal(data, &shipment); err != nil {
		return err
	}
	shipment.Owner = newOwner
	shipment.Status = "In Transit"
	shipment.Timestamp = timestamp
	newData, err := json.Marshal(shipment)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(id, newData)
}

// ReceiveShipment marks shipment as Received
func (s *SmartContract) ReceiveShipment(ctx contractapi.TransactionContextInterface, id, timestamp string) error {
	data, err := ctx.GetStub().GetState(id)
	if err != nil {
		return err
	}
	if data == nil {
		return fmt.Errorf("shipment %s does not exist", id)
	}
	var shipment Shipment
	if err := json.Unmarshal(data, &shipment); err != nil {
		return err
	}
	shipment.Status = "Received"
	shipment.Timestamp = timestamp
	newData, err := json.Marshal(shipment)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(id, newData)
}

// QueryShipment returns the shipment stored in the world state with given id.
func (s *SmartContract) QueryShipment(ctx contractapi.TransactionContextInterface, id string) (*Shipment, error) {
	data, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, fmt.Errorf("shipment %s not found", id)
	}
	var shipment Shipment
	if err := json.Unmarshal(data, &shipment); err != nil {
		return nil, err
	}
	return &shipment, nil
}

// ShipmentExists returns true when shipment with given ID exists in world state
func (s *SmartContract) ShipmentExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	data, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, err
	}
	return data != nil, nil
}

func main() {
	chaincode, err := contractapi.NewChaincode(new(SmartContract))
	if err != nil {
		fmt.Printf("Error create shipment chaincode: %s", err)
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting shipment chaincode: %s", err)
	}
}
