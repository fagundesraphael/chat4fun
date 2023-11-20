package main

import (
	"log"
	"net"
)

const Port = "8080"

func handleConnection(conn net.Conn) {
	defer conn.Close()
	message := []byte("Hello, World!\n")
	n, err := conn.Write(message)
	if err != nil {
		log.Printf("ERROR: could not write to %s: %s", conn.RemoteAddr(), err)
		return
	}
	if n < len(message) {
		log.Printf("The message was not fully written %d/%d\n", n, len(message))
		return
	}
}

func main() {
	ln, err := net.Listen("tcp", ":"+Port)
	if err != nil {
		log.Fatalf("ERROR: could not listen to epic port %s\n", Port, err)
	}
	log.Printf("Listening to TPC connections on port %s ...", Port)
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("Could not accept the connection", err)
		}
		log.Printf("Accepted connection from %s", conn.RemoteAddr())
		go handleConnection(conn)
	}
}
