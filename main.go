package main

import (
	"context"
	"flag"
	"fmt"
	"grpcdemo/pbmsg"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	connectAddr string
	listenPort  uint
	signalChan  = make(chan os.Signal, 1)
)

func init() {
	flag.StringVar(&connectAddr, "connect", "", "connect address")
	flag.UintVar(&listenPort, "listen", 10900, "listen port")
	flag.Parse()

	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
}

// MsgerStream is the handler for a stream error.
func handleStreamError(owner MsgerStream, err error) {
	fmt.Printf("[ERROR] stream error %s: %v\n", owner.ID(), err)
}

// handleMessageA is the handler for message A.
func handleMessageA(owner MsgerStream, msg *pbmsg.Msg_A, err error) {
	if err != nil {
		fmt.Printf("[ERROR] message A from %s with error: %v\n", owner.ID(), err)
		return
	}
	fmt.Printf("[INFO] Got message A from %s: %v\n", owner.ID(), msg)
}

// handleMessageB is the handler for message B.
func handleMessageB(owner MsgerStream, msg *pbmsg.Msg_B, err error) {
	if err != nil {
		fmt.Printf("[ERROR] message B from %s with error: %v\n", owner.ID(), err)
		return
	}
	fmt.Printf("[INFO] Got message B from %s: %v\n", owner.ID(), msg)
}

// handleMessageC is the handler for message C.
func handleMessageC(owner MsgerStream, msg *pbmsg.Msg_C, err error) {
	if err != nil {
		fmt.Printf("[ERROR] message C from %s with error: %v\n", owner.ID(), err)
		return
	}
	fmt.Printf("[INFO] Got message C from %s: %v\n", owner.ID(), msg)
}

// ------------- init process -------------

// messageSender send demo message to peer
func messageSender(stream MsgerStream, src string) {
	messages := []func() error{
		func() error {
			return stream.Send("Msg_A", &pbmsg.Msg_A{Data: "hello", Src: src})
		},
		func() error {
			return stream.Send("Msg_B", &pbmsg.Msg_B{Data: 111, Src: src})
		},
		func() error {
			return stream.Send("Msg_C", &pbmsg.Msg_C{A: "world", B: -222, Src: src})
		},
	}
	index := 0
	defer fmt.Printf("message sender stopped for %s\n", stream.ID())
	for {
		<-time.After(1 * time.Second)
		if index >= len(messages) {
			index = 0
		}
		if err := messages[index](); err != nil {
			fmt.Printf("[ERROR] send message to %s: %v\n", stream.ID(), err)
			return
		}
		index++
	}
}

// serverMode runs the program in server mode.
func serverMode(msger *Msger) {
	fmt.Println("work on server mode, listen port:", listenPort)
	if listenPort == 0 || listenPort > 65535 {
		fmt.Println("invalid listen port")
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	err := MakeServer(ctx, msger, uint16(listenPort), func(stream MsgerStream) {
		fmt.Printf("[INFO] new stream %s\n", stream.ID())
		go messageSender(stream, "[server]")
	})
	if err != nil {
		fmt.Printf("server error: %v\n", err)
		return
	}
	<-signalChan
	cancel()
}

// clientMode runs the program in client mode.
func clientMode(msger *Msger) {
	fmt.Println("work on client mode, connect to:", connectAddr)
	c, err := MakeClient(msger, connectAddr)
	if err != nil {
		fmt.Printf("client error: %v\n", err)
		return
	}
	go messageSender(c, "[client]")
	<-signalChan
	c.Close()
}

func main() {
	msger := &Msger{}
	msger.SetErrHandler(handleStreamError)
	RegHandler(msger, "Msg_A", handleMessageA)
	RegHandler(msger, "Msg_B", handleMessageB)
	RegHandler(msger, "Msg_C", handleMessageC)

	ctx, cancel := context.WithCancel(context.Background())
	msger.Serve(ctx)

	// will run and block until signal received
	if connectAddr == "" { // server mode
		serverMode(msger)
	} else { // client mode
		clientMode(msger)
	}

	cancel()
	time.Sleep(1 * time.Second)
}
