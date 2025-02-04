package handlers

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

func TestHandleConnection(t *testing.T) {
	// Создаем mock-соединение
	conn := &MockConn{}
	conn.WriteString("client123:START_SERVER:map1:v1.0\n")

	// Создаем mock-конфигурацию
	cfg := &config.Config{
		Proxy: struct {
			Ip   string `yaml:"ip"`
			Port string `yaml:"port"`
			Code string `yaml:"code"`
		}{
			Ip:   "127.0.0.1",
			Port: "8080",
			Code: "admin123",
		},
	}

	HandleConnection(cfg, conn)

	// Проверяем результат
	expected := "client123:NEW_SERVER\n"
	if conn.String() != expected {
		t.Errorf("Expected response %q, got %q", expected, conn.String())
	}
}
