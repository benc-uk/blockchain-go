package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"blockchain-go/blockchain"

	"github.com/gorilla/mux"
)

var chain *blockchain.Chain

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r := mux.NewRouter()
	r.Use(commonMiddleware)
	r.HandleFunc("/chain", listChain)
	r.HandleFunc("/chain/validate", validateChain)
	r.HandleFunc("/block", addTransaction).Methods("POST")
	r.HandleFunc("/block/{hash}", get)
	r.HandleFunc("/block/tamper/{hash}", tamper).Methods("PUT")
	r.HandleFunc("/block/validate/{hash}", validateBlock)

	srv := &http.Server{
		Handler:           r,
		Addr:              fmt.Sprintf("0.0.0.0:%s", port),
		WriteTimeout:      25 * time.Second,
		ReadTimeout:       25 * time.Second,
		ReadHeaderTimeout: 25 * time.Second,
	}

	var err error
	chain, err = blockchain.NewChain(5)

	if err != nil {
		log.Fatalf("Could not create or open chain, %v !", err)
	}

	log.Println("### Mini Blockchain API listening on port", port)
	log.Fatal(srv.ListenAndServe())
}

func listChain(w http.ResponseWriter, r *http.Request) {
	blocks, err := chain.GetBlocks()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	type listResult struct {
		Blocks     []blockchain.Block
		LastHash   string
		Difficulty int
	}

	_ = json.NewEncoder(w).Encode(listResult{
		Blocks:     blocks,
		LastHash:   chain.GetLastHash(),
		Difficulty: chain.Difficulty,
	})
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

	hash, err := chain.AddTransaction(transaction.Sender, transaction.Recipient, transaction.Amount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("### Added transaction:", transaction.String())

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ADDED", "hash": hash})
}

func tamper(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hash := vars["hash"]
	log.Println("### Tampering with block", hash)

	b, err := chain.Get(hash)
	if err != nil {
		http.Error(w, "Block not found: "+err.Error(), http.StatusNotFound)
		return
	}

	chain.UpdateBlock(*b, "Tampered data")
}

func validateBlock(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hash := vars["hash"]
	log.Println("### Validating block", hash)

	b, err := chain.Get(hash)
	if err != nil {
		http.Error(w, "Block not found: "+err.Error(), http.StatusNotFound)
		return
	}

	ok := b.Validate(*chain)
	if ok {
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "OK"})
	} else {
		w.WriteHeader(http.StatusNotAcceptable)
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "INTEGRITY ERROR"})
	}
}

func validateChain(w http.ResponseWriter, r *http.Request) {
	ok, blockErr := chain.Validate()
	if ok {
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "OK"})
	} else {
		w.WriteHeader(http.StatusNotAcceptable)
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "INTEGRITY ERROR", "block_hash": blockErr.Hash})
	}
}

func get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hash := vars["hash"]

	b, err := chain.Get(hash)
	if err != nil {
		http.Error(w, "Block not found: "+err.Error(), http.StatusNotFound)
		return
	}

	_ = json.NewEncoder(w).Encode(b)
}

func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.URL.Path)
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
