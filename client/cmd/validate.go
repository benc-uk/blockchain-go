/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"log"
	"net/http"

	"github.com/spf13/cobra"
)

// validateCmd represents the validate command
var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Check integrity",
	Long:  `Validate the block chain integrity`,

	Run: func(cmd *cobra.Command, args []string) {
		resp, err := http.Get(apiEndpoint + "/chain/validate")
		if err != nil {
			log.Fatalln("Error calling blockchain API: ", err)
		}

		if resp.StatusCode != http.StatusOK {
			log.Fatalln("Blockchain integrity failure", resp.StatusCode)
		} else {
			log.Println("Blockchain integrity is OK")
		}
	},
}

func init() {
	rootCmd.AddCommand(validateCmd)
}
