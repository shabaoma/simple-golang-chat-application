package main

import (
	"Simple-Chat-Application/message"
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

func GetParams(s string) (map[string]string, error) {
	// -t id of receiver
	// -m message to send
	ret := make(map[string]string)
	if strings.Contains(s, "-t ") && strings.Contains(s, "-m ") {
		for i, c := range s {
			if c == '-' {
				stack := []rune{}
				for _, c2 := range s[i+3:] {
					if c2 == '-' {
						break
					}
					stack = append(stack, c2)
				}
				ret[s[i:i+2]] = string(stack[:len(stack)-1])
			}
		}
		return ret, nil
	}
	return ret, errors.New("Input format not vilid. Please use -t [id] -m message.")
}

func Process(conn net.Conn) {
	defer conn.Close()

	for {
		buf := make([]byte, 1024)

		// Wait for the server to send data.
		n, err := conn.Read(buf)
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
		fmt.Printf("From: %v\nMessage: %v\n", m.From, m.Text)
	}
}

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8082")
	if err != nil {
		fmt.Println("Client dial error: ", err)
	}
	go Process(conn)

	// Read input from the keyboard.
	reader := bufio.NewReader(os.Stdin)
	for {
		line, err := reader.ReadString('\n') // args []byte not string
		if err != nil {
			fmt.Println("Read line error: ", err)
			continue
		}
		if line == "exit\n" {
			break
		}

		params, err := GetParams(line)
		if err != nil {
			fmt.Println("Get params error: ", err)
			continue
		}

		m := message.New("", params["-t"], params["-m"])
		b, err := json.Marshal(*m)
		if err != nil {
			fmt.Println("Marshal error:", err)
			continue
		}

		_, err = conn.Write(b)
		if err != nil {
			fmt.Println("Send message error:", err)
		}
		// fmt.Printf("Send %v byte to server\n", n)
	}
}
