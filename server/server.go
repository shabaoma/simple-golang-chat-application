package main

import (
	"Simple-Chat-Application/message"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"strconv"
	"sync"
)

type Connection struct {
	Id      string
	Conn    net.Conn
	ConnMap *sync.Map
}

func (c Connection) Process() {
	defer c.Conn.Close()

	for {
		buf := make([]byte, 1024)

		// Wait for the client to send data.
		fmt.Printf("Wait for client: %v\n", c.Conn.RemoteAddr().String())
		n, err := c.Conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				fmt.Println("EOF")
			} else {
				fmt.Println("Read message error:", err)
			}
			break
		}

		m := message.Message{}
		json.Unmarshal(buf[:n], &m)
		m.From = c.Id
		fmt.Printf("Read message %v\n", m)
		b, err := json.Marshal(m)
		if err != nil {
			fmt.Println("Marshal error:", err)
			continue
		}

		if redirect_conn, ok := c.ConnMap.Load(m.To); ok {
			redirect_conn.(net.Conn).Write(b)
		}
	}
}

func main() {
	var connectionMap sync.Map
	id := 0

	go startWebServer()

	server, err := net.Listen("tcp", ":8082")
	if err != nil {
		fmt.Println("Start server error: ", err)
		return
	}
	fmt.Println("TCP server listening on :8082")
	defer server.Close()

	for {
		fmt.Println("Wait for client to connect")
		conn, err := server.Accept()
		if err != nil {
			fmt.Println("Accept error: ", err)
			continue
		}

		fmt.Printf("Assign id %v to %v\n", id, conn.RemoteAddr().String())
		str_id := strconv.Itoa(id)
		connectionMap.Store(str_id, conn)
		connection := Connection{str_id, conn, &connectionMap}
		id++
		go connection.Process()
	}
}
