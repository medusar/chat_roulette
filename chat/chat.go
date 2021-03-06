package chat

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/longda/markov"
	"html/template"
	"io"
	"log"
	"net"
	"net/http"
	"time"
)

const LISTEN_ADDR = "localhost:4000"

func netListen() {
	l, err := net.Listen("tcp", "localhost:4001")
	if err != nil {
		log.Fatal(err)
	}

	for {
		c, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go match(c)
	}
}

var templates = template.Must(template.ParseFiles("chat/ws.html"))

func rootHandler(w http.ResponseWriter, r *http.Request) {
	err := templates.ExecuteTemplate(w, "ws.html", LISTEN_ADDR)
	if err != nil {
		log.Println(err)
	}
}

type wsCon struct {
	ws *websocket.Conn
}

func (s *wsCon) Read(b []byte) (int, error) {
	_, msg, err := s.ws.ReadMessage()
	if err != nil {
		return 0, err
	}
	copy(b, msg)
	return len(b), nil
}

func (s *wsCon) Write(b []byte) (int, error) {
	err := s.ws.WriteMessage(websocket.TextMessage, b)
	if err != nil {
		return 0, err
	}
	return len(b), nil
}

func (s *wsCon) Close() error {
	return s.ws.Close()
}

type socket struct {
	r    io.Reader
	w    io.Writer
	done chan bool
}

func (s socket) Read(b []byte) (int, error) {
	return s.r.Read(b)
}

func (s socket) Write(b []byte) (int, error) {
	return s.w.Write(b)
}

func (s socket) Close() error {
	s.done <- true
	return nil
}

var chain = markov.NewChain(2) // 2-word prefixes

var upgrader = websocket.Upgrader{} // use default options

func socketHandler(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}

	ws := &wsCon{ws: c}

	//!!! how to use Pipe, when and how?
	rd, wr := io.Pipe()
	go func() {
		_, err := io.Copy(io.MultiWriter(wr, chain), ws)
		wr.CloseWithError(err)
	}()

	s := &socket{r: rd, w: ws, done: make(chan bool)}
	go match(s)
	<-s.done
}

var partner = make(chan io.ReadWriteCloser)

func match(c io.ReadWriteCloser) {
	fmt.Fprint(c, "Waiting for a partner...")

	select {
	case partner <- c:
		// now handled by the other goroutine
	case p := <-partner:
		chat(p, c)
	case <-time.After(5 * time.Second):
		chat(Bot(), c)
	}
}

func chat(a, b io.ReadWriteCloser) {
	fmt.Fprintln(a, "Found one! Say hi.")
	fmt.Fprintln(b, "Found one! Say hi.")

	errc := make(chan error, 1)

	go cp(a, b, errc)
	go cp(b, a, errc)

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

// Bot returns an io.ReadWriteCloser that responds to each incoming write with a generated sentence.
func Bot() io.ReadWriteCloser {
	r, out := io.Pipe() // for outgoing data
	return bot{r, out}
}

type bot struct {
	io.ReadCloser
	out io.Writer
}

func (b bot) Write(buf []byte) (int, error) {
	go b.speak()
	return len(buf), nil
}

func (b bot) speak() {
	time.Sleep(time.Second)
	msg := chain.Generate(10) // at most 10 words
	b.out.Write([]byte(msg))
}

func Start() {
	go netListen()

	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/socket", socketHandler)

	err := http.ListenAndServe(LISTEN_ADDR, nil)

	if err != nil {
		log.Fatal(err)
	}
}
