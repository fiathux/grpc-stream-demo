package main

//go:generate protoc --plugin=grpc --go_out=./ --go-grpc_out=./ protos/msgs.proto
//go:generate protoc --plugin=grpc --go_out=./ --go-grpc_out=./ protos/intf.proto
