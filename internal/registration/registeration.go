package registration

import (
	"kiosk-client/config"
	"kiosk-client/pkg/logger"
	"kiosk-client/pkg/models"
	"kiosk-client/pkg/utils"
	"kiosk-client/internal/monitor"
)

func RegisterDevice(cfg *config.Config) string {
	uuid := utils.LoadOrCreateUUID()
	
	logger.Info("Trying to register uuid ", cfg.ServerURL+cfg.RegistrationPath, uuid)
	
	status := monitor.CollectHealthData(&uuid)
	
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
