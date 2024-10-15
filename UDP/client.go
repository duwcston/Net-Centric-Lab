package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	serverEP := "127.0.0.1"
	if len(os.Args) > 1 {
		serverEP = os.Args[1]
	}
	if !strings.Contains(serverEP, ":") {
		serverEP = fmt.Sprintf("%v:8000", serverEP)
	}

	conn, err := net.Dial("udp", serverEP)
	if err != nil {
		fmt.Printf("Dial err %v", err)
		os.Exit(-1)
	}
	defer conn.Close()

	for {
		msg := "Hello, UDP server"
		fmt.Printf("Ping: %v\n", msg)
		if _, err = conn.Write([]byte(msg)); err != nil {
			fmt.Printf("Write err %v", err)
			os.Exit(-1)
		}

		p := make([]byte, 1024)
		nn, err := conn.Read(p)
		if err != nil {
			fmt.Printf("Read err %v\n", err)
			os.Exit(-1)
		}

		fmt.Printf("%v\n", string(p[:nn]))
	}
}
