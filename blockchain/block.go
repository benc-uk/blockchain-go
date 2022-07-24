package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
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

//
// newGenesisBlock creates the first block in the chain
//
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

//
// proofOfWork loops until the block hash starts with the given difficulty
//
func (b *Block) proofOfWork(diff int) {
	for b.Hash[:diff] != strings.Repeat("0", diff) {
		b.Nonce++
		b.Hash = b.calculateHash()
	}
}

//
// calculateHash calculates the hash of the block, note that the hash is not stored in the block
//
func (b *Block) calculateHash() string {
	hashInput := fmt.Sprintf("%s|%s|%s|%d", b.Timestamp, b.PreviousHash, b.Data, b.Nonce)

	h := sha256.New()
	h.Write([]byte(hashInput))
	hashed := h.Sum(nil)

	return hex.EncodeToString(hashed)
}

//
// For debugging blocks
//
func (b *Block) String() string {
	return fmt.Sprintf("Timestamp: %s\nPreviousHash: %s\nData: %s\nNonce: %d\nHash: %s\n",
		b.Timestamp, b.PreviousHash, b.Data, b.Nonce, b.Hash)
}

//
// Encodes the block as gob and returns the bytes
//
func (b *Block) Encode() []byte {
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	_ = enc.Encode(b)

	return buf.Bytes()
}

//
// Validate checks if the block is valid
//
func (b Block) Validate(c Chain) bool {
	if b.calculateHash() != b.Hash {
		return false
	}

	if prevBlock, err := c.Get(b.PreviousHash); err != nil || prevBlock == nil {
		return false
	}

	if thisBlock, err := c.Get(b.Hash); err != nil || thisBlock.Hash != b.Hash {
		return false
	}

	return true
}
