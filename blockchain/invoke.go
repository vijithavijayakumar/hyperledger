package blockchain

import (
	"fmt"
	"time"

	api "github.com/hyperledger/fabric-sdk-go/api"
	fcutil "github.com/hyperledger/fabric-sdk-go/pkg/util"
)

// Create Smartphone - Adding a new record to the ledger
func (setup *FabricSetup) CreateSmartphone(key, value string) (string, error) {

	// Prepare arguments
	var args []string
	args = append(args, "invoke")
	args = append(args, "create")
	args = append(args, key)
	args = append(args, value)
	//args = append(args, `{"make": "Volkswagen", "model": "Vento", "colour": "grey", "owner": "Antonio"}`)

	// Add data that will be visible in the proposal, like a description of the invoke request
	transientDataMap := make(map[string][]byte)
	transientDataMap["result"] = []byte("Transient data in Create Smartphone invoke")

	// Make a next transaction proposal and send it
	transactionProposalResponse, txID, err := fcutil.CreateAndSendTransactionProposal(
		setup.Channel,
		setup.ChaincodeId,
		setup.ChannelId,
		args,
		[]api.Peer{setup.Channel.GetPrimaryPeer()},
		transientDataMap,
	)
	if err != nil {
		return "", fmt.Errorf("Create and send transaction proposal in the invoke - Create Smartphone return error: %v", err)
	}

	// Register the Fabric SDK to listen to the event that will come back when the transaction will be send
	done, fail := fcutil.RegisterTxEvent(txID, setup.EventHub)

	// Send the final transaction signed by endorser
	if _, err := fcutil.CreateAndSendTransaction(setup.Channel, transactionProposalResponse); err != nil {
		return "", fmt.Errorf("Create and send transaction in the invoke - Create Smartphone return error: %v", err)
	}

	// Wait for the result of the submission
	select {
	// Transaction Ok
	case <-done:
		return txID, nil

	// Transaction failed
	case <-fail:
		return "", fmt.Errorf("Error received from eventhub for txid(%s) error(%v)", txID, fail)

	// Transaction timeout
	case <-time.After(time.Second * 30):
		return "", fmt.Errorf("Didn't receive block event for txid(%s)", txID)
	}
}

// Changing Ownership of a Smartphone - Updating a record
func (setup *FabricSetup) ChangeSmartphoneOwner(key, value string) (string, error) {

	// Prepare arguments
	var args []string
	args = append(args, "invoke")
	args = append(args, "invoke")
	args = append(args, "changeOwner")
	args = append(args, key)
	args = append(args, value)

	// Add data that will be visible in the proposal, like a description of the invoke request
	transientDataMap := make(map[string][]byte)
	transientDataMap["result"] = []byte("Transient data in hello invoke")

	// Make a next transaction proposal and send it
	transactionProposalResponse, txID, err := fcutil.CreateAndSendTransactionProposal(
		setup.Channel,
		setup.ChaincodeId,
		setup.ChannelId,
		args,
		[]api.Peer{setup.Channel.GetPrimaryPeer()},
		transientDataMap,
	)
	if err != nil {
		return "", fmt.Errorf("Create and send transaction proposal in the invoke - Change Smartphone Owner return error: %v", err)
	}

	// Register the Fabric SDK to listen to the event that will come back when the transaction will be send
	done, fail := fcutil.RegisterTxEvent(txID, setup.EventHub)

	// Send the final transaction signed by endorser
	if _, err := fcutil.CreateAndSendTransaction(setup.Channel, transactionProposalResponse); err != nil {
		return "", fmt.Errorf("Create and send transaction in the invoke - Change Smartphone Owner return error: %v", err)
	}

	// Wait for the result of the submission
	select {
	// Transaction Ok
	case <-done:
		return txID, nil

	// Transaction failed
	case <-fail:
		return "", fmt.Errorf("Error received from eventhub for txid(%s) error(%v)", txID, fail)

	// Transaction timeout
	case <-time.After(time.Second * 30):
		return "", fmt.Errorf("Didn't receive block event for txid(%s)", txID)
	}
}
