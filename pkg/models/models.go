package models

type RegisterRequest struct {
	DeviceID 	  string `json:"device_id"`
	Temperature   int `json:"temperature"`
	CPULoad       int `json:"cpu_load"`
	MemoryUsage   int `json:"memory_usage"`
	BrowserStatus string `json:"browser_status"`
	Logs          *string `json:"logs"`
}

type URLResponse struct {
	URL string `json:"url"`
}

type HealthRequest struct {
	DeviceID      string `json:"device_id"`
	Temperature   int `json:"temperature"`
	CPULoad       int `json:"cpu_load"`
	MemoryUsage   int `json:"memory_usage"`
	BrowserStatus string `json:"browser_status"`
	Logs          *string `json:"logs"`
}

type CheckUpdateResponse struct {
	UpdateAvailable string `json:"update_available"`
	Version         string `json:"version"`
	UpdateURL       string `json:"update_url"`
	Checksum        string `json:"checksum"`
}
