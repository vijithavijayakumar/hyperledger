package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// ServntireDemoChaincode implementation of Chaincode
type ServntireDemoChaincode struct {
}
type Smartphone struct {
	Make   string `json:"make"`
	Model  string `json:"model"`
	Colour string `json:"colour"`
	Owner  string `json:"owner"`
}

// Init of the chaincode
// This function is called only one when the chaincode is instantiated.
// So the goal is to prepare the ledger to handle future requests.
func (t *ServntireDemoChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("########### ServntireDemoChaincode Init ###########")

	// Get the function and arguments from the request
	function, _ := stub.GetFunctionAndParameters()

	// Check if the request is the init function
	if function != "init" {
		return shim.Error("Unknown function call")
	}

	// Adding to the ledger.
	smartphones := []Smartphone{
		Smartphone{Make: "Lenovo", Model: "Prius", Colour: "blue", Owner: "Tomoko"},
		Smartphone{Make: "Asus", Model: "Mustang", Colour: "red", Owner: "Brad"},
		Smartphone{Make: "Huawei", Model: "Tucson", Colour: "green", Owner: "Jin Soo"},
		Smartphone{Make: "Panasonic", Model: "Passat", Colour: "yellow", Owner: "Max"},
		Smartphone{Make: "Redmi", Model: "S", Colour: "black", Owner: "Adriana"},
		Smartphone{Make: "Apple", Model: "205", Colour: "purple", Owner: "Michel"},
		Smartphone{Make: "Nokia3110", Model: "S22L", Colour: "white", Owner: "Aarav"},
		Smartphone{Make: "iphone7", Model: "Punto", Colour: "violet", Owner: "Pari"},
		Smartphone{Make: "dell", Model: "Nano", Colour: "indigo", Owner: "Valeria"},
		Smartphone{Make: "Holden", Model: "Barina", Colour: "brown", Owner: "Shotaro"},
	}

	i := 0
	for i < len(smartphones) {
		fmt.Println("i is ", i)
		smartphoneAsBytes, _ := json.Marshal(smartphones[i])
		stub.PutState("SMARTPHONE"+strconv.Itoa(i), smartphoneAsBytes)
		fmt.Println("Added", smartphones[i])
		i = i + 1
	}

	// Return a successful message
	return shim.Success(nil)
}

// Invoke
// All future requests named invoke will arrive here.
func (t *ServntireDemoChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("########### ServntireDemoChaincode Invoke ###########")

	// Get the function and arguments from the request
	function, args := stub.GetFunctionAndParameters()

	// Check whether it is an invoke request
	if function != "invoke" {
		return shim.Error("Unknown function call")
	}

	// Check whether the number of arguments is sufficient
	if len(args) < 1 {
		return shim.Error("The number of arguments is insufficient.")
	}

	// In order to manage multiple type of request, we will check the first argument.
	// Here we have one possible argument: query (every query request will read in the ledger without modification)
	if args[0] == "query" {
		return t.query(stub, args)
	}

	// The update argument will manage all update in the ledger
	if args[0] == "invoke" {
		return t.invoke(stub, args)
	}

	// Querying Single Record by Passing SMARTPHONE ID => Key as parameter
	if args[0] == "queryone" {
		return t.queryone(stub, args)
	}

	// Adding a new transaction to the ledger
	if args[0] == "create" {
		return t.createsmartphone(stub, args)
	}

	// If the arguments given don’t match any function, we return an error
	return shim.Error("Unknown action, check the first argument")
}

