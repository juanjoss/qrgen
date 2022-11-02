package main

import (
	"os"

	grpc "github.com/juanjoss/qrgen/pkg/grpc/server"
)

func main() {
	done := make(chan bool)

	// creating and running gRPC server
	server := grpc.NewServer(os.Getenv("GRPC_SERVER_PORT"))
	server.ListenAndServe()

	<-done
}
