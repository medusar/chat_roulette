package tcp

import (
	"fmt"
	"math/rand"
	"net"
	"strings"
	"time"
)

// names from http://listofrandomnames.com/
var names = []string{"Marion Propp", "Miles Rene", "Renay Spiess", "Yer Burton", "Eliseo Brautigam", "Marlyn Miga",
	"Karren Waldorf", "Ciera Just", "Regena Haskell", "Gabriela Viviani", "Garfield Mike", "Pandora Fenimore",
	"Earle Haberle", "Florrie Sellars", "Rosanna Connor", "Anisha Kile", "Tiesha Shelley", "Oda Gilchrest",
	"Rod Guevara", "Karry Firestone",}

// generate a random name based on the address
func RandomName(c net.Conn) string {
	rand.Seed(time.Now().Unix())
	return names[rand.Intn(len(names))] + "@" + strings.Split(c.RemoteAddr().String(), ":")[0]
}

// User with a random name
type User struct {
	Con  net.Conn
	Name string
}

// Wrap a Con into a User
func New(c net.Conn) *User {
	return &User{Con: c, Name: RandomName(c)}
}

//implements io.Reader
func (u *User) Read(p []byte) (n int, err error) {
	n, err = u.Con.Read(p)
	if err != nil {
		return
	}
	// Each time read some data, add userName before the msg
	// It is import to use `p[:n]`, if use `string(p)`, it would
	// contain all the blank data into the string, which is too
	// large and will soon exceed the limit
	msg := fmt.Sprintln("[" + u.Name + "]: " + string(p[:n]))
	// Reread after u.Con.Read, will override the byte slice
	// Is there any better way?
	n, err = strings.NewReader(msg).Read(p)
	return
}

//implements io.Writer
func (u *User) Write(p []byte) (n int, err error) {
	return u.Con.Write(p)
}

//implements io.Closer
func (u *User) Close() error {
	return u.Con.Close()
}

//Write msg to the Con, with a new line character appended
func (u *User) WriteMsg(msg string) (int, error) {
	return u.Con.Write([]byte(fmt.Sprintln(msg)))
}
