package main

import (
	"fmt"
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

		resp := NewRespReader(conn)
		res, err := resp.Read()
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(res)

		conn.Write([]byte("+OK\r\n"))
	}

}
