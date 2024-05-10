package main

import (
	"context"
	"fmt"
	"grpcdemo/rpcintf"
	"net"

	"github.com/google/uuid"
	"google.golang.org/grpc"
)

// MsgService is the service for message handling.
type MsgService struct {
	msger        *Msger
	handleStream func(MsgerStream)
	rpcintf.UnimplementedStreamingServer
}

// MakeServer creates a server side stream processor.
func MakeServer(
	ctx context.Context,
	msger *Msger,
	port uint16,
	streamhnd func(MsgerStream),
) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	srv := grpc.NewServer()
	rpcintf.RegisterStreamingServer(srv, &MsgService{
		msger:        msger,
		handleStream: streamhnd,
	})
	go func() {
		if err := srv.Serve(lis); err != nil {
			fmt.Printf("server terminated: %v\n", err)
		}
	}()
	go func() {
		<-ctx.Done()
		srv.GracefulStop()
	}()
	return nil
}

// TestStream is the test stream handler.
func (m *MsgService) TestStream(
	stream rpcintf.Streaming_TestStreamServer,
) error {
	streamAgent := &msgStream{
		id:     uuid.New().String(),
		origin: stream,
	}
	// new stream handler
	if m.handleStream != nil {
		m.handleStream(streamAgent)
	}
	// message read loop
	for {
		msg, err := stream.Recv()
		errPop := m.msger.popMsg(&MsgReadAgent{
			msg:   msg,
			err:   err,
			agent: streamAgent,
		})
		if err != nil {
			return nil
		}
		if errPop != nil {
			return errPop
		}
	}

	return nil
}
