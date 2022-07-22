package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"time"
)

type Block struct {
	Timestamp    time.Time
	Hash         string
	PreviousHash string
	Data         string
	Nonce        int
}

type Chain struct {
	blocks []Block
	diff   int
}

func NewChain(diff int) *Chain {
	return &Chain{
		blocks: []Block{newGenesisBlock(diff)},
		diff:   diff,
	}
}

func (c *Chain) GetBlocks() []Block {
	return c.blocks
}

func (c *Chain) SetDifficulty(d int) {
	c.diff = d
}

func (c *Chain) AddBlock(data string) {
	b := Block{
		Timestamp:    time.Now(),
		Hash:         "123456789ABCDEF",
		PreviousHash: c.blocks[len(c.blocks)-1].Hash,
		Data:         data,
		Nonce:        0,
	}

	b.proofOfWork(c.diff)
	c.blocks = append(c.blocks, b)
}

func newGenesisBlock(diff int) Block {
	b := Block{
		Timestamp:    time.Now(),
		Hash:         "123456789ABCDEF",
		PreviousHash: "",
		Data:         "Let there be light",
		Nonce:        0,
	}

	b.proofOfWork(diff)
	return b
}

func (b *Block) proofOfWork(diff int) {
	for b.Hash[:diff] != strings.Repeat("0", diff) {
		b.Nonce++
		b.calculateHash()
	}
}

func (b *Block) calculateHash() {
	hashInput := fmt.Sprintf("%s|%s|%s|%d", b.Timestamp, b.PreviousHash, b.Data, b.Nonce)

	h := sha256.New()
	h.Write([]byte(hashInput))
	hashed := h.Sum(nil)
	b.Hash = hex.EncodeToString(hashed)
}
