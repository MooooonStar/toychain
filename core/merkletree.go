package core

import (
	"crypto/sha256"
	"errors"
)

type MerkleTree struct {
	Root *Node
}

func NewMerkleTree(cs [][]byte) (*MerkleTree, error) {
	node, err := buildWithContent(cs)
	if err != nil {
		return nil, err
	}
	return &MerkleTree{node}, nil
}

type Node struct {
	Left    *Node
	Right   *Node
	leaf    bool
	dup     bool
	Hash    Hash
	Content []byte
}

func buildWithContent(contents [][]byte) (*Node, error) {
	if len(contents) == 0 {
		return nil, errors.New("no content found")
	}

	var leafs []*Node
	for _, content := range contents {
		node := &Node{
			leaf:    true,
			Content: content,
			Hash:    sha256.Sum256(content),
		}
		leafs = append(leafs, node)
	}

	if len(leafs)%2 == 1 {
		dup := new(Node)
		*dup = *(leafs[len(leafs)-1])
		dup.dup = true
		leafs = append(leafs, dup)
	}

	return buildIntermediate(leafs), nil
}

func buildIntermediate(nl []*Node) *Node {
	var nodes []*Node
	if len(nl) == 1 {
		return nl[0]
	}
	if len(nl)%2 == 1 {
		dup := new(Node)
		*dup = *nl[len(nl)-1]
		dup.dup = true
		nl = append(nl, dup)
	}

	for i := 0; i < len(nl)/2; i++ {
		left := nl[2*i]
		right := nl[2*i+1]
		data := append(left.Hash[:], right.Hash[:]...)
		node := &Node{
			Left:  left,
			Right: right,
			Hash:  sha256.Sum256(data),
		}
		nodes = append(nodes, node)
	}
	return buildIntermediate(nodes)
}
