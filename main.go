package main

import (
	"embed"
	"os"

	api "github.com/juanjoss/qrgen/pkg/api"
	grpc "github.com/juanjoss/qrgen/pkg/grpc/server"
)

//go:embed view/*
var fs embed.FS

func main() {
	done := make(chan bool)

	// running grpc server
	go func() {
		grpc.ListenAndServe(os.Getenv("GRPC_SERVER_PORT"))
	}()

	// creating and running the http server
	s := api.NewServer(os.Getenv("APP_PORT"), fs)
	s.ListenAndServe()

	<-done
}
