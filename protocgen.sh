#!/usr/bin/env bash

protoc -I proto/ proto/grpc_challenge/server.proto  --go_out=plugins=grpc:types
# move proto files to the right places
#
# Note: Proto files are suffixed with the current binary version.
cd types
cp -r github.com/antstalepresh/grpc-challenge/types/* ./
rm -rf github.com

go mod tidy -compat=1.18
