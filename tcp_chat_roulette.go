package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

var partner = make(chan io.ReadWriteCloser)

func match(c io.ReadWriteCloser) {
	fmt.Fprintln(c, "Waiting for a partner...")
	select {
	case partner <- c:
	case p := <-partner:
		chat(p, c)
	case <-time.After(5 * time.Second):
		fmt.Fprintln(c, "no one is available, close")
		c.Close()
	}
}

func chat(a, b io.ReadWriteCloser) {
	fmt.Fprintln(a, "Found one! Say hi.")
	fmt.Fprintln(b, "Found one! Say hi.")

	errc := make(chan error, 1)
	go cp(a, b, errc)
	go cp(b, a, errc)

	//if an error occurs, it will not block here and close the channels
	if err := <-errc; err != nil {
		log.Println(err)
	}
	a.Close()
	b.Close()
}

func cp(w io.Writer, r io.Reader, errc chan<- error) {
	wr, err := io.Copy(w, r)
	log.Println("wr:", wr)
	errc <- err
}

func main() {
	listener, err := net.Listen("tcp", "localhost:4000")
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go match(conn)
	}
}
