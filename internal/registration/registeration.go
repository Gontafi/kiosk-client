package registration

import (
	"kiosk-client/config"
	"kiosk-client/pkg/logger"
	"kiosk-client/pkg/utils"
)

func RegisterDevice(cfg *config.Config) string {
	uuid := utils.LoadOrCreateUUID()

	resp, err := utils.MakePOSTRequest(cfg.ServerURL+cfg.RegistrationPath, map[string]string{"uuid": uuid})
	if err != nil {
		logger.Error("Registration failed:", err)
	} else {
		logger.Info("Device registered with UUID:", uuid)
	}

	_ = resp

	return uuid
}
