package main

import ()

//
// Internal storage types
//

// User info type
type UserInfo struct {
	GSuiteName     string
	GSuiteEmail    string
	GSuiteMDM      bool
	SlackEmail     string
	SlackName      string
	SlackUserAgent string
	SlackUserId    string
	SlackUserName  string
}

//
// Slack types
//

// Slack access log types
type SlackAccessLog struct {
	Status     bool                  `json:"ok"`
	Logins     []SlackAccessLogEntry `json:"logins"`
	PagingData SlackAccessLogPaging  `json:"paging"`
}

type SlackAccessLogEntry struct {
	UserID    string `json:"user_id"`
	Username  string `json:"username"`
	DateFirst int64  `json:"date_first"`
	DateLast  int64  `json:"date_last"`
	Count     int    `json:"count"`
	IP        string `json:"ip"`
	UserAgent string `json:"user_agent"`
	ISP       string `json:"isp"`
	Country   string `json:"country"`
	Region    string `json:"region"`
}

type SlackAccessLogPaging struct {
	Count int `json:"count"`
	Total int `json:"total"`
	Page  int `json:"page"`
	Pages int `json:"pages"`
}

// Slack user types
type SlackUser struct {
	Ok   bool `json:"ok"`
	User struct {
		ID       string `json:"id"`
		TeamID   string `json:"team_id"`
		Name     string `json:"name"`
		Deleted  bool   `json:"deleted"`
		Color    string `json:"color"`
		RealName string `json:"real_name"`
		Tz       string `json:"tz"`
		TzLabel  string `json:"tz_label"`
		TzOffset int    `json:"tz_offset"`
		Profile  struct {
			Title                 string      `json:"title"`
			Phone                 string      `json:"phone"`
			Skype                 string      `json:"skype"`
			RealName              string      `json:"real_name"`
			RealNameNormalized    string      `json:"real_name_normalized"`
			DisplayName           string      `json:"display_name"`
			DisplayNameNormalized string      `json:"display_name_normalized"`
			Fields                interface{} `json:"fields"`
			StatusText            string      `json:"status_text"`
			StatusEmoji           string      `json:"status_emoji"`
			StatusExpiration      int         `json:"status_expiration"`
			AvatarHash            string      `json:"avatar_hash"`
			ImageOriginal         string      `json:"image_original"`
			IsCustomImage         bool        `json:"is_custom_image"`
			Email                 string      `json:"email"`
			FirstName             string      `json:"first_name"`
			LastName              string      `json:"last_name"`
			Image24               string      `json:"image_24"`
			Image32               string      `json:"image_32"`
			Image48               string      `json:"image_48"`
			Image72               string      `json:"image_72"`
			Image192              string      `json:"image_192"`
			Image512              string      `json:"image_512"`
			Image1024             string      `json:"image_1024"`
			StatusTextCanonical   string      `json:"status_text_canonical"`
			Team                  string      `json:"team"`
		} `json:"profile"`
		IsAdmin           bool `json:"is_admin"`
		IsOwner           bool `json:"is_owner"`
		IsPrimaryOwner    bool `json:"is_primary_owner"`
		IsRestricted      bool `json:"is_restricted"`
		IsUltraRestricted bool `json:"is_ultra_restricted"`
		IsBot             bool `json:"is_bot"`
		IsAppUser         bool `json:"is_app_user"`
		Updated           int  `json:"updated"`
		Has2Fa            bool `json:"has_2fa"`
	} `json:"user"`
}

// EOF
