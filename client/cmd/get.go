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

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a single block by its hash",
	Long:  `Get a single block by its hash`,

	Run: func(cmd *cobra.Command, args []string) {
		hash, _ := cmd.Flags().GetString("hash")
		if hash == "" {
			_ = cmd.Help()
		}

		resp, err := http.Get(apiEndpoint + "/block/" + hash)
		if err != nil {
			log.Fatalln("Error calling blockchain API: ", err)
		}

		if resp.StatusCode != http.StatusOK {
			log.Fatalln("Error calling blockchain API: ", resp.StatusCode)
		}

		defer resp.Body.Close()
		b := &blockchain.Block{}
		if err := json.NewDecoder(resp.Body).Decode(&b); err != nil {
			log.Fatalln("Error decoding JSON response: ", err)
		}

		fmt.Println(b.String())

	},
}

func init() {
	rootCmd.AddCommand(getCmd)

	getCmd.Flags().StringP("hash", "h", "", "Hash of the block to get")
}
