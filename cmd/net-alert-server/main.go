package main

import (
	"net-alert/pkg/server"
	"net-alert/pkg/sniffer"
	"os"
)

func main() {
	if len(os.Args) <= 1 {
		panic("No device name")
	}
	var server server.Server
	go sniffer.GetData()
	server.Start(os.Args[1])
}
