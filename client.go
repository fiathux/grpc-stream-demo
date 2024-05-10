package main

import (
	"context"
	"fmt"
	"grpcdemo/rpcintf"

	"github.com/google/uuid"
	"google.golang.org/grpc"
)

// ClientStream represents a reference of client stream.
type ClientStream struct {
	MsgerStream
	grcpConn *grpc.ClientConn
}

// MakeClient creates a client side stream processor.
func MakeClient(msger *Msger, addr string) (*ClientStream, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	client := rpcintf.NewStreamingClient(conn)
	stream, err := client.TestStream(context.Background())
	if err != nil {
		fmt.Printf("client create stream error: %v\n", err)
		conn.Close()
		return nil, err
	}
	streamAgent := &msgStream{
		id:     uuid.New().String(),
		origin: stream,
	}
	// message read loop
	go func() {
		for {
			msg, err := stream.Recv()
			errPop := msger.popMsg(&MsgReadAgent{
				msg:   msg,
				err:   err,
				agent: streamAgent,
			})
			if err != nil {
				return
			}
			if errPop != nil {
				fmt.Printf("client pop message error: %v\n", errPop)
				return
			}
		}
	}()

	return &ClientStream{
		MsgerStream: streamAgent,
		grcpConn:    conn,
	}, nil
}

// Close closes the client stream.
func (c *ClientStream) Close() {
	c.grcpConn.Close()
}
