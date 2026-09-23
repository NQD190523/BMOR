package protocol

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

func EncodeNil() []byte {
	return []byte("$-1\r\n")
}

func EncodeSingleString(s string) []byte {
	return []byte("+" + s + "\r\n")
}

func EncodeBulkString(s string) []byte {
	return []byte("$" + string(rune(len(s))) + "\r\n" + s + "\r\n")
}

func EncodeErrorString(s string) []byte {
	return []byte("-" + s + "\r\n")
}

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

func ReadBulkString(r *bufio.Reader, line string) ([]string, error) {
	lenStr, err := strconv.Atoi(line[1:])
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
