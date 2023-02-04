package main

import (
	"fmt"
	"os"
)

const defaultBingEndpoint = "https://api.bing.microsoft.com/"

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
}
