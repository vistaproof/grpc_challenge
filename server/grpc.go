package server

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/antstalepresh/grpc-challenge/types"
	tendermint "github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	"github.com/gogo/protobuf/jsonpb"
	"github.com/gogo/protobuf/proto"
	osmosis "github.com/osmosis-labs/osmosis/v15/x/gamm/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Define the server struct
type Server struct {
}

type QueryBalanceRequest struct {
	// address is the address to query balances for.
	Address string `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`
	// denom is the coin denom to query balances for.
	Denom string `protobuf:"bytes,2,opt,name=denom,proto3" json:"denom,omitempty"`
}

// Implement the forward request method
func (s *Server) ForwardRequest(ctx context.Context, req *types.GenericRequest) (*types.GenericResponse, error) {
	// Create a connection to the Osmosis RPC server
	osmosisRPC := types.EndPoint
	conn, err := grpc.Dial(osmosisRPC, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	// Register osmosis & tendermint msg types
	_ = osmosis.NewQueryClient(conn)
	_ = tendermint.NewServiceClient(conn)

	if err != nil {
		return nil, err
	}
	defer conn.Close()

	// Extract the gRPC method name from the request
	methodName := req.Method
	request := new(map[string]interface{})
	err = json.Unmarshal([]byte(req.Message), request)
	if err != nil {
		return nil, err
	}

	// Get request & response type
	names := strings.Split(methodName, "/")
	r1, r2, err := getServiceMethodTypes(names[1])
	if err != nil {
		return nil, err
	}

	subNames := strings.Split(names[0], ".")
	msgPrefix := strings.Join(subNames[:len(subNames)-1], ".")

	// Define request msg
	reqType := proto.MessageType(fmt.Sprintf("%v.%v", msgPrefix, strings.Split(r1.String(), ".")[1]))
	reqMessage := reflect.New(reqType.Elem()).Interface().(proto.Message)
	err = jsonpb.UnmarshalString(req.Message, reqMessage)
	if err != nil {
		return nil, err
	}

	// Define response msg
	resType := proto.MessageType(fmt.Sprintf("%v.%v", msgPrefix, strings.Split(r2.String(), ".")[1]))
	resMessage := reflect.New(resType.Elem()).Interface().(proto.Message)
	err = conn.Invoke(ctx, methodName, reqMessage, resMessage)
	if err != nil {
		return nil, err
	}

	var response types.GenericResponse
	jsonStr, err := protoMessageToJSON(resMessage)
	if err != nil {
		return nil, err
	}
	response.Message = jsonStr

	return &response, nil
}
