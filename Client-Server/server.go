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
	defer conn.Close()
	var (
		buffer = make ([]byte, BUFF_SIZE)
		message string
	)
	for {
		length, err := conn.Read(buffer)
		if ((length == 0) || (err != nil)) {break}
		message += string (buffer[:length])
		if strings.HasSuffix(message, END_BYTES) {
			message = strings.TrimSuffix(message, END_BYTES)
			break
		}
	}
	fmt.Println ("i got this string: " + message)
	message = strings.ToUpper(message)
	fmt.Println ("i send this string: " + message)
	conn.Write([]byte(message))
}