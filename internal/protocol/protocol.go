package protocol

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// EncodeNil encodes a nil value in the Redis protocol format.
func EncodeNil() []byte {
	return []byte("$-1\r\n")
}

// EncodeSingleString encodes a single string in the Redis protocol format.
func EncodeSingleString(s string) []byte {
	return []byte("+" + s + "\r\n")
}

// EncodeBulkString encodes a bulk string in the Redis protocol format.
func EncodeBulkString(s string) []byte {
	return []byte("$" + string(rune(len(s))) + "\r\n" + s + "\r\n")
}

// EncodeErrorString encodes an error string in the Redis protocol format.
func EncodeErrorString(s string) []byte {
	return []byte("-" + s + "\r\n")
}

// ReadCommand reads a command from the reader and returns the arguments.
func ReadCommand(r *bufio.Reader) ([]string, error) {
	line, err := r.ReadString('\n')
	if err != nil {
		return nil, err
	}
	if len(line) == 0 {
		return nil, nil
	}
	if line[0] != '*' {
		return strings.Fields(line), nil
	}
	return ReadBulkString(r, line)
}

// ReadBulkString reads a bulk string from the reader and returns the arguments.
func ReadBulkString(r *bufio.Reader, line string) ([]string, error) {
	lenStr, err := strconv.Atoi(strings.TrimSpace(line[1:]))
	if err != nil {
		return nil, err
	}
	// Read the bulk string data
	args := make([]string, lenStr) // +2 for \r\n
	for i := 0; i < lenStr; i++ {
		typeLine, err := r.ReadString('\n')
		if err != nil {
			return nil, err
		}
		typeLine = strings.TrimRight(typeLine, "\r\n")
		if len(typeLine) == 0 || typeLine[0] != '$' {
			return nil, fmt.Errorf("protocol error: expected '$', got %q", typeLine)
		}
		sizeLine, err := strconv.Atoi(typeLine[1:])
		if err != nil {
			return nil, err
		}
		buf := make([]byte, sizeLine+2) // +2 for \r\n
		_, err = io.ReadFull(r, buf)
		if err != nil {
			return nil, err
		}
		args[i] = string(buf[:sizeLine])
	}
	return args, nil
}
