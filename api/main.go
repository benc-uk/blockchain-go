package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"blockchain-go/blockchain"

	// add spew

	"github.com/gorilla/mux"
)

var chain *blockchain.Chain

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r := mux.NewRouter()
	r.HandleFunc("/list", listChain)
	r.HandleFunc("/add", addTransaction).Methods("POST")

	srv := &http.Server{
		Handler:      r,
		Addr:         fmt.Sprintf("0.0.0.0:%s", port),
		WriteTimeout: 25 * time.Second,
		ReadTimeout:  25 * time.Second,
	}

	chain = blockchain.NewChain(2)

	log.Println("### Mini Blockchain API listening on port", port)
	log.Fatal(srv.ListenAndServe())
}

func listChain(w http.ResponseWriter, r *http.Request) {
	blocks := chain.GetBlocks()
	_ = json.NewEncoder(w).Encode(blocks)
}

func addTransaction(w http.ResponseWriter, r *http.Request) {
	var transaction blockchain.Transaction
	err := json.NewDecoder(r.Body).Decode(&transaction)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !transaction.Validate() {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	hash := chain.AddTransaction(transaction.Sender, transaction.Recipient, transaction.Amount)
	log.Println("### Added transaction:", transaction.String())

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(hash))
}
