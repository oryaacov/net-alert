package main

import (
	"net-alert/pkg/server"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	if len(os.Args) <= 1 {
		panic("No device name")
	}
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	var server server.Server
	go func() {
		<-c
		server.Exit()
		os.Exit(1)
	}()
	server.Start(os.Args[1])

}
