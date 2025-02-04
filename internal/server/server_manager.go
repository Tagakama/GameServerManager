package server

import (
	"GameServerManager/config"
	"fmt"
	"log"
	"net"
	"net/http"
	"os/exec"
	"strconv"
	"sync"
)

// Server represents a game server.
type Server struct {
	ID         int
	Status     string
	PID        int
	IP         string
	Port       int
	MapName    string
	AppVersion string
}

var (
	Servers      []Server
	serversMutex sync.Mutex
	nextServerID = 1
)

// HandleTCPMessage handles the TCP communication to start or find a server.
func HandleTCPMessage(cfg *config.Config, conn net.Conn, id string, message string, mapName string, appVersion string) {
	if message == "START_SERVER" {
		StartServer(cfg, conn, id, mapName, appVersion)
	} else {
		fmt.Fprintf(conn, "%s:INVALID_COMMAND\n", id)
	}
}

// StartServer starts a new server or finds an existing free one.
func StartServer(cfg *config.Config, conn net.Conn, id string, mapName string, appVersion string) {
	serversMutex.Lock()
	defer serversMutex.Unlock()

	newServer := Server{
		ID:         nextServerID,
		Status:     "running",
		IP:         "Server_" + strconv.Itoa(nextServerID),
		Port:       2797 + nextServerID,
		MapName:    mapName,
		AppVersion: appVersion,
	}
	nextServerID++
	Servers = append(Servers, newServer)

	go launchGameServer(cfg, newServer)

	fmt.Fprintf(conn, "%s:NEW_SERVER\n", id)
}

// launchGameServer launches the game server process.
func launchGameServer(cfg *config.Config, server Server) {
	port, errr := FindFreePort()
	if errr != nil {
		log.Fatalf("Failed to find free port: %v", errr)
	}
	server.Port = port

	logFilePath := fmt.Sprintf("Logs/Engine_%d.log", server.Port)
	cmd := exec.Command(cfg.LocalStorage.Directory+server.AppVersion+cfg.LocalStorage.Name,
		"-nographics", "-dedicatedServer", "-batchmode", "-fps", "60", "-dfill", "-UserID", string(server.IP+strconv.Itoa(server.Port)), "-sessionName", string(server.IP+strconv.Itoa(server.Port)), "-logFile", logFilePath,
		"-port", strconv.Itoa(server.Port), "-region eu",
		"-serverName", server.IP, "-scene", server.MapName)

	err := cmd.Start()
	if err != nil {
		fmt.Printf("Failed to start server %d: %v\n", server.ID, err)
		return
	}

	serversMutex.Lock()
	for i, s := range Servers {
		if s.ID == server.ID {
			Servers[i].PID = cmd.Process.Pid
			fmt.Printf("PID: %d Game server %d started on port %d. App server version %s .Map settings - %s .\n", Servers[i].PID, Servers[i].ID, Servers[i].Port, Servers[i].AppVersion, Servers[i].MapName)
			break
		}
	}
	serversMutex.Unlock()

	go func() {
		err := cmd.Wait()
		if err != nil {
			fmt.Printf("Server %d stopped with error: %v\n", server.ID, err)
		}
	}()
}

// StopServerByID stops a server by its ID.
func StopServerByID(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Server : %v  and %v", w, r)
	// Your existing stop
}

func FindFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:0")
	if err != nil {
		return 0, err
	}
	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer listener.Close()
	return listener.Addr().(*net.TCPAddr).Port, nil
}
