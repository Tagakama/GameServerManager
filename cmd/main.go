package main

import (
	"GameServerManager/config"
	"GameServerManager/internal/handlers"
	"fmt"
	"net"
)

// Proxy connection settings

func main() {
	//Load config
	cfg := config.LoadConfig()
	// Start the TCP server
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error starting TCP server:", err)
		return
	}
	defer listener.Close()
	fmt.Println("TCP server listening on port 8080")

	// Start connection to the proxy server as admin
	go connectToProxyAsAdmin(cfg)

	// Handle incoming connections from clients
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go handlers.HandleConnection(cfg, conn)
	}
}

// Connects to the proxy server and authenticates as admin
func connectToProxyAsAdmin(cfg *config.Config) {
	for {
		fmt.Println("Connecting to proxy as admin...")

		// Try to establish a connection with the proxy
		conn, err := net.Dial("tcp", cfg.Proxy.Ip+":"+cfg.Proxy.Port)
		if err != nil {
			fmt.Println("Error connecting to proxy:", err)
			continue
		}
		go handlers.HandleConnection(cfg, conn)

		// Send the admin authentication code
		_, err = conn.Write([]byte(cfg.Proxy.Code + "\n"))
		if err != nil {
			fmt.Println("Error sending admin code:", err)
			conn.Close()
			continue
		}

		fmt.Println("Connected to proxy as admin. Listening for requests...")

		// Handle incoming proxy requests
		handleProxyRequests(conn)

		// Close the connection and wait before retrying
		conn.Close()
	}
}

// Handles incoming requests from the proxy server
func handleProxyRequests(conn net.Conn) {
	for {
		buffer := make([]byte, 1024)
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Disconnected from proxy:", err)
			break
		}
		message := string(buffer[:n])
		fmt.Println("Received request from proxy:", message)
	}
}
