package gsuitemdm

//
// GSuiteMDM various helper functions
//

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"runtime"
	"strings"
	"time"
)

// Sort funcs for devices
func (s DatastoreMobileDevices) Len() int {
	return len(s.Mobiledevices)
}
func (s DatastoreMobileDevices) Less(i, j int) bool {
	return []rune(strings.ToLower(s.Mobiledevices[i].Name))[0] < []rune(strings.ToLower(s.Mobiledevices[j].Name))[0]
}
func (s DatastoreMobileDevices) Swap(i, j int) {
	s.Mobiledevices[i], s.Mobiledevices[j] = s.Mobiledevices[j], s.Mobiledevices[i]
}

// Sort funcs for directory data
func (s AllDirectoryData) Len() int {
	return len(s.Data)
}
func (s AllDirectoryData) Less(i, j int) bool {
	return s.Data[i].Name < s.Data[j].Name
}
func (s AllDirectoryData) Swap(i, j int) {
	s.Data[i], s.Data[j] = s.Data[j], s.Data[i]
}

// Helper function to ask for user confirmation in the CLI
func checkUserConfirmation(s string) bool {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("%s [y/n]: ", s)

		response, err := reader.ReadString('\n')
		if err != nil {
			// TODO: fix this
		}

		response = strings.ToLower(strings.TrimSpace(response))

		if response == "y" || response == "yes" {
			return true
		} else if response == "n" || response == "no" {
			return false
		}
	}
}

// Helper function to strip the domain name from an email address
func getEmailDomain(email string) string {
	components := strings.Split(email, "@")

	return components[1]
}

// Helper function to get a remote IP from an http.Request
func GetIP(r *http.Request) string {
	fwd := r.Header.Get("X-FORWARDED-FOR")
	if fwd != "" {
		return fwd
	}
	return r.RemoteAddr
}

// Load configuration and return a config struct
func loadConfig(config string) (GSuiteMDMConfig, error) {
	var c GSuiteMDMConfig

	jp := json.NewDecoder(strings.NewReader(config))
	jp.Decode(&c)

	return c, nil
}

// Load main configuration file and return a config struct
// TODO: remove this func when all CF's have been converted to use Secret Manager
func loadConfigFile(file string) (GSuiteMDMConfig, error) {

	var c GSuiteMDMConfig

	// Open the main mdmtool configuration file
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

// Helper func to track how long a func takes to execute (found on StackExchange I think!)
func TimeTrack(start time.Time) {
	elapsed := time.Since(start)

	// Skip this function, and fetch the PC and file for its parent.
	pc, _, _, _ := runtime.Caller(1)

	// Retrieve a function object this functions parent.
	funcObj := runtime.FuncForPC(pc)

	// Regex to extract just the function name (and not the module path).
	runtimeFunc := regexp.MustCompile(`^.*\.(.*)$`)
	name := runtimeFunc.ReplaceAllString(funcObj.Name(), "$1")

	log.Println(fmt.Sprintf("DEBUG %s() took %s", name, elapsed))
}

// EOF
