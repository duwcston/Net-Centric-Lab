package main

import (
	"fmt"
	"net"
	"strings"
	"sync"
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

var (
	clients   = make(map[string]*Client)
	clientsMu sync.Mutex
)

func main() {
	serverAddr := HOST + ":" + PORT

	// Resolve server address
	addr, err := net.ResolveUDPAddr("udp4", serverAddr)
	if err != nil {
		fmt.Println("Error resolving address: ", err)
		return
	}

	// Listen on UDP
	conn, err := net.ListenUDP("udp4", addr)
	if err != nil {
		fmt.Println("Error listening: ", err)
		return
	}
	defer conn.Close()

	buffer := make([]byte, 1024)
	fmt.Println("Server listening on:", serverAddr)

	for {
		n, clientAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("Error reading from UDP: ", err)
			continue
		}
		message := strings.TrimSpace(string(buffer[:n]))
		fmt.Println("Received: ", message, " from ", clientAddr)

		// Process commands
		go handleCommand(conn, clientAddr, message)
	}
}

func handleCommand(conn *net.UDPConn, addr *net.UDPAddr, message string) {
	parts := strings.SplitN(message, " ", 2)
	command := parts[0]
	switch command {
	case "LOGIN":
		if len(parts) < 2 {
			return
		}
		username := parts[1]
		registerClient(username, addr)
		conn.WriteToUDP([]byte("Welcome "+username), addr)
	case "LOGOUT":
		if len(parts) < 2 {
			return
		}
		username := parts[1]
		removeClient(username)
		conn.WriteToUDP([]byte("Goodbye "+username), addr)
	case "MSG":
		if len(parts) < 2 {
			return
		}
		handleMessage(conn, addr, parts[1])
	}
}

func registerClient(name string, addr *net.UDPAddr) {
	clientsMu.Lock()
	defer clientsMu.Unlock()
	clients[name] = &Client{Name: name, Addr: addr}
	fmt.Println("Registered client: ", name, addr)
}

func removeClient(name string) {
	clientsMu.Lock()
	defer clientsMu.Unlock()
	delete(clients, name)
	fmt.Println("Removed client: ", name)
}

func handleMessage(conn *net.UDPConn, senderAddr *net.UDPAddr, message string) {
	clientsMu.Lock()
	defer clientsMu.Unlock()

	var senderName string
	for name, client := range clients {
		if client.Addr.String() == senderAddr.String() {
			senderName = name
			break
		}
	}

	// Check for @<username> or @all command
	if strings.HasPrefix(message, "@") {
		parts := strings.SplitN(message, " ", 2)
		if len(parts) < 2 {
			return
		}
		target := parts[0][1:] // Remove "@" prefix
		msg := parts[1]

		// Private message to specific user
		if target != "all" {
			if client, ok := clients[target]; ok {
				privateMsg := fmt.Sprintf("Private from %s: %s", senderName, msg)
				conn.WriteToUDP([]byte(privateMsg), client.Addr)
			} else {
				conn.WriteToUDP([]byte("User "+target+" not found"), senderAddr)
			}
			return
		}

		// Broadcast message to all users
		broadcastMsg := fmt.Sprintf("Broadcast from %s: %s", senderName, msg)
		for _, client := range clients {
			if client.Addr.String() != senderAddr.String() { // Exclude sender
				conn.WriteToUDP([]byte(broadcastMsg), client.Addr)
			}
		}
	}
}
