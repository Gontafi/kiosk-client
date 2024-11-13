package models

type RegisterRequest struct {
	DeviceID string `json:"device_id"`
}

type URLResponse struct {
	URL string `json:"url"`
}

type HealthRequest struct {
	DeviceID      string `json:"device_id"`
	Temperature   string `json:"temperature"`
	CPULoad       string `json:"cpu_load"`
	MemoryUsage   string `json:"memory_usage"`
	BrowserStatus string `json:"browser_status"`
	Logs          string `json:"logs"`
}

type CheckUpdateResponse struct {
	UpdateAvailable string `json:"update_available"`
	Version         string `json:"version"`
	UpdateURL       string `json:"update_url"`
	Checksum        string `json:"checksum"`
}
