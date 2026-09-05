package main

import (
	"fmt"
	"io"
	"log"
	"net"
)

func main() {
	fmt.Println("Listening on port 6379")

	l, err := net.Listen("tcp", ":6379")
	if err != nil {
		log.Fatal(err)
	}

	conn, err := l.Accept()
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	for {
		buf := make([]byte, 1024)

		_, err := conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}

		conn.Write([]byte("+OK\r\n"))
	}

}
