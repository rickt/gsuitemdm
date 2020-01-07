package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

// Constants
const (
	slackaccesslogsurl   string = "https://slack.com/api/team.accessLogs"
	slackuserinfourl     string = "https://slack.com/api/users.info"
	slackmobileuseragent string = "com.tinyspeck.chatlyio"
)

// Globals
var (
	page  int = 1
	pages int = 5
)

// Helper func to do a case-insensitive search
func caseinsensitivecontains(a, b string) bool {
	return strings.Contains(strings.ToUpper(a), strings.ToUpper(b))
}

// Get user info from the Slack API
func getuserinfo(userid string) *SlackUser {
	var su *SlackUser

	// Build the url
	uu := fmt.Sprintf("%s?token=%s&user=%s", slackuserinfourl, os.Getenv("SLACKTOKEN"), userid)
	// Create the request
	req, err := http.NewRequest("GET", uu, nil)
	if err != nil {
		log.Fatal("error: %s", err)
	}
	// Create the http client
	client := &http.Client{}
	// Get the response
	response, err := client.Do(req)
	if err != nil {
		log.Fatal("error: %s", err)
	}
	defer response.Body.Close()
	// Decode the JSON response into our user var
	if err := json.NewDecoder(response.Body).Decode(&su); err != nil {
		log.Println(err)
	}

	return su
}

// Get Slack access logs
func getaccesslogs(page int) *SlackAccessLog {
	var sal *SlackAccessLog

	// Build the url
	logsurl := fmt.Sprintf("%s?token=%s&count=1000&page=%d", slackaccesslogsurl, os.Getenv("SLACKTOKEN"), page)
	// Create the request
	req, err := http.NewRequest("GET", logsurl, nil)
	if err != nil {
		log.Fatal("error: %s", err)
	}
	// Create the http client
	client := &http.Client{}
	// Get the response
	response, err := client.Do(req)
	if err != nil {
		log.Fatal("error: %s", err)
	}
	defer response.Body.Close()
	// Decode the JSON response into our SlackAccessLog var
	if err := json.NewDecoder(response.Body).Decode(&sal); err != nil {
		log.Println(err)
	}

	return sal
}

func main() {
	// User map
	var um map[string]*UserInfo
	um = make(map[string]*UserInfo)

	// Get the Slack access logs one page at a time as per Slack API spec
	for page := 1; page < pages; page++ {
		// Get this page
		var sal *SlackAccessLog
		sal = getaccesslogs(page)
		// Range through the log entries in this page of the response
		for _, le := range sal.Logins {
			// We only care about mobile users
			if caseinsensitivecontains(le.UserAgent, slackmobileuseragent) {
				// Does this user already exist in the user map?
				if _, ok := um[le.UserID]; !ok {
					// It doesn't exist, get more info about the user from the Slack API
					var su *SlackUser
					su = getuserinfo(le.UserID)
					// Add this user to the user map
					um[le.UserID] = &UserInfo{
						SlackEmail:    su.User.Profile.Email,
						SlackName:     le.Username,
						SlackUserId:   le.UserID,
						SlackUserName: su.User.RealName,
					}
				}
			}
		}
	}

	// Show us the users who do not have MDM enabled
	for _, v := range um {
		if v.GSuiteMDM == false {
			fmt.Printf("%10.10s | %16.16s | %-17.17s | %-40.40s | %v\n", v.SlackUserId, v.SlackName, v.SlackUserName, v.SlackEmail, v.GSuiteMDM)
		}
	}
}

// EOF
