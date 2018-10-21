package tcp

import (
	"io"
	"log"
	"net"
)

var partnerChan = make(chan *User)

func start(address string) {
	listener, e := net.Listen("tcp", address)
	if e != nil {
		log.Fatal("error starting tcp server", e)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("error accept", err)
			continue
		}
		go serve(conn)
	}
}

func serve(conn net.Conn) {
	user := New(conn)
	user.WriteMsg("Waiting for a partner")
	select {
	case partnerChan <- user:
	case partner := <-partnerChan:
		privateChat(partner, user)
	}
}

func privateChat(u1, u2 *User) {
	u1.WriteMsg(u2.Name + " has come, try to say hi")
	u2.WriteMsg(u1.Name + " has come, try to say hi")

	errChan := make(chan error)

	go cp(u1, u2, errChan)
	go cp(u2, u1, errChan)

	//if no error occurs, it will block here
	for e := <-errChan; e != nil; {
		log.Println("error chat", e)
	}

	u1.Close()
	u2.Close()
}

func cp(p1 io.ReadWriteCloser, p2 io.ReadWriteCloser, errChan chan error) {
	//io.Copy will block until an error occurs
	_, err := io.Copy(p1, p2)
	errChan <- err
}

func StarTcpRoulette() {
	start("localhost:4000")
}
