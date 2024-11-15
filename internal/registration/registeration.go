package registration

import (
	"kiosk-client/config"
	"kiosk-client/internal/monitor"
	"kiosk-client/pkg/logger"
	"kiosk-client/pkg/models"
	"kiosk-client/pkg/utils"
)

func RegisterDevice(cfg *config.Config) string {
	uuid := utils.LoadOrCreateUUID()

	status := monitor.CollectHealthData(cfg, &uuid)

	_, _, err := utils.MakePOSTRequest(
		cfg.ServerURL+cfg.RegistrationPath,
		models.RegisterRequest{
			DeviceID:      uuid,
			Temperature:   status.Temperature,
			CPULoad:       status.CPULoad,
			MemoryUsage:   status.MemoryUsage,
			BrowserStatus: status.BrowserStatus,
			Logs:          status.Logs,
		},
	)

	if err != nil {
		logger.Error("Registration failed:", err)
	} else {
		logger.Info("Device registered with UUID:", uuid)
	}

	return uuid
}
