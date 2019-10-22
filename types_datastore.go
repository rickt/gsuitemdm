package gsuitemdm

//
// GSuiteMDM types for Datastore
//

// A single mobile device
type DatastoreMobileDevice struct {
	Color             string // Color of device
	CompromisedStatus string // Is the device compromised?
	Domain            string // G Suite domain
	DeveloperMode     bool   // Is the device in developer mode?
	Email             string // Email address of device owner
	IMEI              string // IMEI
	Model             string // Model
	Name              string // Full Name of device owner
	Notes             string // Notes
	OS                string // Operating System
	OSBuild           string // OS Build
	PhoneNumber       string // Telephone number of the device
	RAM               string // RAM in GB
	ResourceId        string // MDM ID for device
	SN                string // Serial number
	Status            string // Device status
	SyncFirst         string // First sync device time
	SyncLast          string // Most recent device sync time
	Type              string // Type of G Suite sync
	UnknownSources    bool   // Are unknown sources of apps allowed on the device?
	USBADB            bool   // Is ADB/USB debugging enabled?
	WifiMac           string // Wifi MAC address
}

// Multiple mobile devices
type DatastoreMobileDevices struct {
	Mobiledevices []DatastoreMobileDevice
}

// EOF
