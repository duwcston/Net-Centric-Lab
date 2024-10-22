package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
)

const (
	HOST = "localhost"
	PORT = "8080"
	TYPE = "tcp"
)

type User struct {
	Username string   `json:"username"`
	Password string   `json:"password"`
	Fullname string   `json:"fullname"`
	Email    []string `json:"email"`
	Address  []string `json:"address"`
}

var randomKey = rand.Intn(1000) + 1
var key = strconv.Itoa(randomKey)

func main() {
	listen, err := net.Listen(TYPE, HOST+":"+PORT)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	// close listener
	defer listen.Close()
	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}

		go handleRequest(conn)
	}
}

func handleRequest(conn net.Conn) {
	defer conn.Close()

	// incoming request
	buffer := make([]byte, 1024)
	_, err := conn.Read(buffer)
	if err != nil {
		log.Fatal(err)
	}

	clientMsg := string(buffer)

	if clientMsg[:5] == "Login" {
		username, password := handleLogin(clientMsg)

		if authenticateUser(username, password) {
			fmt.Println(key + "_Hello Server from client")
			_, err = conn.Write([]byte(key + "_Authentication successful\n"))
			conn.Read(buffer) // Read the next message from the client
			clientMsg := strings.TrimSpace(string(buffer))
			if strings.HasPrefix(clientMsg, "Game") {
				startGuessingGame(err, conn)
			} else if strings.HasPrefix(clientMsg, "Download") {
				startFileTransfer(conn)
			} else if strings.HasPrefix(clientMsg, "Exit") {
				return
			}
		} else {
			fmt.Println(key + "_Login failed")
			_, err = conn.Write([]byte(key + "_Login failed"))
		}

	} else if strings.HasPrefix(clientMsg, "Register") {
		handleRegister(clientMsg)
		conn.Write([]byte(key + "_Register successful"))
	}
}

func encryptPassword(password string) string {
	return base64.StdEncoding.EncodeToString([]byte(password))
}

// func decryptPassword(encryptedPassword string) string {
// 	decoded, _ := base64.StdEncoding.DecodeString(encryptedPassword)
// 	return string(decoded)
// }

func loadUsers(filename string) []User {
	jsonFile, err := os.Open(filename)

	// Check if the file exists
	if err != nil {
		fmt.Println("Error opening file:", err)
		os.Exit(1)
	}

	data, _ := io.ReadAll(jsonFile)
	var users []User
	json.Unmarshal(data, &users)
	return users
}

func authenticateUser(username, password string) bool {
	users := loadUsers("User.json")
	encryptedPassword := encryptPassword(password)
	for _, user := range users {
		if user.Username == username && user.Password == encryptedPassword {
			return true
		}
	}
	return false
}

func handleLogin(clientMsg string) (username, password string) {
	data := strings.Split(clientMsg, "|")
	username = strings.TrimSpace(data[1])
	password = strings.Trim(data[2], "\x00")

	return username, password
}

func saveUser(user User, filename string) {
	users := loadUsers(filename)
	users = append(users, user)
	data, _ := json.Marshal(users)
	_ = os.WriteFile(filename, data, 0644)
}

func handleRegister(clientMsg string) {
	data := strings.Split(clientMsg, "|")
	username := strings.TrimSpace(data[1])
	password := strings.Trim(data[2], "\x00")
	fullname := strings.TrimSpace(data[3])
	email := strings.Split(strings.TrimSpace(data[4]), ",")
	address := strings.Split(strings.Trim(data[5], "\x00"), ",")

	user := User{
		Username: username,
		Password: encryptPassword(password),
		Fullname: fullname,
		Email:    email,
		Address:  address,
	}

	saveUser(user, "User.json")
}

func startGuessingGame(err error, conn net.Conn) {
	for {
		// Send the message to the client, to start the game
		conn.Write([]byte("Start"))

		// Generate random number
		randomNumber := rand.Intn(100) + 1
		target := randomNumber
		fmt.Println("Target number:", target)
		for {
			received := make([]byte, 1024)
			_, err = conn.Read(received)
			if err != nil {
				fmt.Println("Error reading from client:", err)
				break
			}

			clientMsg := strings.TrimSpace(string(received[:]))

			if strings.Contains(clientMsg, "Exit") {
				break
			}

			if strings.Contains(clientMsg, "Again") {
				break
			}

			// Convert received data to integer
			// guessNumber, err := strconv.Atoi(string(received))
			guessNumber := 0
			fmt.Sscanf(string(received), "%d", &guessNumber)
			fmt.Println("Guess number:", guessNumber)

			if guessNumber < target {
				conn.Write([]byte(key + "_Your guess is too low. Try again."))

			} else if guessNumber > target {
				conn.Write([]byte(key + "_Your guess is too high. Try again."))

			} else {
				conn.Write([]byte(key + "_Congratulations! You guessed the number."))
				break
			}
		}
	}
}

func startFileTransfer(conn net.Conn) {
	// Read the requested file name from the client
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println(key+"_Error reading from client:", err)
		return
	}
	fileName := string(buffer[:n])
	fmt.Println("Client requested file:", fileName)

	// Check if the file exists on the server
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		conn.Write([]byte("Error: File not found\n"))
		fmt.Println(key+"_File not found:", fileName)
		return
	}

	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println("Error opening file:", err)
		conn.Write([]byte("Error: Unable to open file\n"))
		return
	}
	defer file.Close()

	conn.Write([]byte("File download starting...\n"))

	// Send file data
	buffer = make([]byte, 1024)
	for {
		n, err := file.Read(buffer)
		if err != nil && err != io.EOF {
			fmt.Println("Error reading file:", err)
			break
		}
		if n == 0 {
			break
		}
		conn.Write(buffer[:n]) // Send the context to the client
	}

	conn.Write([]byte("\nFile download complete\n"))
	fmt.Println("File download completed for:", fileName)
}
