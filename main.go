package main

import (
	"encoding/binary"
	"log"
)

type (
	/*
	 * Node structure:
	 * | type | nkeys |  pointers  |   offsets  | key-values
	 * |  2B  |   2B  | nkeys * 8B | nkeys * 2B | ...
	 */
	BNode struct {
		data []byte
	}

	BTree struct {
		// pointer to disk page
		root uint64

		get func(uint64) BNode // dereference pointer
		new func(BNode) uint64 // allocation
		del func(uint64)       // deallocation
	}
)

const (
	BNODE_NODE = 1
	BNODE_LEAF = 2

	// size constraints (page is limited to 4K bytes)
	HEADER             = 4
	BTREE_PAGE_SIZE    = 4096
	BTREE_MAX_KEY_SIZE = 1000
	BTREE_MAX_VAL_SIZE = 3000
)

func init() {
	node1max := HEADER + 8 + 2 + 4 + BTREE_MAX_KEY_SIZE + BTREE_MAX_VAL_SIZE
	if node1max > BTREE_PAGE_SIZE {
		log.Fatalf("Node exceeds page size: %d Bytes", node1max)
	}
}

func (node *BNode) btype() uint16 {
	return binary.LittleEndian.Uint16(node.data)
}

func (node *BNode) nkeys() uint16 {
	return binary.LittleEndian.Uint16(node.data[2:4])
}

func (node *BNode) setHeader(btype uint16, nkeys uint16) {
	binary.LittleEndian.PutUint16(node.data[0:2], btype)
	binary.LittleEndian.PutUint16(node.data[2:4], nkeys)
}

func (node *BNode) getPtr(idx uint16) uint64 {
	if idx >= node.nkeys() {
		log.Fatalf("Index out of key bound (%d): got %d", node.nkeys(), idx)
	}
	pos := HEADER + 8*idx
	return binary.LittleEndian.Uint64(node.data[pos:])
}

func (node *BNode) setPtr(idx uint16, val uint64) {
	if idx >= node.nkeys() {
		log.Fatalf("Index out of key bound (%d): got %d", node.nkeys(), idx)
	}
	pos := HEADER + 8*idx
	binary.LittleEndian.PutUint64(node.data[pos:], val)
}
