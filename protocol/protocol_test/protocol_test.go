package protocol_test

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/tillknuesting/stdrpc/protocol"
)

func TestSerializeAndUnserializeMessage(t *testing.T) {
	uuid := uuid.New()
	testCases := []struct {
		name   string
		msg    *protocol.Message
		expect *protocol.Message
	}{
		{
			name: "Serialize and Unserialize Message",
			msg: &protocol.Message{
				ID:       uuid,
				Function: "length",
				Parameters: []protocol.Parameter{
					{
						Type:  protocol.StringType,
						Value: "Hello, World!",
					},
				},
			},
			expect: &protocol.Message{
				ID:       uuid,
				Function: "length",
				Parameters: []protocol.Parameter{
					{
						Type:  protocol.StringType,
						Value: "Hello, World!",
					},
				},
			},
		},
		// Add more test cases as needed
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			serialized, err := protocol.SerializeMessage(tc.msg)
			if err != nil {
				t.Errorf("Error serializing message: %v", err)
			}

			unserialized, err := protocol.UnserializeMessage(serialized)
			if err != nil {
				t.Errorf("Error unserializing message: %v", err)
			}

			// Compare the unserialized message with the expected message
			if !areMessagesEqual(unserialized, tc.expect) {
				t.Errorf("Unserialized message does not match expected message")
			}
		})
	}
}

func TestCallFunction(t *testing.T) {
	functionHandlers := map[string]protocol.FunctionHandler{
		"length": lengthHandler,
		"add":    addHandler,
		// Add more function handlers as needed
	}

	uuid := uuid.New()

	testCases := []struct {
		name    string
		msg     *protocol.Message
		expect  *protocol.Message
		handler protocol.FunctionHandler
	}{
		{
			name: "Call length Function",
			msg: &protocol.Message{
				ID:       uuid,
				Function: "length",
				Parameters: []protocol.Parameter{
					{
						Type:  protocol.StringType,
						Value: "Hello, World!",
					},
				},
			},
			expect: &protocol.Message{
				ID:       uuid,
				Function: "length",
				Parameters: []protocol.Parameter{
					{
						Type:  protocol.IntType,
						Value: 13,
					},
				},
			},
			handler: lengthHandler,
		},
		{
			name: "Call add Function",
			msg: &protocol.Message{
				ID:       uuid,
				Function: "add",
				Parameters: []protocol.Parameter{
					{
						Type:  protocol.IntType,
						Value: 10,
					},
					{
						Type:  protocol.IntType,
						Value: 20,
					},
				},
			},
			expect: &protocol.Message{
				ID:       uuid,
				Function: "add",
				Parameters: []protocol.Parameter{
					{
						Type:  protocol.IntType,
						Value: 30,
					},
				},
			},
			handler: addHandler,
		},
		// Add more test cases as needed
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			response := protocol.CallFunction(tc.msg, functionHandlers)
			// Compare the response with the expected response
			if !areMessagesEqualFunctionCall(response, tc.expect) {
				t.Errorf("Response does not match expected response")
			}
		})
	}
}

func areMessagesEqual(msg1, msg2 *protocol.Message) bool {
	if msg1 == nil || msg2 == nil {
		return msg1 == nil && msg2 == nil
	}

	if msg1.ID != msg2.ID || msg1.Function != msg2.Function {
		return false
	}

	if len(msg1.Parameters) != len(msg2.Parameters) {
		return false
	}

	for i := range msg1.Parameters {
		param1 := msg1.Parameters[i]
		param2 := msg2.Parameters[i]

		if param1.Type != param2.Type || param1.Value != param2.Value {
			return false
		}
	}

	return true
}

func areMessagesEqualFunctionCall(msg1, msg2 *protocol.Message) bool {
	if msg1 == nil || msg2 == nil {
		return msg1 == nil && msg2 == nil
	}

	if msg1.ID != msg2.ID || msg1.Function != msg2.Function {
		return false
	}

	if len(msg2.Parameters) != 1 {
		return false
	}

	if msg2.Function == "add" {
		if msg2.Parameters[0].Type != protocol.IntType || msg2.Parameters[0].Value != 30 {
			return false
		}
	}

	if msg2.Function == "length" {
		if msg2.Parameters[0].Type != protocol.IntType || msg2.Parameters[0].Value != 13 {
			return false
		}
	}

	return true
}

func lengthHandler(params ...interface{}) (interface{}, error) {
	if len(params) == 0 {
		return nil, errors.New("Missing parameter")
	}

	str, ok := params[0].(string)
	if !ok {
		return nil, errors.New("Invalid parameter type")
	}

	return len(str), nil
}

func addHandler(params ...interface{}) (interface{}, error) {
	if len(params) < 2 {
		return nil, errors.New("Missing parameters")
	}

	var sum int
	for _, param := range params {
		num, ok := param.(int)
		if !ok {
			return nil, errors.New("Invalid parameter type")
		}
		sum += num
	}

	return sum, nil
}
