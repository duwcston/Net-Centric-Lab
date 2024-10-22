package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

const (
	HOST = "localhost"
	PORT = "8080"
	TYPE = "tcp"
)

func main() {
	tcpServer, err := net.ResolveTCPAddr(TYPE, HOST+":"+PORT)

	if err != nil {
		println("ResolveTCPAddr failed:", err.Error())
		os.Exit(1)
	}

	conn, err := net.DialTCP(TYPE, nil, tcpServer)
	if err != nil {
		println("Dial failed:", err.Error())
		os.Exit(1)
	}

	defer conn.Close() // Ensure connection is closed when the main function returns

	// Login or register
	fmt.Println("1. Login")
	fmt.Println("2. Register")
	fmt.Println("3. Exit")
	fmt.Print("Enter your choice: ")
	var choice int
	fmt.Scanln(&choice)

	switch choice {
	case 1:
		username, password := login()

		// Send the login data to the server
		_, err = conn.Write([]byte("Login" + "|" + username + "|" + password))
		if err != nil {
			println("Write data failed:", err.Error())
			os.Exit(1)
		}

		received := make([]byte, 1024)
		_, err = conn.Read(received)
		if err != nil {
			println("Read data failed:", err.Error())
			os.Exit(1)
		}
		fmt.Println(string(received))

		if strings.Contains(string(received), "successful") {
			fmt.Println("Hello client from server")
			fmt.Println("1. Play Game")
			fmt.Println("2. Download File")
			fmt.Println("3. Exit")

			fmt.Print("Enter your choice: ")
			var choice int
			fmt.Scanln(&choice)

			switch choice {
			case 1:
				_, err = conn.Write([]byte("Game"))
				fmt.Println("Welcome to Guessing Game!")
				playGame(err, conn)

			case 2:
				_, err = conn.Write([]byte("Download"))
				fmt.Print("Enter the name of the file to download: ")
				var fileName string
				fmt.Scanln(&fileName)

				// Send download request
				_, err = conn.Write([]byte(fileName))
				if err != nil {
					fmt.Println("Error sending file request:", err)
					return
				}

				// Receive the file data
				requestFile(conn, fileName)
			case 3:
				_, err = conn.Write([]byte("Exit"))
				exit()
			default:
				println("Invalid choice")
				os.Exit(0)
			}
		}

	case 2:
		username, password, fullname, email, address := register()

		// Send the register data to the server
		_, err = conn.Write([]byte("Register" + "|" + username + "|" + password + "|" + fullname + "|" + email + "|" + address))
		if err != nil {
			println("Write data failed:", err.Error())
			os.Exit(1)
		}
	case 3:
		exit()
	default:
		println("Invalid choice")
		return
	}

	conn.Close()

}

func login() (username, password string) {
	// Login to server
	fmt.Print("Enter username: ")
	// var username string
	fmt.Scanln(&username)

	fmt.Print("Enter password: ")
	// var password string
	fmt.Scanln(&password)

	return username, password
}

func register() (username, password, fullname, email, address string) {
	// Register to server
	fmt.Print("Enter username: ")
	fmt.Scanln(&username)

	fmt.Print("Enter password: ")
	fmt.Scanln(&password)

	fmt.Print("Enter fullname: ")
	fmt.Scanln(&fullname)

	fmt.Print("Enter email: ")
	fmt.Scanln(&email)

	fmt.Print("Enter address: ")
	fmt.Scanln(&address)

	return username, password, fullname, email, address
}

func exit() {
	return
}

func playGame(err error, conn *net.TCPConn) {
	received := make([]byte, 1024)
	_, err = conn.Read(received)
	for {
		fmt.Printf("Guess a number between 1 and 100 (press 0 to exit): ")
		var guessNumber int
		fmt.Scanln(&guessNumber)

		if guessNumber == 0 {
			_, err = conn.Write([]byte("Exit"))
			return
		}

		_, err = conn.Write([]byte(fmt.Sprintf("%d", guessNumber)))
		if err != nil {
			println("Write data failed:", err.Error())
			os.Exit(1)
		}

		// Read the response from the server
		received := make([]byte, 1024)
		_, err = conn.Read(received)
		if err != nil {
			println("Read data failed:", err.Error())
			os.Exit(1)
		}

		response := strings.TrimSpace(string(received[:]))
		fmt.Printf("Received response: '%s'\n", response)

		if strings.Contains(response, "_Congratulations!") {

			fmt.Println("Do you want to play again? (y/n)")
			var playAgain string
			fmt.Scanln(&playAgain)

			if strings.ToLower(playAgain) == "n" {
				_, err = conn.Write([]byte("Exit")) // Send "Exit" to the server
				break
			} else if strings.ToLower(playAgain) == "y" {
				// If the user says yes, send "Again" to the server
				_, err = conn.Write([]byte("Again"))
				_, err = conn.Read(received) // Wait for server's response to start the game again
				fmt.Println("Starting a new game...")
				conn.Read(received)
			}
		}
	}
}

func requestFile(conn net.Conn, fileName string) {
	// Create or open the file for writing (only if the file exists)
	outFile, err := os.Create("downloaded_" + fileName)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer outFile.Close()

	// Buffer to receive the file data
	buffer := make([]byte, 1024)

	// Read the response from the server
	for {
		n, err := conn.Read(buffer)
		// Handle server responses
		serverResponse := string(buffer[:n])
		// fmt.Print("Server response: ", serverResponse)
		if err != nil {
			if err == io.EOF {
				fmt.Println("File download completed.")
				break
			}
			fmt.Println("Error reading from server:", err)
			return
		}

		// Ignore notification messages from the server
		if serverResponse == "File download starting...\n" || serverResponse == "\nFile download complete\n" {
			continue
		}
		if serverResponse == "Error: File not found\n" {
			fmt.Println("Server response: ", serverResponse)
			return // Exit if the file does not exist
		}

		// Write file data to the local file if it's not an error message
		outFile.Write(buffer[:n])
	}
}
