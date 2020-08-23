package main

import (
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

var logger = shim.NewLogger("SimpleChaincode")

// SimpleChaincode object
type SimpleChaincode struct {
}

// Init to initialize chaincode
func (cc *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	logger.SetLevel(shim.LogDebug)
	logger.Info("SimpleChaincode.Init")
	return shim.Success(nil)
}

// Invoke is called for CRUD operations
func (cc *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	logger.Info("SimpleChaincode.Invoke")

	function, args := stub.GetFunctionAndParameters()
	logger.Debugf("function: %s", function)

	if function == "put" {
		return cc.put(stub, args)
	} else if function == "get" {
		return cc.get(stub, args)
	} else if function == "del" {
		return cc.del(stub, args)
	}

	message := fmt.Sprintf("unknown function name: %s, expected one of {get, put, del}", function)
	logger.Error(message)
	return pb.Response{Status: 400, Message: message}
}

func (cc *SimpleChaincode) put(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	logger.Info("SimpleChaincode.put")

	if len(args) != 2 {
		message := fmt.Sprintf("wrong number of arguments: passed %d, expected %d", len(args), 2)
		logger.Error(message)
		return pb.Response{Status: 400, Message: message}
	}

	key, value := args[0], args[1]
	logger.Debugf("key: %s, value: %s", key, value)

	if key == "" {
		message := "key must be a non-empty string"
		logger.Error(message)
		return pb.Response{Status: 400, Message: message}
	}

	if err := stub.PutState(key, []byte(value)); err != nil {
		message := fmt.Sprintf("unable to put a key-value pair: %s", err.Error())
		logger.Error(message)
		return shim.Error(message)
	}

	logger.Info("SimpleChaincode.put exited successfully")
	return shim.Success(nil)
}

func (cc *SimpleChaincode) get(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	logger.Info("SimpleChaincode.get")

	if len(args) != 1 {
		message := fmt.Sprintf("wrong number of arguments: passed %d, expected %d", len(args), 1)
		logger.Error(message)
		return pb.Response{Status: 400, Message: message}
	}

	key := args[0]
	logger.Debugf("key: %s", key)

	if key == "" {
		message := "key must be a non-empty string"
		logger.Error(message)
		return pb.Response{Status: 400, Message: message}
	}

	valueBytes, err := stub.GetState(key)
	if err != nil {
		message := fmt.Sprintf("unable to get a value for the key %s: %s", key, err.Error())
		logger.Error(message)
		return shim.Error(message)
	}

	if valueBytes == nil {
		message := fmt.Sprintf("a value for the key %s not found", key)
		logger.Error(message)
		return pb.Response{Status: 404, Message: message}
	}

	logger.Info("SimpleChaincode.get exited successfully")
	return shim.Success(valueBytes)
}

func (cc *SimpleChaincode) del(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	logger.Info("SimpleChaincode.del")

	if len(args) != 1 {
		message := fmt.Sprintf("wrong number of arguments: passed %d, expected %d", len(args), 1)
		logger.Error(message)
		return pb.Response{Status: 400, Message: message}
	}

	key := args[0]
	logger.Debugf("key: %s", key)

	if key == "" {
		message := "key must be a non-empty string"
		logger.Error(message)
		return pb.Response{Status: 400, Message: message}
	}

	if err := stub.DelState(key); err != nil {
		message := fmt.Sprintf("unable to delete a pair associated with the key %s: %s", key, err.Error())
		logger.Error(message)
		return shim.Error(message)
	}

	logger.Info("SimpleChaincode.del exited successfully")
	return shim.Success(nil)
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting SimpleChaincode: %s", err)
	}
}
