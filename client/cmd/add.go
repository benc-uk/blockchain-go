/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"blockchain-go/blockchain"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a transaction to the blockchain",
	Long:  `Adds a transaction to the blockchain, by calculating the hash of the block using proof of work.`,

	Run: func(cmd *cobra.Command, args []string) {
		sender, _ := cmd.Flags().GetString("sender")
		recipient, _ := cmd.Flags().GetString("recipient")
		amount, _ := cmd.Flags().GetFloat64("amount")

		if sender == "" || recipient == "" || amount <= 0 {
			_ = cmd.Help()
			return
		}

		block, err := blockchain.NewTransactionBlock(sender, recipient, amount)
		if err != nil {
			log.Fatalln("Error creating transaction block: ", err)
			return
		}

		spinnerCharSet := []string{"█▒▒▒▒▒▒▒▒▒", "██▒▒▒▒▒▒▒▒", "███▒▒▒▒▒▒▒", "████▒▒▒▒▒▒", "█████▒▒▒▒▒", "██████▒▒▒▒", "███████▒▒▒", "████████▒▒", "█████████▒", "██████████"}
		spinner := spinner.New(spinnerCharSet, 100*time.Millisecond)
		_ = spinner.Color("magenta")
		spinner.Prefix = "Mining in progress "
		spinner.Start()
		start := time.Now()

		// Do the actual work !!!
		block.Mine(difficulty)

		spinner.Stop()
		elapsed := time.Since(start)
		fmt.Println("Block mined:", block.Nonce, block.Hash)
		fmt.Println("Process took:", elapsed)

		data, err := json.Marshal(block)
		if err != nil {
			log.Fatalln("JSON error: ", err)
		}

		_, err = http.Post(apiEndpoint+"/block", "application/json", bytes.NewBuffer(data))
		if err != nil {
			log.Fatalln("Error sending POST request: ", err)
		}

		fmt.Println("Transaction added to the blockchain!")
	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	addCmd.Flags().StringP("sender", "s", "", "Name of sender")
	addCmd.Flags().StringP("recipient", "r", "", "Name of receiver")
	addCmd.Flags().Float64P("amount", "a", 0.0, "Amount to send")
}
