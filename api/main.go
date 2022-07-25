package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
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

	diff := os.Getenv("CHAIN_DIFFICULTY")
	if diff == "" {
		diff = "5"
	}

	r := mux.NewRouter()
	r.Use(commonMiddleware)
	r.HandleFunc("/chain", listChain)
	r.HandleFunc("/chain/difficulty", getDifficulty)
	r.HandleFunc("/chain/validate", validateChain)
	r.HandleFunc("/block", add).Methods("POST")
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

	diffInt, _ := strconv.Atoi(diff)
	chain, err = blockchain.NewChain(diffInt)

	if err != nil {
		log.Fatalf("Could not create or open chain, %v !", err)
	}

	log.Println("### Blockchain API listening on port", port)
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

func add(w http.ResponseWriter, r *http.Request) {
	var block blockchain.Block

	err := json.NewDecoder(r.Body).Decode(&block)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = block.ValidateSimple(); err != nil {
		log.Println("### Error adding block", err)
		http.Error(w, "Invalid block payload "+err.Error(), http.StatusBadRequest)

		return
	}

	err = chain.AddBlock(&block)
	if err != nil {
		log.Println("### Error adding block", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	_ = json.NewEncoder(w).Encode(block)
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
		log.Println("### Chain integrity error, invalid block found: ", blockErr)
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

func getDifficulty(w http.ResponseWriter, r *http.Request) {
	_ = json.NewEncoder(w).Encode(map[string]int{"difficulty": chain.Difficulty})
}

func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.URL.Path)
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
