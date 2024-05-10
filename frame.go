package main

import (
	"context"
	"fmt"
	"grpcdemo/rpcintf"
	"sync"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

// MsgerStream is the interface for a stream of messages.
type MsgerStream interface {
	ID() string
	Send(typename string, m proto.Message) error
}

// MsgerStreamOrigin is the reference of original grpc stream.
type MsgerStreamOrigin interface {
	Send(*rpcintf.MsgAgent) error
}

// MsgStream is the stream of messages.
type msgStream struct {
	id     string
	origin MsgerStreamOrigin
}

// MsgAgent is the agent for message handling.
type MsgReadAgent struct {
	msg   *rpcintf.MsgAgent
	err   error
	agent MsgerStream
}

// Msger is the framework for message handling.
type Msger struct {
	ctx       context.Context
	cancel    context.CancelFunc
	mtx       sync.RWMutex
	msgParser map[string]func(MsgerStream, *anypb.Any)
	chRead    chan *MsgReadAgent
	errhnd    func(MsgerStream, error)
}

// RegHandler registers a handler for a message type.
func RegHandler[M any](
	msger *Msger, name string,
	hnd func(MsgerStream, *M, error),
) error {
	// test whether handle is available
	if !func() bool {
		var m M
		_, ok := any(&m).(proto.Message)
		return ok
	}() {
		return fmt.Errorf("A handle must be accepted a protobuf message")
	}

	msger.mtx.Lock()
	defer msger.mtx.Unlock()
	if msger.ctx != nil {
		return fmt.Errorf("Msger is already started")
	}
	// make proto parser
	if msger.msgParser == nil {
		msger.msgParser = make(map[string]func(MsgerStream, *anypb.Any))
	}
	if _, ok := msger.msgParser[name]; ok {
		return fmt.Errorf("handler for %s already exists", name)
	}
	msger.msgParser[name] = func(owner MsgerStream, a *anypb.Any) {
		var m M
		um := proto.Message((any(&m)).(proto.Message))
		if err := a.UnmarshalTo(um); err != nil {
			hnd(owner, nil, err)
		}
		hnd(owner, &m, nil)
	}
	return nil
}

// ------------------ Msger ------------------

// Serve starts the Msger.
func (m *Msger) Serve(ctx context.Context) error {
	if m == nil {
		return fmt.Errorf("Msger is nil")
	}
	// check and init
	if err := func() error {
		m.mtx.Lock()
		defer m.mtx.Unlock()
		if m.ctx != nil {
			return fmt.Errorf("Msger is already started")
		}
		m.ctx, m.cancel = context.WithCancel(ctx)
		m.chRead = make(chan *MsgReadAgent, 10)
		return nil
	}(); err != nil {
		return err
	}

	go func() {
		defer func() {
			m.mtx.Lock()
			defer m.mtx.Unlock()
			m.ctx = nil
			m.cancel = nil
			close(m.chRead)
		}()
		for {
			select {
			case <-m.ctx.Done():
				return
			case msg := <-m.chRead:
				m.dispatch(msg)
			}
		}
	}()
	return nil
}

// Started returns whether the Msger is started.
func (m *Msger) Started() bool {
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	return m.ctx != nil
}

// Close closes the Msger.
func (m *Msger) Close() {
	if m == nil {
		return
	}
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	if m.ctx == nil {
		return
	}
	m.cancel()
}

func (m *Msger) SetErrHandler(hnd func(MsgerStream, error)) error {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	if m.ctx != nil {
		return fmt.Errorf("Msger is already started")
	}
	m.errhnd = hnd
	return nil
}

// popMsg pops a message to the message handler.
func (m *Msger) popMsg(msg *MsgReadAgent) error {
	if m == nil {
		return fmt.Errorf("Msger is nil")
	}
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	if m.ctx == nil {
		return fmt.Errorf("Msger is not started")
	}
	m.chRead <- msg
	return nil
}

// dispatch dispatches a message to the handler.
func (m *Msger) dispatch(msg *MsgReadAgent) {
	if msg.err != nil {
		if m.errhnd != nil {
			m.errhnd(msg.agent, msg.err)
		}
		return
	}
	if hnd := m.msgParser[msg.msg.Type]; hnd != nil {
		hnd(msg.agent, msg.msg.Body)
	}
}

// ------------------ msgStream ------------------

// ID returns the ID of the stream.
func (m *msgStream) ID() string {
	return m.id
}

// Send sends a message.
func (m *msgStream) Send(typename string, msg proto.Message) error {
	apb, err := anypb.New(msg)
	if err != nil {
		return err
	}
	return m.origin.Send(&rpcintf.MsgAgent{
		Type: typename,
		Body: apb,
	})
}
