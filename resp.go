package main

import (
	"bufio"
	"errors"
	"io"
	"strconv"
	"strings"
)

const (
	STRING  = "+"
	ERROR   = "-"
	INTEGER = ":"
	BULK    = "$"
	ARRAY   = "*"
)

type Token struct {
	typ         string
	strValue    string
	intValue    int
	tokenValues []Token
}
type header struct {
	typ  string
	size int
}

type RespReader struct {
	reader *bufio.Reader
}

func NewRespReader(rd io.Reader) *RespReader {
	return &RespReader{reader: bufio.NewReader(rd)}
}
func (r *RespReader) Read() (Token, error) {
	return r.read("", 0)
}

func (r *RespReader) read(typ string, size int) (token Token, err error) {

	h := header{typ: typ, size: size}
	if typ == "" {
		h, err = r.readHeader()
		if err != nil {
			return Token{}, err
		}
	}

	token.typ = h.typ

	switch h.typ {
	case STRING:
		token.strValue, err = r.readString()
		if err != nil {
			return Token{}, err
		}

		return token, nil
	case ERROR:
		token.strValue, err = r.readString()
		if err != nil {
			return Token{}, err
		}

		return token, nil
	case INTEGER:
		token.intValue, err = r.readInt()
		if err != nil {
			return Token{}, err
		}

		return token, nil
	case BULK:

		if h.size == -1 {
			return token, nil
		}

		token.strValue, err = r.readBulk(h.size)
		if err != nil {
			return Token{}, err
		}

		return token, nil
	case ARRAY:

		if h.size == -1 {
			return token, nil
		}

		tokens := make([]Token, h.size)
		for i := 0; i < h.size; i++ {

			t, err := r.read("", 0)
			if err != nil {
				return Token{}, err
			}

			tokens[i] = t
		}

		token.tokenValues = tokens

		return token, nil
	default:
		return Token{}, errors.New("unknown_value")
	}
}

func (r *RespReader) readHeader() (header, error) {

	h := header{}

	typ, err := r.readByte()
	if err != nil {
		return header{}, err
	}

	h.typ = typ

	if h.typ == "*" || h.typ == "$" {

		var intStrVal strings.Builder

		for {
			val, err := r.readByte()
			if err != nil {
				return header{}, err
			}

			if val == "\r" {
				break
			}

			intStrVal.WriteString(val)
		}

		size, err := strconv.Atoi(string(intStrVal.String()))
		if err != nil {
			return header{}, err
		}

		h.size = size

		_, err = r.readByte()
		if err != nil {
			return header{}, err
		}

	}

	return h, nil
}

func (r *RespReader) readByte() (string, error) {
	buf := make([]byte, 1)
	_, err := io.ReadFull(r.reader, buf)
	if err != nil {
		return "", err
	}
	return string(buf[0]), nil
}

func (r *RespReader) readInt() (int, error) {
	numStr, err := r.readString()
	if err != nil {
		return 0, err
	}
	i, err := strconv.Atoi(numStr)
	if err != nil {
		return 0, err
	}
	return i, nil
}

func (r *RespReader) readString() (string, error) {
	res, err := r.reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSuffix(res, "\r\n"), nil
}

func (r *RespReader) readBulk(size int) (string, error) {

	buf := make([]byte, size+2)
	_, err := io.ReadFull(r.reader, buf)
	if err != nil {
		return "", err
	}

	return string(buf[:size]), nil
}
