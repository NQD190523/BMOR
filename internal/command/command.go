package command

import (
	"Build_my_own_redis/internal/protocol"
	"strings"
	"sync"
)

var store = struct {
	sync.Mutex
	data map[string]string
}{
	data: make(map[string]string),
}

func Handle(args []string) []byte {
	if len(args) == 0 {
		//return protocol is empty
		return protocol.EncodeNil()
	}

	switch strings.ToUpper(args[0]) {
	case "PING":
		return protocol.EncodeSingleString("PONG")
	case "SET":
		if len(args) < 3 {
			return protocol.EncodeErrorString("Error: SET command requires a key and a value")
		}
		key := args[1]
		value := args[2]
		store.Lock()
		store.data[key] = value
		store.Unlock()
		return protocol.EncodeSingleString("OK")
	case "GET":
		if len(args) < 2 {
			return protocol.EncodeErrorString("Error: GET command requires a key")
		}
		key := args[1]
		store.Lock()
		value := store.data[key]
		store.Unlock()
		return protocol.EncodeSingleString(value)
	default:
		return protocol.EncodeErrorString("Error: Unknown command")
	}
}
