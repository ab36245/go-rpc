package rpc

import (
	"github.com/ab36245/go-msgpack"
)

type Input struct {
	*msgpack.Decoder
	Flags Flags
	call  *Call
}
