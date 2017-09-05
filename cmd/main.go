package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/djavorszky/rlog"

	"google.golang.org/grpc"
)

func main() {
	start()
}

func start() {
	port := flag.Int("port", 1338, "Define port on which server should listen")

	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	rlog.RegisterLogServer(grpcServer, &rlog.Server{})

	grpcServer.Serve(lis)
}
