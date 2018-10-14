package main

import (
	"fmt"
	"io"
	"log"
	"net"
)

const listenAddr = "localhost:4000"

func echoServer() {
	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		//copy data from reader to writer
		//but will be blocked here, waiting for reading data from the 1st connection
		//other connections will not work
		io.Copy(conn, conn)
	}
}

func concurrentEchoServer() {
	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go io.Copy(conn, conn)
	}
}

func sayHiServer() {
	l, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintln(conn, "hello")
		conn.Close()
	}
}

func main() {
	//sayHiServer()
	//echoServer()
	concurrentEchoServer()
}
