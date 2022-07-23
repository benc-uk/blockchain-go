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
	Index        int
}

type Chain struct {
	blocks []Block
	diff   int
}

// NewChain creates a new chain with the given difficulty
func NewChain(diff int) *Chain {
	return &Chain{
		blocks: []Block{newGenesisBlock(diff)},
		diff:   diff,
	}
}

// GetBlocks returns the blocks of the chain
func (c *Chain) GetBlocks() []Block {
	return c.blocks
}

// SetDifficulty sets the difficulty of the chain
func (c *Chain) SetDifficulty(d int) {
	c.diff = d
}

// AddBlock adds a block to the chain with the given data
func (c *Chain) AddBlock(data string) string {
	b := Block{
		Timestamp:    time.Now(),
		Hash:         "123456789ABCDEF",
		PreviousHash: c.blocks[len(c.blocks)-1].Hash,
		Data:         data,
		Nonce:        0,
		Index:        len(c.blocks),
	}

	b.proofOfWork(c.diff)
	c.blocks = append(c.blocks, b)

	return b.Hash
}

// newGenesisBlock creates the first block in the chain
func newGenesisBlock(diff int) Block {
	b := Block{
		Timestamp:    time.Now(),
		Hash:         "123456789ABCDEF",
		PreviousHash: "",
		Data:         "Let there be light",
		Nonce:        0,
		Index:        0,
	}

	b.proofOfWork(diff)

	return b
}

// proofOfWork loops until the block hash starts with the given difficulty
func (b *Block) proofOfWork(diff int) {
	for b.Hash[:diff] != strings.Repeat("0", diff) {
		b.Nonce++
		b.Hash = b.calculateHash()
	}
}

// calculateHash calculates the hash of the block, note that the hash is not stored in the block
func (b *Block) calculateHash() string {
	hashInput := fmt.Sprintf("%s|%s|%s|%d", b.Timestamp, b.PreviousHash, b.Data, b.Nonce)

	h := sha256.New()
	h.Write([]byte(hashInput))
	hashed := h.Sum(nil)

	return hex.EncodeToString(hashed)
}

// For debugging blocks
func (b *Block) String() string {
	return fmt.Sprintf("Timestamp: %s\nPreviousHash: %s\nData: %s\nNonce: %d\nHash: %s\n",
		b.Timestamp, b.PreviousHash, b.Data, b.Nonce, b.Hash)
}

// UpdateBlock updates the data of the block
func (c *Chain) UpdateBlock(b Block, newData string) {
	//b := c.FindBlock(hash)
	c.blocks[b.Index].Data = newData
}

// FindBlock finds the block with the given hash
func (c *Chain) FindBlock(hash string) *Block {
	for _, block := range c.blocks {
		if block.Hash == hash {
			return &block
		}
	}

	return nil
}

// Validate checks if the block is valid
func (b Block) Validate(c Chain) bool {
	if b.calculateHash() != b.Hash {
		return false
	}

	if c.FindBlock(b.PreviousHash) == nil {
		return false
	}

	if c.blocks[b.Index].Hash != b.Hash {
		return false
	}

	return true
}

// Validate checks if the chain is valid
func (c *Chain) Validate() (bool, *Block) {
	for i := 1; i < len(c.blocks); i++ {
		if !c.blocks[i].Validate(*c) {
			return false, &c.blocks[i]
		}
	}

	return true, nil
}
