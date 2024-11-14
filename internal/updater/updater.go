package updater

import (
	"kiosk-client/config"
	"kiosk-client/pkg/logger"
	"time"
)

func CheckForUpdates(cfg *config.Config, uuid *string) {
	for range time.Tick(cfg.PollInterval) {
		updateAvailable := checkUpdateAvailable(cfg, uuid)
		if updateAvailable {
			logger.Info("New update found, starting update process")
			applyUpdate()
		}
	}
}

func checkUpdateAvailable(cfg *config.Config, uuid *string) bool {
	// Logic to check with server if an update is available
	return false
}

func applyUpdate() {
	// Logic to download and apply update
}
