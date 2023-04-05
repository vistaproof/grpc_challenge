package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"

	"testing"

	"github.com/antstalepresh/grpc-challenge/types"
	osmosis "github.com/osmosis-labs/osmosis/v15/x/gamm/types"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

type GrpcTestSuite struct {
	suite.Suite
}

func (s *GrpcTestSuite) SetupTest() {

}

func (s *GrpcTestSuite) convertReqToGenericReq(methodName string, methodReq interface{}) (*types.GenericRequest, error) {
	bz, err := json.Marshal(methodReq)
	if err != nil {
		return nil, err
	}

	req := types.GenericRequest{
		Method:  methodName,
		Message: string(bz),
	}
	return &req, err
}

func (s *GrpcTestSuite) TestGrpcServer() {
	// Create direct client for osmosis
	osmosisAddr := types.EndPoint
	conn, err := grpc.Dial(osmosisAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	s.Require().NoError(err)

	osmosisClient := osmosis.NewQueryClient(conn)

	// Create proxy server for osmosis
	srv := Server{}

	// Compare two result
	//--------------- osmosis.gamm.v1beta1.Query/PoolType ---------------------
	methodName := "osmosis.gamm.v1beta1.Query/PoolType"
	methodReq1 := osmosis.QueryPoolTypeRequest{
		PoolId: 1,
	}

	req, err := s.convertReqToGenericReq(methodName, methodReq1)
	s.Require().NoError(err)

	res1, err := srv.ForwardRequest(context.Background(), req)
	s.Require().NoError(err)

	res2, err := osmosisClient.PoolType(context.Background(), &methodReq1)
	s.Require().NoError(err)
	res2Str, err := protoMessageToJSON(res2)
	s.Require().NoError(err)

	s.Require().Equal(res1.Message, res2Str)

	//--------------- osmosis.gamm.v1beta1.Query/NumPools ---------------------
	methodName = "osmosis.gamm.v1beta1.Query/NumPools"
	methodReq2 := osmosis.QueryNumPoolsRequest{}

	req, err = s.convertReqToGenericReq(methodName, methodReq2)
	s.Require().NoError(err)

	res1, err = srv.ForwardRequest(context.Background(), req)
	s.Require().NoError(err)

	res3, err := osmosisClient.NumPools(context.Background(), &methodReq2)
	s.Require().NoError(err)
	res3Str, err := protoMessageToJSON(res3)
	s.Require().NoError(err)

	s.Require().Equal(res1.Message, res3Str)
}

func TestGrpcTestSuite(t *testing.T) {
	testSuite := *new(GrpcTestSuite)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", types.ServerPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Set up the gRPC server
	srv := grpc.NewServer()
	types.RegisterGenericServiceServer(srv, &Server{})
	reflection.Register(srv)

	// Start the server
	go srv.Serve(lis)
	suite.Run(t, &testSuite)
	srv.Stop()
}
