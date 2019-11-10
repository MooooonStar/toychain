package main

import (
	"encoding/json"

	"github.com/MooooonStar/minichain/core"
)

type Node struct {
	bc *core.Blockchain
}

func NewNode() *Node {
	return &Node{
		bc: core.NewBlockchain("17XLQvEM5uKPvuFPqfN8op2GQ6zs51Rqjv", "Long Live The Bitcoin."),
	}
}

func (node *Node) Run(port int, dest string, debug bool) error {
	in, _, err := StartP2P(port, dest, debug)
	if err != nil {
		return err
	}
	for {
		select {
		case data := <-in:
			var msg Message
			if err := json.Unmarshal(data, &msg); err != nil {
				continue
			}
			switch msg.Type {
			case MessageTypeTransaction:
				var tx core.Transaction
				if err := json.Unmarshal(msg.Data, &tx); err != nil {
					return err
				}
			case MessageTypeBlock:
				var block core.Block
				if err := json.Unmarshal(msg.Data, &block); err != nil {
					return err
				}
				if core.CheckProofOfWork(block) {
					node.bc.AddBlock(&block)
				}
			default:
			}
		}
	}
}
