package protocol

import (
	"encoding/binary"
	"errors"
	"github.com/google/uuid"
)

type MessageType byte

const (
	StringType MessageType = iota + 1
	BoolType
	IntType
)

type Parameter struct {
	Type  MessageType
	Value interface{}
}

type Message struct {
	ID         uuid.UUID
	Function   string
	Parameters []Parameter
}

type FunctionHandler func(...interface{}) (interface{}, error)

func SerializeMessage(msg *Message) ([]byte, error) {
	// Serialize ID (16 bytes)
	idBytes := msg.ID[:]
	serialized := append([]byte{}, idBytes...)

	// Serialize Function
	serialized = append(serialized, byte(len(msg.Function)))
	serialized = append(serialized, []byte(msg.Function)...)

	// Serialize Parameters
	serialized = append(serialized, byte(len(msg.Parameters)))
	for _, param := range msg.Parameters {
		serialized = append(serialized, byte(param.Type))

		switch param.Type {
		case StringType:
			str := param.Value.(string)
			serialized = append(serialized, byte(len(str)))
			serialized = append(serialized, []byte(str)...)
		case BoolType:
			boolValue := param.Value.(bool)
			if boolValue {
				serialized = append(serialized, byte(1))
			} else {
				serialized = append(serialized, byte(0))
			}
		case IntType:
			intValue := param.Value.(int)
			intBytes := make([]byte, 4)
			binary.BigEndian.PutUint32(intBytes, uint32(intValue))
			serialized = append(serialized, intBytes...)
		default:
			return nil, errors.New("Invalid parameter type")
		}
	}

	return serialized, nil
}

func UnserializeMessage(data []byte) (*Message, error) {
	if len(data) < 16 {
		return nil, errors.New("Invalid message length")
	}

	msg := &Message{}
	copy(msg.ID[:], data[:16])
	data = data[16:]

	if len(data) < 1 {
		return nil, errors.New("Invalid message length")
	}

	functionLen := int(data[0])
	data = data[1:]
	if len(data) < functionLen {
		return nil, errors.New("Invalid message length")
	}

	msg.Function = string(data[:functionLen])
	data = data[functionLen:]

	if len(data) < 1 {
		return nil, errors.New("Invalid message length")
	}

	paramCount := int(data[0])
	data = data[1:]
	msg.Parameters = make([]Parameter, paramCount)

	for i := 0; i < paramCount; i++ {
		if len(data) < 1 {
			return nil, errors.New("Invalid message length")
		}

		paramType := MessageType(data[0])
		data = data[1:]

		switch paramType {
		case StringType:
			if len(data) < 1 {
				return nil, errors.New("Invalid message length")
			}

			valueLen := int(data[0])
			data = data[1:]
			if len(data) < valueLen {
				return nil, errors.New("Invalid message length")
			}

			strValue := string(data[:valueLen])
			data = data[valueLen:]

			msg.Parameters[i] = Parameter{
				Type:  paramType,
				Value: strValue,
			}
		case BoolType:
			if len(data) < 1 {
				return nil, errors.New("Invalid message length")
			}

			boolByte := data[0]
			data = data[1:]

			boolValue := false
			if boolByte == 1 {
				boolValue = true
			}

			msg.Parameters[i] = Parameter{
				Type:  paramType,
				Value: boolValue,
			}
		case IntType:
			if len(data) < 4 {
				return nil, errors.New("Invalid message length")
			}

			intValue := int(binary.BigEndian.Uint32(data[:4]))
			data = data[4:]

			msg.Parameters[i] = Parameter{
				Type:  paramType,
				Value: intValue,
			}
		default:
			return nil, errors.New("Invalid parameter type")
		}
	}

	return msg, nil
}

func CallFunction(msg *Message, functionHandlers map[string]FunctionHandler) *Message {
	handler, exists := functionHandlers[msg.Function]
	if !exists {
		return CreateErrorResponse(msg.ID, "Unknown function")
	}

	params := make([]interface{}, len(msg.Parameters))
	for i, param := range msg.Parameters {
		params[i] = param.Value
	}

	result, err := handler(params...)
	if err != nil {
		return CreateErrorResponse(msg.ID, err.Error())
	}

	responseParams := []Parameter{}
	if result != nil {
		responseParams = []Parameter{
			{Type: InferParameterType(result), Value: result},
		}
	}

	response := &Message{
		ID:         msg.ID,
		Function:   msg.Function,
		Parameters: responseParams,
	}

	return response
}

func CreateErrorResponse(id uuid.UUID, errorMsg string) *Message {
	return &Message{
		ID: id,
		Parameters: []Parameter{
			{Type: StringType, Value: errorMsg},
		},
	}
}

func InferParameterType(value interface{}) MessageType {
	switch value.(type) {
	case string:
		return StringType
	case bool:
		return BoolType
	case int:
		return IntType
	default:
		return StringType // Default to string type
	}
}
