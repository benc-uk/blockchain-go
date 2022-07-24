package blockchain

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"time"

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
// AddBlock adds a block to the chain with the given data
//
func (c *Chain) AddBlock(data string) (string, error) {
	log.Println("Adding block c.last is ", c.last)

	b := Block{
		Timestamp:    time.Now().UTC(),
		Hash:         "FFFFFFFFFFFFFFF", // initial hash, will be overwritten
		PreviousHash: c.last,
		Data:         data,
		Nonce:        0,
	}

	// Calculate the hash of the block
	b.proofOfWork(c.Difficulty)

	// Incase another block is added while we are adding this one
	b.PreviousHash = c.last
	b.Hash = b.calculateHash()

	// Write the block to the db
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
		return "", err
	}

	return b.Hash, nil
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
// WITHOUT updating the hashes - this leads to a broken chain
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

	for i := 1; i < len(blocks); i++ {
		if !blocks[i].Validate(*c) {
			if blocks[i].PreviousHash == "" {
				continue
			}

			return false, &blocks[i]
		}
	}

	return true, nil
}
