package gsuitemdm

//
// GSuiteMDM types for HTTP requests
//

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
