package rpc

import "github.com/ab36245/go-msgpack"

type Output struct {
	*msgpack.Encoder
	Flags Flags
	call  *Call
}

func (o Output) Send() {
	o.call.send(o)
}

func (o Output) SendAndClose() {
	o.Flags |= ClosedFlag
	o.Send()
}
