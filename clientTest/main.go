package main

import (
	"log"

	"github.com/djavorszky/rlog"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("127.0.0.1:1338", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed dialing server: %v", err)
	}
	defer conn.Close()

	client := rlog.NewLogClient(conn)

	resp, err := client.Register(context.Background(), &rlog.RegisterRequest{App: "testApp", Service: "testService"})
	if err != nil {
		log.Fatalf("failed calling 'Register' on client: %v", err)
	}

	_, err = client.Fatal(context.Background(), &rlog.LogMessage{Id: resp.Id, Message: "YO MAMMA SO FAT"})
	if err != nil {
		log.Fatalf("failed calling 'INFO' on client: %v", err)
	}
}
