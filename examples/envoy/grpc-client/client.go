package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/LorenzHW/gnostic-grpc/examples/envoy/bookstore"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
	"log"
	"time"
)

var (
	serverAddr = flag.String("server_addr", "127.0.0.1:51051", "The server address in the format of host:port")
)

func main() {
	flag.Parse()

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	conn, err := grpc.Dial(*serverAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()

	client := bookstore.NewBookstoreClient(conn)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	res, err := client.ListShelves(ctx, &empty.Empty{})
	if res != nil {
		fmt.Println("The themes of your shelves:")
		for _, shelf := range res.Ok.Shelves {
			fmt.Println(shelf.Theme)
		}
	}
}
