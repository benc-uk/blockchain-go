package blockchain

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"

	//"github.com/boltdb/bolt"

	bolt "go.etcd.io/bbolt"
)

type Chain struct {
	db         *bolt.DB
	last       string
	Difficulty int
}

//
// NewChain creates a new chain with the given difficulty
//
func NewChain(diff int) (*Chain, error) {
	// Open the db
	db, err := bolt.Open("blockchain.db", 0600, nil)
	if err != nil {
		log.Fatalf("Could not open db, %v !", err)
	}

	// Create the bucket, and check last hash
	lastHash := ""
	err = db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("blocks"))

		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}

		lastHash = string(bucket.Get([]byte("last")))

		// Empty chain, create the genesis block
		if lastHash == "" {
			log.Println("No last block found!")
			log.Println("Initialising a new chain with genesis block")
			genesis := newGenesisBlock(diff)

			err := bucket.Put([]byte(genesis.Hash), genesis.Encode())
			if err != nil {
				return fmt.Errorf("put genesis block: %s", err)
			}

			err = bucket.Put([]byte("last"), []byte(genesis.Hash))
			if err != nil {
				return fmt.Errorf("update last: %s", err)
			}

			lastHash = genesis.Hash
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	log.Println("Chain db opened, last block:", lastHash)
	log.Println("Current chain difficulty:", diff)

	return &Chain{
		db:         db,
		Difficulty: diff,
		last:       lastHash,
	}, nil
}

func (c Chain) GetLastHash() string {
	return c.last
}

//
// GetBlocks returns the blocks of the chain, no paging, will NOT scale!
//
func (c *Chain) GetBlocks() ([]Block, error) {
	blocks := []Block{}
	err := c.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("blocks"))

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			// Last is a special case, doesn't hold a block, just a hash
			if string(k) == "last" {
				continue
			}

			block := Block{}
			decoder := gob.NewDecoder(bytes.NewBuffer(v))
			if err := decoder.Decode(&block); err != nil {
				return fmt.Errorf("decode: %s", err)
			}

			blocks = append(blocks, block)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return blocks, nil
}

//
// AddBlock adds a mined block to the chain, checking the block first
//
func (c *Chain) AddBlock(b *Block) error {
	// count leading zeros
	zeros := 0

	for _, v := range b.Hash {
		if v == '0' {
			zeros++
		} else {
			break
		}
	}

	// Check the hash of the block
	if b.CalculateHash() != b.Hash {
		return fmt.Errorf("block hash does not match")
	}

	// Validate the block's hash difficulty
	if zeros < c.Difficulty {
		return fmt.Errorf("block not fully mined to current difficulty")
	}

	// Write the block to the db
	b.PreviousHash = c.last
	err := c.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("blocks"))

		err := bucket.Put([]byte(b.Hash), b.Encode())
		if err != nil {
			return fmt.Errorf("put block: %s", err)
		}

		err = bucket.Put([]byte("last"), []byte(b.Hash))
		if err != nil {
			return fmt.Errorf("update last: %s", err)
		}
		c.last = b.Hash

		return nil
	})

	if err != nil {
		return err
	}

	log.Printf("Added block:\n%s", b.String())

	return nil
}

//
// Get a block by its hash
//
func (c *Chain) Get(hash string) (*Block, error) {
	b := &Block{}
	err := c.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("blocks"))
		blockBytes := bucket.Get([]byte(hash))
		if blockBytes == nil {
			return fmt.Errorf("block not found")
		}

		decoder := gob.NewDecoder(bytes.NewBuffer(blockBytes))
		if err := decoder.Decode(b); err != nil {
			return fmt.Errorf("decode: %s", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return b, nil
}

//
// UpdateBlock updates the data of the block,
// WITHOUT updating the hashes - this leads to an invalid chain
//
func (c *Chain) UpdateBlock(b Block, newData string) {
	b.Data = newData
	// Write the block to the db
	_ = c.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("blocks"))

		err := bucket.Put([]byte(b.Hash), b.Encode())
		if err != nil {
			return fmt.Errorf("put block: %s", err)
		}

		err = bucket.Put([]byte("last"), []byte(b.Hash))
		if err != nil {
			return fmt.Errorf("update last: %s", err)
		}

		return nil
	})
}

//
// Validate checks if the chain is valid
//
func (c *Chain) Validate() (bool, *Block) {
	blocks, err := c.GetBlocks()
	if err != nil {
		return false, nil
	}

	for i := 0; i < len(blocks); i++ {
		// Skip the genesis block
		if blocks[i].PreviousHash == "" {
			continue
		}

		if !blocks[i].Validate(*c) {
			return false, &blocks[i]
		}
	}

	return true, nil
}
