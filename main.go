package main

import (
	"embed"
	"os"

	grpc "github.com/juanjoss/qrgen/handlers/grpc/server"
	"github.com/juanjoss/qrgen/handlers/http"
)

//go:embed views/*
var fs embed.FS

func main() {
	done := make(chan bool)

	// running grpc server
	go func() {
		grpc.ListenAndServe(os.Getenv("GRPC_SERVER_PORT"))
	}()

	// running http server
	http.ListenAndServe(os.Getenv("APP_PORT"), fs)

	<-done
}
