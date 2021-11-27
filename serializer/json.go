package serializer

import (
	"github.com/golang/protobuf/jsonpb"
	"google.golang.org/protobuf/proto"
)

// ProtobufToJSON converts protocol buffer message to JSON string
func ProtobufToJSON(message proto.Message) (string, error) {
	marshaler := jsonpb.Marshaler{
		EnumAsInt:    false,
		EmitDefaults: true,
		Indent:       " ",
		OrigName:     true,
	}

	return marshaler.MarshalToString(message)
}

// JSONToProtobufMessage converts JSON string to protocol buffer message
func JSONToProtobufMessage(data string, message proto.Message) error {
	return jsonpb.UnmarshalString(data, message)
}
