package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/djavorszky/disco"

	"github.com/djavorszky/rlog"

	"google.golang.org/grpc"
)

func main() {
	start()
}

func start() {
	port := flag.Int("port", 1338, "Define port on which server should listen")
	logFile := flag.String("log", "logregator.log", "Set the file in which logging will happen")
	flag.Parse()

	fmt.Printf("Starting to listen to tcp connections on port %d\n", *port)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	defer lis.Close()

	file, err := os.OpenFile(*logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		rlog.SetOut(file)
	} else {
		fmt.Printf("Using default stderr because: %v\n", err)
	}
	defer file.Close()

	grpcServer := grpc.NewServer()

	rlog.RegisterLogServer(grpcServer, &rlog.Server{})

	err = disco.Announce("224.0.0.1:9999", fmt.Sprintf("192.168.221.219:%d", *port), "rlog")
	if err != nil {
		log.Fatalf("Could not announce myself: %v", err)
	}

	grpcServer.Serve(lis)
}
