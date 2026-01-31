package server

import (
	"bedis/internal/handler"
	"bedis/internal/storage"
	"bufio"
	"net"
	"strings"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
)

type Client struct {
	conn   net.Conn
	reader *bufio.Reader
}

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
	for range 50 {
		if server.listener != nil {
			break
		}
		time.Sleep(50 * time.Millisecond)
	}
	addr = server.listener.Addr().String()
	t.Log("server created on", addr)
	return server, addr
}

func request(t *testing.T, command string, client *Client) (string, error) {
	_, err := client.conn.Write([]byte(command))
	if err != nil {
		t.Logf("write failed: %v", err)
		return "", err
	}
	err = client.conn.SetReadDeadline(time.Now().Add(time.Second))
	if err != nil {
		t.Logf("error to set timeout on read: %v", err)
		return "", err
	}
	resp, err := client.reader.ReadString('\n')
	if err != nil {
		t.Logf("read failed: %v", err)
		return "", err
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
	client := Client{
		conn:   conn,
		reader: bufio.NewReader(conn),
	}
	command := "SET " + gofakeit.Noun() + " " + gofakeit.Adverb() + "\n"
	resp, err := request(t, command, &client)
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
	client := Client{
		conn:   conn,
		reader: bufio.NewReader(conn),
	}
	for i := 0; i < 100; i++ {
		key := gofakeit.Noun()
		value := gofakeit.Adverb()
		command := "SET " + key + " " + value + "\n"
		_, err := request(t, command, &client)
		if err != nil {
			t.Fatalf("error while SET request: %v", err)
		}
		command = "GET " + key + "\n"
		res, err := request(t, command, &client)
		if err != nil {
			t.Fatalf("error while GET request: %v", err)
		}
		if strings.TrimSpace(res) != value {
			t.Fatalf("SET and GET values dont match: %v, %s", res, value)
		}
	}
}

func TestUnknownCommand(t *testing.T) {
	server, addr := startTestServer(t)
	defer server.Stop()

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		t.Fatalf("dial failed: %s", err)
	}

	defer func() {
		err := conn.Close()
		if err != nil {
			t.Errorf("conn close failed: %s", err)
		}
	}()
	client := Client{
		conn:   conn,
		reader: bufio.NewReader(conn),
	}
	command := gofakeit.Noun() + "\n"
	t.Log(command)
	resp, err := request(t, command, &client)
	if err != nil {
		t.Fatalf("request failed: %s", err)
	}
	if strings.TrimSpace(resp) != "unknown command" {
		t.Fatalf("unknown command wrong response: %s", resp)
	}
}
