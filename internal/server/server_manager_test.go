package server

import (
	"GameServerManager/config"
	"bytes"
	"net"
	"testing"
	"time"
)

// MockConn реализует net.Conn для тестирования
type MockConn struct {
	bytes.Buffer
}

func (m *MockConn) Close() error                       { return nil }
func (m *MockConn) RemoteAddr() net.Addr               { return &net.TCPAddr{} }
func (m *MockConn) LocalAddr() net.Addr                { return &net.TCPAddr{} }
func (m *MockConn) SetDeadline(t time.Time) error      { return nil }
func (m *MockConn) SetReadDeadline(t time.Time) error  { return nil }
func (m *MockConn) SetWriteDeadline(t time.Time) error { return nil }

func TestStartServer(t *testing.T) {
	// Создаем mock-соединение
	conn := &MockConn{}

	// Создаем mock-конфигурацию
	cfg := &config.Config{
		LocalStorage: struct {
			Directory string `yaml:"directory"`
			Name      string `yaml:"filename_exe"`
		}{
			Directory: "/www/",
			Name:      "map1",
		},
	}

	// Очищаем глобальное состояние перед тестом
	Servers = []Server{}
	nextServerID = 1

	// Вызываем тестируемую функцию
	StartServer(cfg, conn, "client123", "map1", "v1.0")

	// Проверяем, что сервер был добавлен
	if len(Servers) != 1 {
		t.Errorf("Expected 1 server, got %d", len(Servers))
	}

	// Проверяем состояние добавленного сервера
	server := Servers[0]
	if server.ID != 1 {
		t.Errorf("Expected server ID 1, got %d", server.ID)
	}
	if server.Status != "running" {
		t.Errorf("Expected server status 'running', got %s", server.Status)
	}
	if server.Port != 2798 { // 2797 + nextServerID (1)
		t.Errorf("Expected server port 2798, got %d", server.Port)
	}
	if server.MapName != "map1" {
		t.Errorf("Expected server map name 'map1', got %s", server.MapName)
	}
	if server.AppVersion != "v1.0" {
		t.Errorf("Expected server app version 'v1.0', got %s", server.AppVersion)
	}

	// Проверяем ответ клиенту
	expected := "client123:NEW_SERVER\n"
	if conn.String() != expected {
		t.Errorf("Expected response %q, got %q", expected, conn.String())
	}
}
