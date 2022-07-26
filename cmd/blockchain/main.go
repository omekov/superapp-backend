package main

import (
	"fmt"

	"github.com/omekov/superapp-backend/internal/blockchain"
)

func main() {
	// create a new blockchain instance with a mining difficulty of 2
	blockchain := blockchain.CreateBlockChain(2)

	// record transactions on the blockchain for Alice, Bob, and John
	blockchain.AddBlock("Alice", "Bob", 5)
	blockchain.AddBlock("John", "Bob", 2)

	// check if the blockchain is valid; expecting true
	fmt.Println(blockchain.IsValid())
}
