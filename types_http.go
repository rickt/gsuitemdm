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
	PhoneNumber string `json:"phonenumbner"`
}

// Search
type SearchRequest struct {
	Debug  bool   `json:"debug"`
	Domain string `json:"domain"`
	Key    string `json:"key"`
	QType  string `json:"qtype"`
	Q      string `json:"q"`
}

// Update
type UpdateRequest struct {
	Debug bool   `json:"debug"`
	Key   string `json:"key"`
}

// EOF
