package gsuitemdm

//
// GSuiteMDM types for HTTP requests
//

// Action (Approve, Block, Delete, Wipe)
type ActionRequest struct {
	Action  string `json:"action"`
	Confirm bool   `json:"confirm"`
	Debug   bool   `json:"debug"`
	Domain  string `json:"domain"`
	IMEI    string `json:"imei"`
	Key     string `json:"key"`
	SN      string `json:"sn"`
}

// Individual directory entry
type DirectoryData struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phonenumber"`
}

// Multiple directory entries
type AllDirectoryData struct {
	Data []DirectoryData
}

// Search
type SearchRequest struct {
	Debug        bool   `json:"debug"`
	Domain       string `json:"domain"`
	ReturnFormat string `json:"format"`
	Key          string `json:"key"`
	QType        string `json:"qtype"`
	Q            string `json:"q"`
	SlackToken   string `json:"slacktoken"`
}

// Slack Search Request (nicked from https://github.com/nlopes/slack)
type SlackRequest struct {
	Token          string `json:"token"`
	TeamID         string `json:"team_id"`
	TeamDomain     string `json:"team_domain"`
	EnterpriseID   string `json:"enterprise_id,omitempty"`
	EnterpriseName string `json:"enterprise_name,omitempty"`
	ChannelID      string `json:"channel_id"`
	ChannelName    string `json:"channel_name"`
	UserID         string `json:"user_id"`
	UserName       string `json:"user_name"`
	Command        string `json:"command"`
	Text           string `json:"text"`
	ResponseURL    string `json:"response_url"`
	TriggerID      string `json:"trigger_id"`
}

// Update
type UpdateRequest struct {
	Debug bool   `json:"debug"`
	Key   string `json:"key"`
}

// EOF
