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
	conn.Write([]byte(InputString () + END_BYTES))

	var (
		buffer = make ([]byte, BUFF_SIZE)
		message string
	)
	for {
		length, err := conn.Read(buffer)
		if ((length == 0) || (err != nil)) {break}
		message += string(buffer[:length])
	}
	fmt.Println(message)
}

func InputString () string {
	msg, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	return strings.Replace(msg, "\n", "", -1)
}