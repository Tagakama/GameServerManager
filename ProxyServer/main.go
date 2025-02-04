package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
)

var (
	adminConn     net.Conn
	isAdminActive bool
	mu            sync.Mutex
	clientConnMap = make(map[string]net.Conn) // Stores pairs Client - Administrator by client ID
)

func main() {
	// We start a server that listens on port 8088
	listener, err := net.Listen("tcp", ":8088")
	if err != nil {
		fmt.Println("Error starting server:", err)
		os.Exit(1)
	}
	defer listener.Close()

	fmt.Println("Server is listening on port 8088")

	for {
		// We are waiting for clients to connect
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		go handleConnection(conn)
	}
}

// Handling client connections
func handleConnection(conn net.Conn) {
	defer conn.Close()
	fmt.Println("New client connected:", conn.RemoteAddr())

	scanner := bufio.NewScanner(conn)
	if scanner.Scan() {
		code := scanner.Text()

		if code == "admin" {
			// Establishing an administrator connection
			mu.Lock()
			isAdminActive = true
			adminConn = conn
			mu.Unlock()

			fmt.Println("Admin connected:", conn.RemoteAddr())
			handleAdminConnection(conn)
		} else {
			// We process regular requests, identifying the client with a code
			handleClientRequest(conn, code)
		}
	}
}

// Processing requests from the administrator and forwarding the response to the client
func handleAdminConnection(conn net.Conn) {
	defer func() {
		mu.Lock()
		isAdminActive = false
		adminConn = nil
		mu.Unlock()
		fmt.Println("Admin disconnected")
	}()

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		message := scanner.Text()
		fmt.Println("Admin sent:", message)

		// We parse the message and find the client ID for the response
		clientID, response := parseAdminMessage(message)
		sendResponseToClient(clientID, response)
	}
}

// Processing requests from regular clients and relaying them to the administrator
func handleClientRequest(conn net.Conn, clientID string) {
	// We save the client connection by its ID
	mu.Lock()
	clientConnMap[clientID] = conn
	mu.Unlock()

	defer func() {
		mu.Lock()
		delete(clientConnMap, clientID)
		mu.Unlock()
		fmt.Println("Client disconnected:", clientID)
	}()

	fmt.Println("Client ID:", clientID, "connected")

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		message := scanner.Text()
		fmt.Printf("Received from client %s: %s\n", clientID, message)
		complete := retransmitToAdmin(fmt.Sprintf("%s:%s", clientID, message))
		if !complete {
			sendResponseToClient(clientID, "SERVER_NOT_FOUND")
		}
	}
}

// Relay message to administrator
func retransmitToAdmin(message string) bool {
	mu.Lock()
	defer mu.Unlock()

	if isAdminActive && adminConn != nil {
		_, err := adminConn.Write([]byte(message + "\n"))
		if err != nil {
			fmt.Println("Error sending to admin:", err)
			isAdminActive = false
			adminConn = nil
			return false
		}
		return true
	} else {

		fmt.Println("No active admin to retransmit message")
		return false
	}
}

// Parses the admin message to obtain the client ID and response
func parseAdminMessage(message string) (string, string) {
	parts := strings.SplitN(message, ":", 2)
	if len(parts) < 2 {
		// If the message format is incorrect, return empty values
		return "", message
	}
	clientID := strings.TrimSpace(parts[0])
	response := strings.TrimSpace(parts[1])
	return clientID, response
}

// Sending a response to a client by ID
func sendResponseToClient(clientID, response string) {
	mu.Lock()
	clientConn, exists := clientConnMap[clientID]
	mu.Unlock()

	if exists {
		_, err := clientConn.Write([]byte(response + "\n"))
		if err != nil {
			fmt.Println("Error sending response to client:", err)
		} else {
			fmt.Printf("Response sent to client %s: %s\n", clientID, response)
		}
	} else {
		fmt.Printf("Client %s not found or disconnected\n", clientID)
	}
}
