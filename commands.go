package main

import (
	"errors"
	"strings"
)

func ParseCommand(value Value) (string, error) {

	if value.typ != ARRAY || len(value.values) == 0 {
		return "", errors.New("incorrect_command")
	}

	return strings.ToUpper(value.values[0].strValue), nil
}
