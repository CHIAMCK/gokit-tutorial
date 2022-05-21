package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"

	pb "github.com/go-kit-tutorial/pb"
)

const (
	address     = "localhost:50051"
	defaultName = "world"
)

func main() {
	// set up a connection to the server
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	fmt.Println(err)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	// close connection after finish
	defer conn.Close()

	// create a client to access a gRPC service
	c := pb.NewMathServiceClient(conn)

	// name := defaultName
	// if len(os.Args) > 1 {
	// 	name =
	// }

	// set timeout and cancel the request after certain period of time
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	// Calling the CancelFunc cancels the child and its children,
	// removes the parent’s reference to the child, and stops any associated timers
	defer cancel()
	// client.methodName(define in proto file)
	r, err := c.Add(ctx, &pb.MathRequest{
		NumA: 13,
		NumB: 12,
	})
	fmt.Println(err)
	if err != nil {
		log.Fatalf("could not .....: %v", err)
	}
	fmt.Println("ending")
	log.Printf("Calculating: %f", r.GetResult())
}

// create a client connection to the given target
// initialize service client, pass client connection to it
// call method on service client
