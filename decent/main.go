package main

import (
	"os"
	"fmt"
	"net"
	"bufio"
	"strings"
	"encoding/json"
)

const (
	BUFF_SIZE = 512
)

type addr_t struct {
	IPv4 string
	Port string
}

type node_t struct {
	Connections map[string]bool
	Address addr_t
}

type pack_t struct {
	To string
	From string
	Data string
}

func init () {
	if len (os.Args) != 2 {panic ("len args != 2")}
}

func main () {
	NewNode (os.Args[1]).Run (handleServer, handleClient)
}

func NewNode (address string) *node_t {
	splited := strings.Split(address, ":")
	if len (splited) != 2 {return nil}
	return &node_t {
		Connections: make (map[string]bool),
		Address: addr_t {
			IPv4: splited[0],
			Port: ":" + splited[1],
		},
	}
}

func (node *node_t) Run (handleServer func(*node_t), handleClient func(*node_t)) {
	go handleServer (node)
	handleClient (node)
}

func handleServer (node *node_t) {
	listen, err := net.Listen ("tcp", "0.0.0.0" + node.Address.Port)
	if err != nil {
		panic ("listen error")
	}
	defer listen.Close ()
	for {
		conn, err := listen.Accept ()
		if err != nil {break}
		go handleConnection (node, conn)
	}
}

func handleConnection (node *node_t, conn net.Conn) {
	defer conn.Close ()
	var (
		buffer = make ([]byte, BUFF_SIZE)
		message string
		pack pack_t
	)
	for {
		length, err := conn.Read (buffer)
		if err != nil {break}
		message += string (buffer[:length])
	}
	err := json.Unmarshal ([]byte (message), &pack)
	if err != nil {return}

	node.ConnectTo ([]string{pack.From})
	fmt.Println (pack.Data)
}

func handleClient (node *node_t) {
	for {
		message := InputString ()
		splited := strings.Split (message, " ")
		switch splited[0] {
			case "/exit": os.Exit (0)
			case "/connect": node.ConnectTo (splited[1:])
			case "/network": node.PrintNetwork ()
			default: node.SendToAll (message)
		}
	}	
}

func (node *node_t) PrintNetwork () {
	for addr := range node.Connections {
		fmt.Println ("|", addr)
	}
}

func (node *node_t) ConnectTo (address_arr []string) {
	for _, addr := range address_arr {
		node.Connections[addr] = true
	}
}

func (node *node_t) SendToAll (message string) {
	var new_pack = pack_t {
		From: node.Address.IPv4 + node.Address.Port,
		Data: message,
	}
	for addr := range node.Connections {
		new_pack.To = addr
		node.Send (&new_pack) 
	}
}

func (node *node_t) Send (pack *pack_t) {
	conn, err := net.Dial ("tcp", pack.To)
	if err != nil {
		delete (node.Connections, pack.To)
		return
	}
	defer conn.Close ()
	json_pack, _ := json.Marshal (*pack)
	conn.Write (json_pack)
}

func InputString () string {
	msg, _ := bufio.NewReader (os.Stdin).ReadString('\n')
	return strings.Replace (msg, "\n", "", -1)
}
