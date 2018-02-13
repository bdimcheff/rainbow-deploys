package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"time"
)

func getColor() string {
	return os.Getenv("COLOR")
}

func tcpHandler(c net.Conn) {
	for {
		c.Write([]byte(fmt.Sprintf("The color is #%s", getColor())))
		c.Write([]byte(fmt.Sprintln()))
		time.Sleep(5 * time.Second)

	}
}

func serveTCP() {
	ln, err := net.Listen("tcp", ":8081")
	if err != nil {
		// handle error but not today
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			// handle error but not today
		}
		go tcpHandler(conn)
	}
}

func httpHandler(w http.ResponseWriter, r *http.Request) {
	color := getColor()

	fmt.Printf("Serving color: #%s", color)
	fmt.Println()

	fmt.Fprintf(w, "<body bgcolor=\"#%s\"><h1>#%s</h1></body>", color, color)
}

func main() {
	color := getColor()

	fmt.Printf("Booted with color: #%s", color)
	fmt.Println()

	go serveTCP()

	http.HandleFunc("/", httpHandler)
	fmt.Println("listening with http on :8080 and tcp on :8081")
	http.ListenAndServe(":8080", nil)
}
