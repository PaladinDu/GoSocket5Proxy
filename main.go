package main

import (
	"GoSocket5Proxy/Socket5Proxy"
	"net"
	"os"
	"time"
)

func main() {
	if len(os.Args) != 4 || os.Args[1] == "help" {
		println("Socket5Proxy {listenAddr} {userID} {password}")
		return
	}
	listenAddr := os.Args[1]
	userID := os.Args[2]
	password := os.Args[3]
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
		go Socket5Proxy.Socket5Proxy(client, userID, password)
	}
}
