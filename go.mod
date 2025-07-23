module github.com/ab36245/go-rpc

go 1.24.4

replace github.com/ab36245/go-websocket => ../go-websocket

replace github.com/ab36245/go-errors => ../go-errors

require (
	github.com/ab36245/go-msgpack v0.0.0-20250708012415-aa1507c506e3
	github.com/rs/zerolog v1.34.0
)

require (
	github.com/ab36245/go-errors v0.0.0-20250428061939-8b056c3b905e // indirect
	github.com/ab36245/go-websocket v0.0.0-20250714021031-87e7ab40c492 // indirect
	github.com/gorilla/websocket v1.5.3 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	golang.org/x/sys v0.12.0 // indirect
)

replace github.com/ab36245/go-msgpack => ../go-msgpack
