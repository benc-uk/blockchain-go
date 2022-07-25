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
//
//
func NewBlock(data string) *Block {
	b := Block{
		Timestamp:    time.Now().UTC(),
		Hash:         "fffffffffffffff", // initial hash, will be overwritten
		PreviousHash: "fffffffffffffff",
		Data:         data,
		Nonce:        0,
	}

	return &b
}

//
// newGenesisBlock creates the first block in the chain
//
func newGenesisBlock(diff int) Block {
	b := Block{
		Timestamp:    time.Now().UTC(),
		Hash:         "fffffffffffffff",
		PreviousHash: "",
		Data:         "Let there be light",
		Nonce:        0,
	}

	b.Mine(diff)

	return b
}

//
// Mine carries out proof of work hashing, until the block hash starts with the given difficulty
//
func (b *Block) Mine(diff int) {
	for b.Hash[:diff] != strings.Repeat("0", diff) {
		b.Nonce++
		b.Hash = b.CalculateHash()
	}
}

//
// CalculateHash calculates the hash of the block, note that the hash is not stored in the block
//
func (b *Block) CalculateHash() string {
	// **** SIMPLIFICATION!!! *****************************************
	// We DO NOT include the previous hash in the hash calculation
	// In a simplified system such as this, the previous hash is not known to mining clients
	// The last hash is known server side and could be used and sent/fetched to mining clients
	// However chances are will be out of date by the time it has been mined as other blocks will have been added
	// ****************************************************************
	hashInput := fmt.Sprintf("%s|%s|%d", b.Timestamp, b.Data, b.Nonce)

	h := sha256.New()
	h.Write([]byte(hashInput))
	hashed := h.Sum(nil)

	return hex.EncodeToString(hashed)
}

//
// For debugging blocks
//
func (b *Block) String() string {
	return fmt.Sprintf("Hash: %s\nPreviousHash: %s\nData: %s\nNonce: %d\nTimestamp: %s\n",
		b.Hash, b.PreviousHash, b.Data, b.Nonce, b.Timestamp)
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
// Validate checks if the block on the chain is valid
//
func (b Block) Validate(c Chain) bool {
	if b.CalculateHash() != b.Hash {
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

//
// ValidateSimple checks if the block is ready to be added
//
func (b Block) ValidateSimple() error {
	ch := b.CalculateHash()
	if ch != b.Hash {
		return fmt.Errorf("invalid hash: %s", ch)
	}

	if b.Timestamp.IsZero() || b.Data == "" || b.Hash == "" || b.Nonce < 0 {
		return fmt.Errorf("missing field(s)")
	}

	return nil
}
