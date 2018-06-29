package blockchain

import (
	"fmt"

	api "github.com/hyperledger/fabric-sdk-go/api"
	fcutil "github.com/hyperledger/fabric-sdk-go/pkg/util"
)

// QueryAll query the chaincode to get the state of hello
func (setup *FabricSetup) QueryAll() (string, error) {

	// Prepare arguments
	var args []string
	args = append(args, "invoke")
	args = append(args, "query")
	args = append(args, "all")

	// Make the proposal and submit it to the network (via our primary peer)
	transactionProposalResponses, _, err := fcutil.CreateAndSendTransactionProposal(
		setup.Channel,
		setup.ChaincodeId,
		setup.ChannelId,
		args,
		[]api.Peer{setup.Channel.GetPrimaryPeer()}, // Peer contacted when submitted the proposal
		nil,
	)
	if err != nil {
		return "", fmt.Errorf("Create and send transaction proposal return error in the query all records: %v", err)
	}
	return string(transactionProposalResponses[0].ProposalResponse.GetResponse().Payload), nil
}

//QueryOne - Retrieves a single record by accepting Key
func (setup *FabricSetup) QueryOne(value string) (string, string, error) {

	// Prepare arguments
	var args []string
	args = append(args, "invoke")
	args = append(args, "queryone")
	args = append(args, value)

	// Make the proposal and submit it to the network (via our primary peer)
	transactionProposalResponses, txID, err := fcutil.CreateAndSendTransactionProposal(
		setup.Channel,
		setup.ChaincodeId,
		setup.ChannelId,
		args,
		[]api.Peer{setup.Channel.GetPrimaryPeer()}, // Peer contacted when submitted the proposal
		nil,
	)
	if err != nil {
		return "", "", fmt.Errorf("Create and send transaction proposal return error in the query one record: %v", err)
	}
	return string(transactionProposalResponses[0].ProposalResponse.GetResponse().Payload), txID, nil
}

// GetHistoryofSmartphone - Retrieves history of transaction by passing Key
func (setup *FabricSetup) GetHistoryofSmartphone(value string) (string, error) {

	// Prepare arguments
	var args []string
	args = append(args, "invoke")
	args = append(args, "gethistory")
	args = append(args, value)

	// Make the proposal and submit it to the network (via our primary peer)
	transactionProposalResponses, _, err := fcutil.CreateAndSendTransactionProposal(
		setup.Channel,
		setup.ChaincodeId,
		setup.ChannelId,
		args,
		[]api.Peer{setup.Channel.GetPrimaryPeer()}, // Peer contacted when submitted the proposal
		nil,
	)
	if err != nil {
		return "", fmt.Errorf("Create and send transaction proposal return error in the get history of the record: %v", err)
	}
	return string(transactionProposalResponses[0].ProposalResponse.GetResponse().Payload), nil
}
