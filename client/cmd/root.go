/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

var apiEndpoint = "http://localhost:8080"
var difficulty = 1

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "blockchain",
	Short: "A command line client for the blockchain",
	Long: `A command line client for the blockchain, allowing you to add, list and verify blocks

 ** Configuration **
Set API_ENDPOINT env var to the address of the blockchain API server, e.g. http://myhost:9000
The default is localhost:8080`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolP("help", "", false, "help for this command")

	resp, err := http.Get(apiEndpoint + "/chain/difficulty")
	if err != nil {
		log.Println("Error calling blockchain API: ", err)
		log.Fatalln("Make sure the API is running and the API_ENDPOINT env var is set to the correct address")
	}

	diffResp := struct {
		Difficulty int
	}{}

	if err := json.NewDecoder(resp.Body).Decode(&diffResp); err != nil {
		log.Fatalln("Error decoding JSON response: ", err)
	}

	difficulty = diffResp.Difficulty
	fmt.Println("Block chain current difficulty:", difficulty)

	envVarEndpoint := os.Getenv("API_ENDPOINT")
	if envVarEndpoint != "" {
		apiEndpoint = envVarEndpoint
	}
}
