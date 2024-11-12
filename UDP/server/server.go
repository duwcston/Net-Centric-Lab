package main

import (
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"strings"
)

const (
	HOST = "localhost"
	PORT = "8080"
	TYPE = "udp"
)

type Client struct {
	Name string
	Addr *net.UDPAddr
}

var clients = make(map[string]*Client)

func random(min, max int) int {
	return rand.Intn(max-min) + min
}

func main() {
	arguments := []string{HOST + ":" + PORT}
	if len(arguments) == 1 {
		fmt.Println("Please provide a port number!")
		return
	}
	PORT := ":" + arguments[1]

	s, err := net.ResolveUDPAddr("udp4", PORT)
	if err != nil {
		fmt.Println(err)
		return
	}

	connection, err := net.ListenUDP("udp4", s)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer connection.Close()
	buffer := make([]byte, 1024)

	for {
		n, addr, err := connection.ReadFromUDP(buffer)
		fmt.Print("-> ", string(buffer[0:n-1]))

		// If the client is not in the map, add it
		if clients[string(buffer[0:n-1])] == nil {

		}

		if strings.TrimSpace(string(buffer[0:n])) == "STOP" {
			fmt.Println("Exiting UDP server!")
			return
		}

		data := []byte(strconv.Itoa(random(1, 1001)))
		fmt.Printf("data: %s\n", string(data))
		_, err = connection.WriteToUDP(data, addr)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}
