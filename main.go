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

		writer := NewWriter(conn)

		reader := NewReader(conn)

		res, err := reader.Read()
		if err != nil {
			err := writer.Write(Value{typ: ERROR, strValue: ERROR + err.Error()})
			if err != nil {
				fmt.Println(err)
			}
			continue
		}

		command, err := ParseCommand(res)
		if err != nil {
			err := writer.Write(Value{typ: ERROR, strValue: ERROR + err.Error()})
			if err != nil {
				fmt.Println(err)
			}
			continue
		}

		handler, ok := Handlers[command]
		if !ok {
			err := writer.Write(Value{typ: ERROR, strValue: ERROR + "no handler for " + command})
			if err != nil {
				fmt.Println(err)
			}
			continue
		}

		resp, err := handler(res.values)
		if err != nil {
			err := writer.Write(Value{typ: ERROR, strValue: ERROR + err.Error()})
			if err != nil {
				fmt.Println(err)
			}
			continue
		}

		err = writer.Write(resp)
		if err != nil {
			fmt.Println(err)
			continue
		}
	}

}
