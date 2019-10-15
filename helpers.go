package gsuitemdm

//
// Various helper functions
//

import (
	"bufio"
	"cloud.google.com/go/logging"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"runtime"
	"strings"
	"time"
)

// Case-insensitive sort helper funcs
func (s DatastoreMobileDevices) Len() int {
	return len(s.Mobiledevices)
}

func (s DatastoreMobileDevices) Less(i, j int) bool {
	// return s.Mobiledevices[i].Name < s.Mobiledevices[j].Name
	return []rune(strings.ToLower(s.Mobiledevices[i].Name))[0] < []rune(strings.ToLower(s.Mobiledevices[j].Name))[0]
}

func (s DatastoreMobileDevices) Swap(i, j int) {
	s.Mobiledevices[i], s.Mobiledevices[j] = s.Mobiledevices[j], s.Mobiledevices[i]
}

// Helper function to check errors
func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// Helper function to check errors, Cloud Logging Edition (tm)
func CheckErrorCTX(err error, log *logging.Logger) {
	if err != nil {
		// TODO: remove this
		fmt.Printf("%s\n", err)
		log.Log(logging.Entry{
			Payload:  err,
			Severity: logging.Error,
		})
	}
}

// Helper function to ask for user confirmation in the CLI
func checkUserConfirmation(s string) bool {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("%s [y/n]: ", s)

		response, err := reader.ReadString('\n')
		checkError(err)

		response = strings.ToLower(strings.TrimSpace(response))

		if response == "y" || response == "yes" {
			return true
		} else if response == "n" || response == "no" {
			return false
		}
	}
}

// Helper function to do case-insensitive partial-string searches
func ciContains(a, b string) bool {

	return strings.Contains(strings.ToUpper(a), strings.ToUpper(b))
}

// Helper function to strip the domain name from an email address
func getEmailDomain(email string) string {

	components := strings.Split(email, "@")

	return components[1]
}

// Load main configuration file and return a config struct
func loadConfig(file string) (GSuiteMDMConfig, error) {

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
