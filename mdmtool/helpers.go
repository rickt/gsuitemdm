package main

//
// MDMTool helper funcs
//

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
)

// Helper function to ask for user confirmation
func checkUserConfirmation(s string) bool {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("%s [y/n]: ", s)

		response, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		response = strings.ToLower(strings.TrimSpace(response))

		if response == "y" || response == "yes" {
			return true
		} else if response == "n" || response == "no" {
			return false
		}
	}
}

// Load MDMTool configuration
func loadMDMToolConfig(file string) (MDMToolConfig, error) {
	var c MDMToolConfig

	// Open the MDMTool configuration file
	cf, err := os.Open(file)
	defer cf.Close()
	if err != nil {
		return c, err
	}

	// Decode the JSON
	jp := json.NewDecoder(cf)
	jp.Decode(&c)

	return c, nil

}

// EOF
