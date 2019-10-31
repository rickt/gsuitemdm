package main

//
// MDMTool helper funcs
//

import (
	"encoding/json"
	"os"
)

// Load MDMTool configuration
func loadMDMToolURLs(file string) (MDMToolURLs, error) {
	var c MDMToolURLs

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
