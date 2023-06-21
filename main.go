package main

import (
	"errors"
	"fmt"
	"github.com/tillknuesting/stdrpc/protocol"
	"log"

	"github.com/google/uuid"
)

func main() {
	// Function Handlers
	functionHandlers := map[string]protocol.FunctionHandler{
		"length": lengthHandler,
		"add":    addHandler,
		"print":  printHandler,
	}

	err := call(functionHandlers)
	lenCall(err, functionHandlers)
	addCall(err, functionHandlers)
}

func addCall(err error, functionHandlers map[string]protocol.FunctionHandler) {
	// Create a sample Message for add function
	addMsg := protocol.Message{
		ID:       uuid.New(),
		Function: "add",
		Parameters: []protocol.Parameter{
			{Type: protocol.IntType, Value: 10},
			{Type: protocol.IntType, Value: 20},
		},
	}

	// Serialize the add Message
	serializedAdd, err := protocol.SerializeMessage(&addMsg)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Serialized Add: %v\n", serializedAdd)

	// Unserialize the add Message
	unserializedAdd, err := protocol.UnserializeMessage(serializedAdd)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Unserialized Add:")
	fmt.Printf("ID: %s\n", unserializedAdd.ID)
	fmt.Printf("Function: %s\n", unserializedAdd.Function)
	fmt.Printf("Parameters: %v\n", unserializedAdd.Parameters)

	// Call the add function
	addResponse := protocol.CallFunction(unserializedAdd, functionHandlers)

	// Serialize the add response
	serializedAddResponse, err := protocol.SerializeMessage(addResponse)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Serialized Add Response: %v\n", serializedAddResponse)

	// Unserialize the add response
	unserializedAddResponse, err := protocol.UnserializeMessage(serializedAddResponse)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Unserialized Add Response:")
	fmt.Printf("ID: %s\n", unserializedAddResponse.ID)
	fmt.Printf("Function: %s\n", unserializedAddResponse.Function)
	fmt.Printf("Parameters: %v\n", unserializedAddResponse.Parameters)
}

func lenCall(err error, functionHandlers map[string]protocol.FunctionHandler) {
	// Create a sample Message for length function
	lengthMsg := protocol.Message{
		ID:       uuid.New(),
		Function: "length",
		Parameters: []protocol.Parameter{
			{Type: protocol.StringType, Value: "Hello, World!"},
		},
	}

	// Serialize the length Message
	serializedLength, err := protocol.SerializeMessage(&lengthMsg)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Serialized Length: %v\n", serializedLength)

	// Unserialize the length Message
	unserializedLength, err := protocol.UnserializeMessage(serializedLength)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Unserialized Length:")
	fmt.Printf("ID: %s\n", unserializedLength.ID)
	fmt.Printf("Function: %s\n", unserializedLength.Function)
	fmt.Printf("Parameters: %v\n", unserializedLength.Parameters)

	// Call the length function
	lengthResponse := protocol.CallFunction(unserializedLength, functionHandlers)

	// Serialize the length response
	serializedLengthResponse, err := protocol.SerializeMessage(lengthResponse)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Serialized Length Response: %v\n", serializedLengthResponse)

	// Unserialize the length response
	unserializedLengthResponse, err := protocol.UnserializeMessage(serializedLengthResponse)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Unserialized Length Response:")
	fmt.Printf("ID: %s\n", unserializedLengthResponse.ID)
	fmt.Printf("Function: %s\n", unserializedLengthResponse.Function)
	fmt.Printf("Parameters: %v\n", unserializedLengthResponse.Parameters)
}

func call(functionHandlers map[string]protocol.FunctionHandler) error {
	// Create a sample Message for length function
	printMsg := protocol.Message{
		ID:       uuid.New(),
		Function: "print",
		Parameters: []protocol.Parameter{
			{Type: protocol.StringType, Value: "Hello, World!"},
		},
	}

	// Serialize the length Message
	serializedPrint, err := protocol.SerializeMessage(&printMsg)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Serialized print: %v\n", serializedPrint)

	// Unserialize the length Message
	unserializedPrint, err := protocol.UnserializeMessage(serializedPrint)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Unserialized print:")
	fmt.Printf("ID: %s\n", unserializedPrint.ID)
	fmt.Printf("Function: %s\n", unserializedPrint.Function)
	fmt.Printf("Parameters: %v\n", unserializedPrint.Parameters)

	// Call the length function
	printResponse := protocol.CallFunction(unserializedPrint, functionHandlers)

	// Serialize the length response
	serializedPrintResponse, err := protocol.SerializeMessage(printResponse)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Serialized print Response: %v\n", serializedPrintResponse)

	// Unserialize the print response
	unserializedPrintResponse, err := protocol.UnserializeMessage(serializedPrintResponse)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Unserialized print Response:")
	fmt.Printf("ID: %s\n", unserializedPrintResponse.ID)
	fmt.Printf("Function: %s\n", unserializedPrintResponse.Function)
	fmt.Printf("Parameters: %v\n", unserializedPrintResponse.Parameters)
	return err
}

func lengthHandler(params ...interface{}) (interface{}, error) {
	if len(params) != 1 {
		return nil, errors.New("Invalid number of parameters")
	}

	str, ok := params[0].(string)
	if !ok {
		return nil, errors.New("Invalid parameter type")
	}

	result := len(str)
	return result, nil
}

func addHandler(params ...interface{}) (interface{}, error) {
	if len(params) != 2 {
		return nil, errors.New("Invalid number of parameters")
	}

	num1, ok1 := params[0].(int)
	num2, ok2 := params[1].(int)
	if !ok1 || !ok2 {
		return nil, errors.New("Invalid parameter type")
	}

	result := num1 + num2
	return result, nil
}

func printHandler(params ...interface{}) (interface{}, error) {
	if len(params) != 1 {
		return nil, errors.New("Invalid number of parameters")
	}

	stringToPrint, ok1 := params[0].(string)
	if !ok1 {
		return nil, errors.New("Invalid parameter type")
	}
	fmt.Println(stringToPrint)
	return nil, nil
}
