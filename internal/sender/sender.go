package sender

import (
	"fmt"
	"net"
)

type Sender struct {
	addr string
	conn *net.TCPConn
}

func New(addr string) *Sender {
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		fmt.Println("Dial failed:", err.Error())
		panic(err)

	}
	return &Sender{addr: addr, conn: conn}
}

func (s *Sender) Send(msg string) {
	_, err := s.conn.Write([]byte(msg))
	if err != nil {
		fmt.Println("Write failed:", err.Error())
	}
	reply := make([]byte, 1024)
	_, err = s.conn.Read(reply)
	if err != nil {
		fmt.Println("Read failed:", err.Error())

	}

	fmt.Println(string(reply))
	err = s.conn.Close()
	if err != nil {
		return
	}
}
