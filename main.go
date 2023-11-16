package main

import (
	"GoSocket5Proxy/Socket5Proxy"
	"net"
	"os"
	"time"
)

func main() {
	if len(os.Args) != 2 || os.Args[1] == "help" {
		println("Socket5Proxy {listenAddr}")
		return
	}
	listenAddr := os.Args[1]
	socket, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return
	}
	var tempDelay time.Duration
	for {
		client, err := socket.Accept()

		if err != nil {
			if ne, ok := err.(interface {
				Temporary() bool
			}); ok && ne.Temporary() {
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}
				if max := 1 * time.Second; tempDelay > max {
					tempDelay = max
				}
				timer := time.NewTimer(tempDelay)
				select {
				case <-timer.C:
				}
				continue
			}
			break
		}
		go Socket5Proxy.Socket5Proxy(client)
	}
}