// query
// Every readonly functions in the ledger will be here
func (t *ServntireDemoChaincode) query(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	// Check whether the number of arguments is sufficient
	if len(args) < 2 {
		return shim.Error("The number of arguments is insufficient.")
	}

	// Retrieves All the records here.
	if args[1] == "all" {

		//startKey := "SMARTPHONE0"
		//endKey := "SMARTPHONE999"

		// GetState by passing lower and upper limits
		resultsIterator, err := stub.GetStateByRange("", "")
		if err != nil {
			return shim.Error(err.Error())
		}
		defer resultsIterator.Close()

		// buffer is a JSON array containing QueryResults
		var buffer bytes.Buffer
		buffer.WriteString("[")

		bArrayMemberAlreadyWritten := false
		for resultsIterator.HasNext() {
			queryResponse, err := resultsIterator.Next()
			if err != nil {
				return shim.Error(err.Error())
			}
			// Add a comma before array members, suppress it for the first array member
			if bArrayMemberAlreadyWritten == true {
				buffer.WriteString(",")
			}
			buffer.WriteString("{\"Key\":")
			buffer.WriteString("\"")
			buffer.WriteString(queryResponse.Key)
			buffer.WriteString("\"")

			buffer.WriteString(", \"Record\":")
			// Record is a JSON object, so we write as-is
			buffer.WriteString(string(queryResponse.Value))
			buffer.WriteString("}")
			bArrayMemberAlreadyWritten = true
		}
		buffer.WriteString("]")

		fmt.Printf("- queryAllSmartphones:\n%s\n", buffer.String())

		return shim.Success(buffer.Bytes())
	}

	// If the arguments given don’t match any function, we return an error
	return shim.Error("Unknown query action, check the second argument.")
}

// invoke
// Every functions that read and write in the ledger will be here
func (t *ServntireDemoChaincode) invoke(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) < 2 {
		return shim.Error("The number of arguments is insufficient.")
	}

	// Changing Ownership of a Smartphone by Accepting Key and Value
	if args[1] == "changeOwner" && len(args) == 4 {

		/*
			@@@ Editing Single Field @@@
		*/
		smartphoneAsBytes, _ := stub.GetState(args[2])
		smartphone := Smartphone{}

		json.Unmarshal(smartphoneAsBytes, &smartphone)
		smartphone.Owner = args[3]

		smartphoneAsBytes, _ = json.Marshal(smartphone)
		stub.PutState(args[2], smartphoneAsBytes)
		return shim.Success(nil)
	}

	// If the arguments given don’t match any function, we return an error
	return shim.Error("Unknown invoke action, check the second argument.")
}

//  Retrieves a single record from the ledger by accepting Key value
func (t *ServntireDemoChaincode) queryone(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) < 2 {
		return shim.Error("The number of arguments is insufficient.")
	}

	// GetState retrieves the data from ledger using the Key
	smartphoneAsBytes, _ := stub.GetState(args[1])

	// Transaction Response
	return shim.Success(smartphoneAsBytes)

}

// Adds a new transaction to the ledger
func (s *ServntireDemoChaincode) createsmartphone(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) < 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	var newSmartphone Smartphone
	json.Unmarshal([]byte(args[2]), &newSmartphone)
	var smartphone = Smartphone{Make: newSmartphone.Make, Model: newSmartphone.Model, Colour: newSmartphone.Colour, Owner: newSmartphone.Owner}
	smartphoneAsBytes, _ := json.Marshal(smartphone)
	stub.PutState(args[1], smartphoneAsBytes)

	return shim.Success(nil)
}

// Get History of a transaction by passing Key
func (s *ServntireDemoChaincode) gethistory(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) < 2 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	smartphoneKey := args[1]
	fmt.Printf("##### start History of Record: %s\n", smartphoneKey)

	resultsIterator, err := stub.GetHistoryForKey(smartphoneKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing historic values for the marble
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"TxId\":")
		buffer.WriteString("\"")
		buffer.WriteString(response.TxId)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Value\":")
		// if it was a delete operation on given key, then we need to set the
		//corresponding value null. Else, we will write the response.Value
		//as-is (as the Value itself a JSON marble)
		if response.IsDelete {
			buffer.WriteString("null")
		} else {
			buffer.WriteString(string(response.Value))
		}

		buffer.WriteString(", \"Timestamp\":")
		buffer.WriteString("\"")
		buffer.WriteString(time.Unix(response.Timestamp.Seconds, int64(response.Timestamp.Nanos)).String())
		buffer.WriteString("\"")

		buffer.WriteString(", \"IsDelete\":")
		buffer.WriteString("\"")
		buffer.WriteString(strconv.FormatBool(response.IsDelete))
		buffer.WriteString("\"")

		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- getHistoryForMarble returning:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

func main() {
	// Start the chaincode and make it ready for futures requests
	err := shim.Start(new(ServntireDemoChaincode))
	if err != nil {
		fmt.Printf("Error starting Servntire Demo chaincode: %s", err)
	}
}
