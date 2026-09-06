package main

import (
	"strconv"
	"strings"
)

func ping(args []Value) (Value, error) {

	if len(args) == 1 {
		return Value{typ: STRING, strValue: "PONG"}, nil
	}

	return Value{typ: STRING, strValue: parseValues(args)}, nil
}

func parseValues(args []Value) string {

	var res strings.Builder

	for i := range len(args) {
		switch args[i].typ {
		case STRING:
			res.WriteString(args[i].strValue)
		case ERROR:
			res.WriteString(args[i].strValue)
		case INTEGER:
			res.WriteString(strconv.Itoa(args[i].intValue))
		case BULK:
			res.WriteString(args[i].strValue)
		case ARRAY:
			for range len(args[i].values) {
				res.WriteString(parseValues(args[i].values))
			}
		}
		res.WriteString(" ")
	}

	return res.String()
}

var Handlers = map[string]func([]Value) (Value, error){
	"PING": ping,
}
