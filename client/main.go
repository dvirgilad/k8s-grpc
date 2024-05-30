package main

import (
	"context"
	"flag"
	"fmt"
	pb "github.com/dvirgilad/grpcNode/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"
)

var (
	addr = flag.String("addr", "localhost:9999", "the address to connect to")
)

func main() {
	flag.Parse()
	// Set up a connection to the server.
	conn, err := grpc.NewClient(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewNodeServiceClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	r, err := c.GetNodes(ctx, &pb.NodeRequest{})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	fmt.Println("Name\t version\t Ready")
	for _, node := range r.GetNodes() {
		fmt.Printf("%s \t %s \t %t \t", node.Name, node.Version, node.Ready)
		fmt.Println()
	}
	if err != nil {
		log.Fatalf("could not greet again: %v", err)
	}
	log.Print("got nodes")

}
