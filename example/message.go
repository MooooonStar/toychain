package main

import (
	"encoding/json"
)

const (
	MessageTypeTransaction = "TRANSACTION"
	MessageTypeBlock       = "BLOCK"
)

type Message struct {
	Type string
	Data json.RawMessage
}
