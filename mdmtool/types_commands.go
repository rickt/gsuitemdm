package main

//
// Types for mdmtool commands
//

// ApproveCommand ...
type ApproveCommand struct {
	Domain string
	IMEI   string
	SN     string
}

// BlockCommand ...
type BlockCommand struct {
	Domain string
	IMEI   string
	SN     string
}

// DeleteCommand ...
type DeleteCommand struct {
	Domain string
	IMEI   string
	SN     string
}

// DirectoryCommand ...
type DirectoryCommand struct {
	Email string
	Name  string
}

// SearchCommand ...
type SearchCommand struct {
	All     bool
	Domain  string
	Email   string
	IMEI    string
	Name    string
	Notes   string
	Phone   string
	SN      string
	Status  string
	Verbose bool
}

// ShowDomainsCommand ..
type ShowDomainsCommand struct {
	Verbose bool
}

// UpdateSheetCommand ...
type UpdateSheetCommand struct {
	Verbose bool
}

// UpdateDatastoreCommand ...
type UpdateDatastoreCommand struct {
	Verbose bool
}

// WipeCommand ...
type WipeCommand struct {
	Domain string
	IMEI   string
	SN     string
}

// EOF
