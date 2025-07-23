package rpc

import (
	"errors"

	"github.com/rs/zerolog/log"

	"github.com/ab36245/go-msgpack"
	"github.com/ab36245/go-websocket"
)

func NewServer(socket websocket.Socket, handlers map[uint]Handler) *Server {
	return &Server{
		calls:    make(map[uint]*Call, 0),
		closing:  false,
		control:  make(chan controlMessage),
		handlers: handlers,
		output:   make(chan []byte),
		socket:   socket,
	}
}

type Server struct {
	calls    map[uint]*Call
	closing  bool
	control  chan controlMessage
	handlers map[uint]Handler
	output   chan []byte
	socket   websocket.Socket
}

func (s *Server) Run() {
	go s.doInput()
	go s.doOutput()
	s.doControl()
}

func (s *Server) doControl() {
	for {
		cm, ok := <-s.control
		if !ok {
			log.Error().Msg("control channel has closed unexpectedly!")
			break
		}
		switch cm.kind {
		case cmCloseCall:
			log.Trace().Uint("cid", cm.id).Msg("closing call")
			call := s.calls[cm.id]
			if call == nil {
				log.Trace().Uint("cid", cm.id).Msg("call already closed")
				continue
			}
			delete(s.calls, cm.id)
		}
	}
}

func (s *Server) doInput() {
	for {
		log.Trace().Msg("waiting for socket")
		msg, err := s.socket.Read()
		if errors.Is(err, websocket.ClosedError) {
			log.Debug().Msg("client has closed socket")
			break
		}
		if err != nil {
			log.Error().Err(err).Msg("socket read returned an error")
			break
		}
		log.Trace().Stringer("kind", msg.Kind).Msg("read message")
		if !msg.IsBinary() {
			log.Error().Stringer("kind", msg.Kind).Msg("can't handle message")
			break
		}
		mpd := msgpack.NewDecoder(msg.Data)

		log.Trace().Msg("reading call id (cid)")
		n, err := mpd.GetUint()
		if err != nil {
			log.Error().Err(err).Msg("bad cid")
			break
		}
		cid := uint(n)
		log.Trace().Uint("cid", cid).Send()

		log.Trace().Msg("reading flags")
		n, err = mpd.GetUint()
		if err != nil {
			log.Error().Err(err).Msg("bad flags")
			break
		}
		flags := Flags(n)
		log.Trace().Stringer("flags", flags).Send()

		var call *Call
		if flags&NewFlag == NewFlag {
			log.Trace().Msg("new flag set")
			if _, ok := s.calls[cid]; ok {
				log.Error().Uint("cid", cid).Msg("cid already in use")
				break
			}

			n, err := mpd.GetUint()
			if err != nil {
				log.Error().Err(err).Msg("bad hid")
				break
			}
			hid := uint(n)
			log.Trace().Uint("hid", hid).Msg("new call")

			handler := s.handlers[hid]
			if handler == nil {
				log.Error().Uint("hid", hid).Msg("no handler for hid")
				break
			}

			call = &Call{
				id:     cid,
				input:  make(chan Input),
				output: s.output,
				server: s,
			}
			s.calls[cid] = call
			go handler(call)
		} else {
			log.Trace().Msg("new flag not set")
			call = s.calls[cid]
			if call == nil {
				log.Trace().Uint("cid", cid).Msg("unknown cid")
				break
			}
		}
		call.recv(flags, mpd)
	}
	log.Trace().Msg("stopping")
	s.closing = true
}

func (s *Server) doOutput() {
	for {
		bytes, ok := <-s.output
		if !ok {
			log.Debug().Msg("output channel is closed")
			break
		}
		log.Trace().Int("bytes", len(bytes)).Msg("read bytes")
		err := s.socket.WriteBinary(bytes)
		if err != nil {
			log.Error().Err(err).Msg("writing to socket failed")
			break
		}
	}
}

func (s *Server) send(bytes []byte) {
	s.output <- bytes
}
