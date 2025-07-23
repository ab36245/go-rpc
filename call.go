package rpc

import (
	"iter"

	"github.com/ab36245/go-msgpack"
	"github.com/rs/zerolog/log"
)

type Call struct {
	id     uint
	input  chan Input
	output chan<- []byte
	server *Server
}

func (c *Call) Input() iter.Seq[Input] {
	return func(yield func(Input) bool) {
		log.Debug().Msg("starting iterator")
		for {
			input, ok := <-c.input
			if !ok {
				log.Debug().Msg("ok is false")
				break
			}
			if !yield(input) {
				log.Debug().Msg("yield returned false")
				break
			}
		}
		log.Debug().Msg("stopping iterator")
	}
}

func (c *Call) Output() Output {
	return Output{
		Encoder: msgpack.NewEncoder(),
		Flags:   0,
		call:    c,
	}
}

func (c *Call) close() {
	c.server.control <- controlMessage{
		kind: cmCloseCall,
		id:   c.id,
	}
}

func (c *Call) recv(flags Flags, mpd *msgpack.Decoder) {
	input := Input{
		Decoder: mpd,
		call:    c,
		Flags:   flags,
	}
	c.input <- input
}

func (c *Call) send(output Output) {
	mpe := msgpack.NewEncoder()
	mpe.PutUint(uint64(c.id))
	mpe.PutUint(uint64(output.Flags))
	mpe.PutBytes(output.Bytes())
	c.server.send(mpe.Bytes())
	if output.Flags&ClosedFlag == ClosedFlag {
		c.close()
	}
}
