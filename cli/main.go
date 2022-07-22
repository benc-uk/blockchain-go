package main

import (
	"blockchain-go/blockchain"
	"fmt"

	"github.com/davecgh/go-spew/spew"
)

func main() {
	chain := blockchain.NewChain(2)
	chain.AddBlock("Send 1 BTC to Ivan")
	chain.AddBlock("Send 1453434 BTC to BEN")
	chain.AddTransaction("Ivan", "Alice", 1)
	chain.AddTransaction("Ben", "Phil", 547)
	chain.AddTransaction("David", "Monty", 1312)
	fmt.Println("==================================")
	spew.Dump(chain)
}
