package main

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/antstalepresh/grpc-challenge/types"

	tendermint "github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	"github.com/gogo/protobuf/jsonpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type BlockData struct {
	Height int64  `json:"height"`
	Hash   string `json:"hash"`
}

type FileData struct {
	TestResult []BlockData `json:"test_result"`
}

func main() {
	localRPC := fmt.Sprintf("localhost:%v", types.ServerPort)
	conn, err := grpc.Dial(localRPC, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	client := types.NewGenericServiceClient(conn)
	req := types.GenericRequest{
		Method:  "cosmos.base.tendermint.v1beta1.Service/GetLatestBlock",
		Message: "{}",
	}

	// Create a timer to periodically fetch the block data
	timer := time.NewTicker(5 * time.Second)
	defer timer.Stop()

	fileData := FileData{}

	// Loop forever and periodically fetch the latest block data
	for range timer.C {
		res, err := client.ForwardRequest(context.Background(), &req)
		if err != nil {
			panic(err)
		}

		var data tendermint.GetLatestBlockResponse
		err = jsonpb.UnmarshalString(res.Message, &data)
		// err = json.Unmarshal([]byte(res.Message), &data)
		if err != nil {
			panic(err)
		}

		fileData.TestResult = append(fileData.TestResult, BlockData{
			Height: data.Block.Header.Height,
			Hash:   hex.EncodeToString(data.BlockId.Hash),
		})

		// limit to 5 block data
		if len(fileData.TestResult) > 5 {
			fileData.TestResult = fileData.TestResult[1:]
		}

		// Convert the file data to a JSON string
		jsonData, err := json.MarshalIndent(fileData, "", "  ")
		if err != nil {
			log.Printf("Failed to convert block data to JSON: %s", err)
			continue
		}

		// Write the JSON string to a file
		err = ioutil.WriteFile("block_data.json", jsonData, 0644)
		if err != nil {
			log.Printf("Failed to write block data to file: %s", err)
			continue
		}

		log.Printf("Saved latest block data to file: %s", string(jsonData))
	}
}
