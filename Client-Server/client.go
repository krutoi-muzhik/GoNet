package main

import (
	"os"
	"fmt"
	"net"
	"bufio"
	"strings"
)

const (
	BUFF_SIZE = 512
	END_BYTES = "\000BEEBA\000"
	ADDR_SERVER = ":8080"
)

func main () {
	conn, err := net.Dial("tcp", ADDR_SERVER)
	if (err != nil) {
		panic ("cant connect to server")
	}
	defer conn.Close()
	go ClientOutput(conn)
	ClientInput(conn)
}

func ClientInput (conn net.Conn) {
	for {
		conn.Write([]byte(InputString() + END_BYTES))
	}
}

func ClientOutput (conn net.Conn) {
	var (
		buffer = make([]byte, 512)
		message string
	)
	close: for {
		message = ""
		for {
			length, err := conn.Read(buffer)
			if err != nil {break close}
			message += string(buffer[:length])
			if strings.HasSuffix(message, END_BYTES) {
				message = strings.TrimSuffix(message, END_BYTES)
				break
			}
		}
		fmt.Println(message)
	}
}

func InputString () string {
	msg, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	return strings.Replace(msg, "\n", "", -1)
}