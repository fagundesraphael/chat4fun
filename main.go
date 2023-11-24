package main

import (
	"log"
	"net"
)

const (
	Port     = "8080"
	SafeMode = true
)

func safeRemoteAddr(conn net.Conn) string {
	if SafeMode {
		return "[REDACTED]"
	} else {
		return conn.RemoteAddr().String()
	}
}

type MessageType int

const (
	ClientConnected MessageType = iota + 1
	DisconnectClient
	NewMessage
)

type Message struct {
	Type MessageType
	Conn net.Conn
	Text string
}

func server(messages chan Message) {
	conns := map[string]net.Conn{}
	for {
		msg := <-messages
		switch msg.Type {
		case ClientConnected:
			conns[msg.Conn.RemoteAddr().String()] = msg.Conn
		case DisconnectClient:
			delete(conns, msg.Conn.RemoteAddr().String())
		case NewMessage:
			for _, conn := range conns {
				_, err := conn.Write([]byte(msg.Text))
				if err != nil {
					// TODO: remove the client from the list
					log.Printf("Could not send data to %s: %s", safeRemoteAddr(conn), err)
				}
			}
		}
	}
}

func client(conn net.Conn, messages chan Message) {
	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			conn.Close()
			messages <- Message{
				Type: DisconnectClient,
				Conn: conn,
			}
			return
		}
		messages <- Message{
			Type: NewMessage,
			Text: string(buffer[:n]),
			Conn: conn,
		}
	}
}

func main() {
	ln, err := net.Listen("tcp", ":"+Port)
	if err != nil {
		log.Fatalf("ERROR: could not listen to epic port %s\n", Port, err)
	}
	log.Printf("Listening to TPC connections on port %s ...", Port)

	messages := make(chan Message)
	go server(messages)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("Could not accept the connection", err)
		}
		log.Printf("Accepted connection from %s", safeRemoteAddr(conn))

		messages <- Message{Type: ClientConnected, Conn: conn}

		go client(conn, messages)
	}
}
