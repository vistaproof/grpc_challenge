package server

import (
	"fmt"
	"reflect"

	tendermint "github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	"github.com/gogo/protobuf/jsonpb"
	"github.com/gogo/protobuf/proto"
	osmosis "github.com/osmosis-labs/osmosis/v15/x/gamm/types"
)

var serverMap = map[string]any{
	"osmosis":    (*osmosis.QueryServer)(nil),
	"tendermint": (*tendermint.ServiceServer)(nil),
}

// Get request & response type for a given method name
func getServiceMethodTypes(methodName string) (requestType reflect.Type, responseType reflect.Type, err error) {
	var method reflect.Method
	found := false

	for _, server := range serverMap {
		// Get the reflect.Type for the Registered service
		serviceType := reflect.TypeOf(server).Elem()

		// Get the reflect.Method for the service method
		method, found = serviceType.MethodByName(methodName)

		if found {
			break
		}
	}
	if !found {
		err = fmt.Errorf("service method %s not found", methodName)
		return
	}

	// Get the request message type from the method's second input parameter
	requestType = method.Type.In(1).Elem()

	// Get the response message type from the method's output parameter
	responseType = method.Type.Out(0).Elem()

	return
}

func protoMessageToJSON(protoMsg proto.Message) (string, error) {
	jsonpbMarshaler := jsonpb.Marshaler{OrigName: true}
	jsonString, err := jsonpbMarshaler.MarshalToString(protoMsg)
	if err != nil {
		return "", err
	}
	return jsonString, nil
}
