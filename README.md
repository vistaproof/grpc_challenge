# grpc_challenge

## Requirement
### Step 1
1. The gRPC Server should listen to requests from the user (you can send the requests using `grpcurl` package) and should use your implemented client to forward the request to the Osmosis public RPC endpoint
2. The code should be as generic as possible allowing the Osmosis endpoint to be replaced with another gRPC endpoint without many changes.
3. The server should listen to requests.
4. The project should contain a test_grpc.go file that will execute a few gRPC requests and print the results. the test should compare the server results with the direct client to the Osmosis endpoint results.
5. Implement support for at least cosmos.base.tendermint.v1beta1.Service service.

### Step 2
1. Now that you have a gRPC server running, design a state tracker that will query your server for the latest block information using the `cosmos.base.tendermint.v1beta1.Service.GetLatestBlock` API.
After getting the response you should parse the response and extract the following information:
a. block height - could be found under the key “block” → “height” 
b. block hash - could be found under the key “block_id” → “hash”

2. Parse the information for a duration of 5 blocks, for each block that passes save the data into a JSON file with the following structure:
```
{
 "test_result": [
  {"height": X, "hash" Y}, {"height": X+1, "hash" Z}, etc...
 ]
}
```

## Commands
`make protogen` will generate protobuf go files.

`make test` will run the unit tests

`make run-tracker` will run the tracker service

`make run-server` will run the gRPC service

`make build` will build two binary files(tracker, server)

`make lint` will check linter issues and fix

## How to test
1. Run `make protogen` first and it will generate necessary protobuf go fiels and download all dependencies.
2. Run the server and then the tracker using above commands(Make sure you start the gRPC server first).
3. Send query requests to local gRPC server through `grpcurl` using following format.

```
grpcurl -plaintext -d '{"method": "[method name]", "message": "[request params]"}' localhost:9090 grpc_challenge.GenericService/ForwardRequest
```

Here is some example queries with the ones using direct osmosis client.

```
grpcurl -plaintext -d '{"method": "osmosis.gamm.v1beta1.Query/TotalPoolLiquidity", "message": "{\"pool_id\": 1}"}' localhost:9090 grpc_challenge.GenericService/ForwardRequest
grpcurl -plaintext -d '{"pool_id": 1}' grpc.osmosis.zone:9090 "osmosis.gamm.v1beta1.Query/TotalPoolLiquidity"
```

```
grpcurl -plaintext -d '{"method": "osmosis.gamm.v1beta1.Query/NumPools", "message": "{}"}' localhost:9090 grpc_challenge.GenericService/ForwardRequest
grpcurl -plaintext -d '{}' grpc.osmosis.zone:9090 "osmosis.gamm.v1beta1.Query/NumPools"
```

```
grpcurl -plaintext -d '{"method": "cosmos.base.tendermint.v1beta1.Service/GetLatestBlock", "message": "{}"}' localhost:9090 grpc_challenge.GenericService/ForwardRequest
grpcurl -plaintext -d '{}' grpc.osmosis.zone:9090 "cosmos.base.tendermint.v1beta1.Service/GetLatestBlock"
```

## How to replace Osmosis with another network
1. Change `EndPoint` in types/keys.go file.
2. Open server/grpc.go and register necessary msg types by creating new query client for a chain's module.

```
	// Register osmosis & tendermint msg types
	_ = osmosis.NewQueryClient(conn)
	_ = tendermint.NewServiceClient(conn)
```
3. Open server/utils.go and replace osmosis query server with a new query server in `serverMap`.
```
var serverMap = map[string]any{
	"osmosis":    (*osmosis.QueryServer)(nil),
	"tendermint": (*tendermint.ServiceServer)(nil),
}
```

