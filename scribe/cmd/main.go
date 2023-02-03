package main

import (
	"fmt"
	"os"
)

const defaultSpecialty = "RADIOLOGY"

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
}
