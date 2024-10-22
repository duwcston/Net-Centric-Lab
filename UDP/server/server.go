package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
)

func main() {
	serverPort := 8000
	if len(os.Args) > 1 {
		if v, err := strconv.Atoi(os.Args[1]); err != nil {
			fmt.Printf("Invalid port %v, err %v", os.Args[1], err)
			os.Exit(-1)
		} else {
			serverPort = v
		}
	}

	addr := net.UDPAddr{
		Port: serverPort,
		IP:   net.ParseIP("0.0.0.0"),
	}
	server, err := net.ListenUDP("udp", &addr)
	if err != nil {
		fmt.Printf("Listen err %v\n", err)
		os.Exit(-1)
	}
	fmt.Printf("Listen at %v\n", addr.String())

	for {
		p := make([]byte, 1024)
		nn, raddr, err := server.ReadFromUDP(p)
		if err != nil {
			fmt.Printf("Read err  %v", err)
			continue
		}

		msg := p[:nn]
		fmt.Printf("Received %v %s\n", raddr, msg)

		go func(conn *net.UDPConn, raddr *net.UDPAddr, msg []byte) {
			_, err := conn.WriteToUDP([]byte(fmt.Sprintf("Pong: %s", msg)), raddr)
			if err != nil {
				fmt.Printf("Response err %v", err)
			}
		}(server, raddr, msg)
	}
}
