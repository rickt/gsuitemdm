package gsuitemdm

//
// GSuiteMDM types for HTTP requests
//

// Action (Approve, Block, Delete, Wipe)
type DeviceActionRequest struct {
	Action string `json:"action"`
	Debug  bool   `json:"debug"`
	IMEI   string `json:"imei"`
	Key    string `json:"key"`
	SN     string `json:"sn"`
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
