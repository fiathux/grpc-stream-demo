package main

import (
	"grpcdemo/pbmsg"
	"testing"
)

func TestRegHandler(t *testing.T) {
	// TestRegHandler is a test for RegHandler
	err := RegHandler(&Msger{}, "somehandler",
		func(owner MsgerStream, m *pbmsg.Msg_A, err error) {
		})
	if err != nil {
		t.Errorf("RegHandler failed: %v", err)
	}
}
