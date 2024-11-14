package registration

import (
	"kiosk-client/config"
	"kiosk-client/pkg/logger"
	"kiosk-client/pkg/models"
	"kiosk-client/pkg/utils"
)

func RegisterDevice(cfg *config.Config) string {
	uuid := utils.LoadOrCreateUUID()

	_, _, err := utils.MakePOSTRequest(
		cfg.ServerURL+cfg.RegistrationPath,
		models.RegisterRequest{DeviceID: uuid})

	if err != nil {
		logger.Error("Registration failed:", err)
	} else {
		logger.Info("Device registered with UUID:", uuid)
	}

	return uuid
}
