package rpc

type controlMessage struct {
	kind controlMessageKind
	id   uint
}

type controlMessageKind int

const (
	cmCloseAll controlMessageKind = iota
	cmCloseCall
)
