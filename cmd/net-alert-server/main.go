package main

import (
	"net-alert/pkg/server"
	"os"
)

func main() {
	if len(os.Args) <= 1 {
		panic("No device name")
	}
	var server server.Server
	server.Start(os.Args[1])
}
