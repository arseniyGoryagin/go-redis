package main

import (
	"errors"
	"strconv"
	"strings"
	"sync"
)

type Handler struct {
	storageSet      map[string]string
	storageSetMutex sync.RWMutex
}

func NewHandler() Handler {
	return Handler{
		storageSet: make(map[string]string),
	}
}

func (h *Handler) Handle(command string, args []Value) (Value, error) {
	switch command {
	case "PING":
		return h.ping(args)
	case "GET":
		return h.get(args)
	case "SET":
		return h.set(args)
	case "HGET":
		return h.hget(args)
	case "HSET":
		return h.hset(args)
	default:
		return Value{}, errors.New("NO handler for command " + command)
	}
}

func (h *Handler) ping(args []Value) (Value, error) {

	if len(args) == 1 {
		return Value{typ: STRING, strValue: "PONG"}, nil
	}

	return Value{typ: STRING, strValue: h.parseValues(args)}, nil
}

func (h *Handler) parseValues(args []Value) string {

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
				res.WriteString(h.parseValues(args[i].values))
			}
		}
		res.WriteString(" ")
	}

	return res.String()
}

func (h *Handler) set(args []Value) (Value, error) {

	if len(args) != 3 {
		return Value{}, errors.New("wrong number of arguments")
	}

	Set(args[1].strValue, args[2].strValue)

	return Value{typ: STRING, strValue: "OK"}, nil

}

func (h *Handler) get(args []Value) (Value, error) {

	if len(args) != 2 {
		return Value{}, errors.New("wrong number of arguments")
	}

	found, val := Get(args[1].strValue)

	if !found {
		return Value{typ: NULL}, nil
	}

	return Value{typ: BULK, strValue: val}, nil

}

func (h *Handler) hset(args []Value) (Value, error) {

	if len(args) != 4 {
		return Value{}, errors.New("wrong number of arguments")
	}

	Hset(args[1].strValue, args[2].strValue, args[3].strValue)

	return Value{typ: STRING, strValue: "OK"}, nil

}

func (h *Handler) hget(args []Value) (Value, error) {

	if len(args) != 3 {
		return Value{}, errors.New("wrong number of arguments")
	}

	found, val := Hget(args[1].strValue, args[2].strValue)

	if !found {
		return Value{typ: NULL}, nil
	}

	return Value{typ: BULK, strValue: val}, nil

}
