package server

import (
	"bedis/internal/handler"
	"bedis/internal/storage"
	"bufio"
	"github.com/brianvoe/gofakeit/v7"
	"net"
	"strings"
	"testing"
	"time"
)

func startTestServer(t *testing.T) (server *Server, addr string) {
	t.Helper()
	s := storage.New(nil)
	h := handler.New(s, nil)
	server = New("localhost:0", h, nil)

	go func() {
		if err := server.Start(); err != nil {
			t.Log(err)
			return
		}
	}()
	for i := 0; i < 50; i++ {
		if server.listener != nil {
			break
		}
		time.Sleep(50 * time.Millisecond)
	}
	addr = server.listener.Addr().String()
	t.Log("server created on", addr)
	return server, addr
}

func request(t *testing.T, command string, conn net.Conn) (string, error) {
	_, err := conn.Write([]byte(command))
	if err != nil {
		t.Logf("write failed: %v", err)
		return "", err
	}
	err = conn.SetReadDeadline(time.Now().Add(time.Second))
	if err != nil {
		t.Logf("error to set timeout on read: %v", err)
		return "", err
	}
	resp, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		t.Logf("read failed: %v", err)
		return "", nil
	}
	return resp, nil

}

func TestServerSingleCommand(t *testing.T) {
	server, addr := startTestServer(t)
	defer server.Stop()

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		t.Fatalf("dial failed: %v", err)
	}
	defer func() {
		err := conn.Close()
		if err != nil {
			t.Errorf("error on client close: %v", err)
		}
	}()
	command := "SET " + gofakeit.Noun() + " " + gofakeit.Adverb() + "\n"
	resp, err := request(t, command, conn)
	if err != nil {
		t.Fatalf("error while request: %v", err)
	}

	if strings.TrimSpace(resp) != "OK" {
		t.Fatalf("wrong response: %v", resp)
	}
}

func TestServerManyCommandsPerConnection(t *testing.T) {
	server, addr := startTestServer(t)
	defer server.Stop()

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		t.Fatalf("dial failed: %v", err)
	}
	defer func() {
		err := conn.Close()
		if err != nil {
			t.Errorf("error on client close: %v", err)
		}
	}()
	for i := 0; i < 100; i++ {
		key := gofakeit.Noun()
		value := gofakeit.Adverb()
		command := "SET " + key + " " + value + "\n"
		_, err := request(t, command, conn)
		if err != nil {
			t.Fatalf("error while SET request: %v", err)
		}
		command = "GET " + key + "\n"
		res, err := request(t, command, conn)
		if err != nil {
			t.Fatalf("error while GET request: %v", err)
		}
		if strings.TrimSpace(res) != value {
			t.Fatalf("SET and GET values dont match: %v, %s", res, value)
		}
	}

}
