package bookstore

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/golang/glog"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var (
	bookstoreEndpoint = flag.String("bookstoreEndpoint", "localhost:50051", "endpoint of YourService")
)

func run() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := RegisterBookstoreHandlerFromEndpoint(ctx, mux, *bookstoreEndpoint, opts)
	if err != nil {
		return err
	}

	fmt.Print("\nProxy listening on 8081\n")
	return http.ListenAndServe(":8081", mux)
}

func RunProxy() {
	flag.Parse()
	defer glog.Flush()
	if err := run(); err != nil {
		glog.Fatal(err)
	}
}
