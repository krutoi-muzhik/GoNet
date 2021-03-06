package main

import (
	"fmt"
	"net"
	"strings"
)

const (
	END_BYTES = "\000BEEBA\000"
	BUFF_SIZE = 512
	PORT = ":8080"
)

var (
	Connections = make(map[net.Conn]bool)
)

func main () {
	listen, err := net.Listen("tcp", PORT)
	if (err != nil) {
		panic ("server error")
	}
	defer listen.Close()
	for {
		conn, err := listen.Accept()
		if (err != nil) {break}
		go handleConnect (conn)
	}
}

func handleConnect (conn net.Conn) { 
	Connections[conn] = true
	var (
		buffer = make ([]byte, BUFF_SIZE)
		message string
	)
	close: for {
		message = ""
		for {
			length, err := conn.Read(buffer)
			if (err != nil) {break close}
			message += string (buffer[:length])
			if strings.HasSuffix(message, END_BYTES) {
				message = strings.TrimSuffix(message, END_BYTES)
				break
			}
		}
		fmt.Println ("i got this string: " + message)
		for addr := range Connections {
			if (addr == conn) {continue}
			addr.Write([]byte(strings.ToUpper(message) + END_BYTES))
		}
	}

	delete (Connections, conn)
}