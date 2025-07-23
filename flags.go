package rpc

import "fmt"

type Flags uint8

const (
	// These must agree with the rpc client
	NewFlag    Flags = 0x01
	CloseFlag  Flags = 0x02
	ClosedFlag Flags = 0x04
)

func (f Flags) String() string {
	s := fmt.Sprintf("%#02x", uint8(f))
	if f&NewFlag == NewFlag {
		s += " New"
	}
	if f&CloseFlag == CloseFlag {
		s += " Close"
	}
	if f&ClosedFlag == ClosedFlag {
		s += " Closed"
	}
	return s
}
