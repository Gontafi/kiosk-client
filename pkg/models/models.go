package models

type RegisterRequest struct {
	DeviceID      string  `json:"device_id"`
	Temperature   float64 `json:"temperature"`
	CPULoad       float64 `json:"cpu_load"`
	MemoryUsage   float64 `json:"memory_usage"`
	BrowserStatus string  `json:"browser_status"`
	Logs          *string `json:"logs"`
}

type URLResponse struct {
	URL string `json:"url"`
}

type HealthRequest struct {
	DeviceID      string  `json:"device_id"`
	Temperature   float64 `json:"temperature"`
	CPULoad       float64 `json:"cpu_load"`
	MemoryUsage   float64 `json:"memory_usage"`
	BrowserStatus string  `json:"browser_status"`
	Logs          *string `json:"logs"`
}
