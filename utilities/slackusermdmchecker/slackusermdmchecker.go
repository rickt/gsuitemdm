package slackusermdmchecker

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"
)

// Constants
const (
	datelayout                  = "January 2, 2006"
	slackaccesslogsurl   string = "https://slack.com/api/team.accessLogs"
	slackuserinfourl     string = "https://slack.com/api/users.info"
	slackmobileuseragent string = "com.tinyspeck.chatlyio"
)

// Globals
var (
	mdmstatusurl string = os.Getenv("GSUITEMDMURL") + "/SearchDatastore"
	page         int    = 1
	pages        int    = 5
)

// Sort funcs for users
func (s Users) Len() int {
	return len(s)
}
func (s Users) Less(i, j int) bool {
	return s[i].SlackName < s[j].SlackName
}
func (s Users) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Build HTML email
func buildhtmlemail(cu, ncu Users) []byte {
	// Basic setup
	from := mail.NewEmail(os.Getenv("FROM_NAME"), os.Getenv("FROM_ADDR"))
	to := mail.NewEmail(os.Getenv("RECIPIENTS_NAME"), os.Getenv("RECIPIENTS_ADDR"))
	t := time.Now()
	subject := "Company MDM Compliance Report for " + t.Format(datelayout)
	replyto := mail.NewEmail(os.Getenv("REPLYTO_NAME"), os.Getenv("REPLYTO_ADDR"))

	// Create the plaintext email body
	var body string

	// Start with non-compliant users
	body = fmt.Sprintf("<p><strong>(%d) Active GIGANIK Slack staff using a personal phone or company phone with no MDM to login to Slack:</strong><br>", len(ncu))
	for _, x := range ncu {
		body = body + fmt.Sprintf("&nbsp;&nbsp;&nbsp;%s (<a href=\"%s=%s\">@%s</a> &lt;%s&gt;)<br>", x.SlackUserName, os.Getenv("SLACKURL"), x.SlackUserId, x.SlackName, x.SlackEmail)
	}
	body = body + "</p>"
	// Now the users with MDM
	body = body + fmt.Sprintf("<p><strong>(%d) Active GIGANIK Slack staff using an MDM-protected company phone to login to Slack:\n</strong><br>", len(cu))
	for _, x := range cu {
		body = body + fmt.Sprintf("&nbsp;&nbsp;&nbsp;%s (<a href=\"%s=%s\">@%s</a> &lt;%s&gt;)<br>", x.SlackUserName, os.Getenv("SLACKURL"), x.SlackUserId, x.SlackName, x.SlackEmail)
	}
	body = body + "</p>"
	body = body + fmt.Sprintf("<p><br>This email generated & sent every weekday at 06:45PST by %s/%s.<br>", os.Getenv("GSUITEMDMURL"), os.Getenv("APPNAME"))
	body = body + fmt.Sprintf("Source code available <a href=\"%s\">here</a>.<br>", os.Getenv("SRCURL"))
	body = body + "</p>"

	// Build the message
	content := mail.NewContent("text/html", body)
	m := mail.NewV3MailInit(from, subject, to, content)
	m.SetReplyTo(replyto)

	return mail.GetRequestBody(m)
}

// Helper func to do a case-insensitive search
func caseinsensitivecontains(a, b string) bool {
	return strings.Contains(strings.ToUpper(a), strings.ToUpper(b))
}

// Get MDM status
func getmdmstatus(email string) bool {
	var mr MDMRequest

	mr.Key = os.Getenv("GSUITEMDMTOKEN")
	mr.QType = "email"
	mr.Q = email

	req, err := json.Marshal(mr)
	if err != nil {
		log.Fatal("error: %s", err)
	}

	// Create the request
	resp, err := http.Post(mdmstatusurl, "application/json", bytes.NewBuffer(req))
	if err != nil {
		log.Fatal("error: %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		// User has MDM
		return true
	} else {
		// User does not have MDM
		return false
	}
}

// Get Slack access logs
func getslackaccesslogs(page int) *SlackAccessLog {
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

// Get user info from the Slack API
func getslackuserinfo(userid string) *SlackUser {
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

// Send email
func sendemail(cu, ncu Users) {
	// Create the Sendgrid API request
	req := sendgrid.GetRequest(os.Getenv("SENDGRID_API_KEY"), "/v3/mail/send", "https://api.sendgrid.com")
	req.Method = "POST"

	// Create the email body
	var Body []byte = buildhtmlemail(cu, ncu)
	req.Body = Body

	// Make the request to the Sendgrid API
	resp, err := sendgrid.API(req)
	if err != nil {
		log.Fatal("error: %s", err)
	}
	if resp.StatusCode == 202 {
		fmt.Printf("\nEmail sent.\n\n")
	}
}

// Generate a list of naughty and nice users
func SlackUserMDMChecker(w http.ResponseWriter, r *http.Request) {
	// User map
	var um map[string]*UserInfo
	um = make(map[string]*UserInfo)

	// Get the Slack access logs one page at a time as per Slack API spec
	for page := 1; page < pages; page++ {
		// Get this page
		var sal *SlackAccessLog
		sal = getslackaccesslogs(page)
		// Range through the log entries in this page of the response
		for _, le := range sal.Logins {
			// We only care about mobile users
			if caseinsensitivecontains(le.UserAgent, slackmobileuseragent) {
				// Does this user already exist in the user map?
				if _, ok := um[le.UserID]; !ok {
					// It doesn't exist, get more info about the user from the Slack API
					var su *SlackUser
					su = getslackuserinfo(le.UserID)
					// Add this user to the user map
					um[le.UserID] = &UserInfo{
						GSuiteMDM:     getmdmstatus(su.User.Profile.Email),
						SlackEmail:    su.User.Profile.Email,
						SlackName:     le.Username,
						SlackUserId:   le.UserID,
						SlackUserName: su.User.RealName,
					}
				}
			}
		}
	}

	// Range through our map and find MDM-compliant/non-compliant users
	var cu, ncu Users
	for _, v := range um {
		switch {
		case v.GSuiteMDM == true:
			// User has MDM
			cu = append(cu, v)
			break

		case v.GSuiteMDM == false:
			// User does not have MDM
			ncu = append(ncu, v)
			break
		}
	}

	// Sort the data
	sort.Sort(cu)
	sort.Sort(ncu)

	// Print out a report
	fmt.Fprintf(w, "(%d) Active GIGANIK Slack staff using a personal phone or company phone with no MDM to login to Slack.\n", len(ncu))
	fmt.Fprintf(w, "\n")
	fmt.Fprintf(w, "(%d) Active GIGANIK Slack staff using an MDM-protected company phone to login to Slack\n", len(cu))

	// Send email
	sendemail(cu, ncu)
	fmt.Fprintf(w, "\n")
	fmt.Fprintf(w, "Email sent.\n")
	fmt.Fprintf(w, "\n")
	fmt.Fprintf(w, "%s Success\n", os.Getenv("APPNAME"))

	return

}

// EOF
