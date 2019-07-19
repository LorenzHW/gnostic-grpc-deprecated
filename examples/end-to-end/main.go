package main

import (
	"github.com/LorenzHW/gnostic-grpc/examples/end-to-end/bookstore"
)

func main() {
	// Run server inside goroutine so we don't block the main thread.
	go bookstore.RunServer()
	bookstore.RunProxy()
}
