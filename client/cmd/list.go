/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"blockchain-go/blockchain"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/spf13/cobra"
)

type ListResp struct {
	Blocks     []blockchain.Block
	LastHash   string
	Difficulty int
}

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List the whole block chain",
	Long:  `List the whole block chain`,

	Run: func(cmd *cobra.Command, args []string) {
		resp, err := http.Get(apiEndpoint + "/chain")
		if err != nil {
			log.Fatalln("Error calling blockchain API: ", err)
		}

		defer resp.Body.Close()
		listData := &ListResp{}
		if err := json.NewDecoder(resp.Body).Decode(&listData); err != nil {
			log.Fatalln("Error decoding JSON response: ", err)
		}

		for _, block := range listData.Blocks {
			fmt.Println(block.String())
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
