package handlers

import (
	"GameServerManager/config"
	"GameServerManager/internal/server"
	"bufio"
	"fmt"
	"net"
	"strings"
	"time"
)

var tcpRequests []Request // Query storage
type Request struct {
	IP   string
	Time string
}

// HandleConnection processes the incoming TCP request.
func HandleConnection(cfg *config.Config, conn net.Conn) {
	defer conn.Close()
	var IDclient, message, mapName, appVersion string

	reader := bufio.NewReader(conn)
	rawMessage, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading from connection:", err)
		return
	}
	fmt.Println("Full message received:", rawMessage)

	// Parse message to separate IDclient and message
	IDclient, message, mapName, appVersion = parseAdminMessage(strings.TrimSpace(rawMessage))
	fmt.Printf("Parsed IDclient: %s, message: %s\n", IDclient, message)

	if message == "START_SERVER" {
		newRequest := Request{
			IP:   conn.RemoteAddr().String(),
			Time: time.Now().Format("2006-01-02 15:04:05"),
		}
		tcpRequests = append(tcpRequests, newRequest)
	}

	// Process the message
	server.HandleTCPMessage(cfg, conn, IDclient, message, mapName, appVersion)
}

func parseAdminMessage(message string) (string, string, string, string) {
	parts := strings.SplitN(message, ":", 4)
	if len(parts) < 4 {
		// If the message format is incorrect, return blank values.
		return "", message, "", ""
	}
	clientID := strings.TrimSpace(parts[0])
	response := strings.TrimSpace(parts[1])
	mapName := strings.TrimSpace(parts[2])
	appVersion := strings.TrimSpace(parts[3])
	return clientID, response, mapName, appVersion
}
