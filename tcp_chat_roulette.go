package main

import (
	"fmt"
	"io"
	"log"
	"net"
)

var partner = make(chan io.ReadWriteCloser)

func match(c io.ReadWriteCloser) {
	fmt.Fprint(c, "Waiting for a partner...")
	select {
	case partner <- c:
		// now handled by the other goroutine
	case p := <-partner:
		chat(p, c)
	}
}

//func chat(a, b io.ReadWriteCloser) {
//	fmt.Fprintln(a, "Found one! Say hi.")
//	fmt.Fprintln(b, "Found one! Say hi.")
//	go io.Copy(a, b) //async send data from b to a
//	io.Copy(b, a) //sync send data from a to b
//}

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
	_, err := io.Copy(w, r)
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
