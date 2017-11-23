package config

import (
	"os"

	"github.com/tokubai/kinu/logger"
)

var (
	BackwardCompatibleMode = false
)

func init() {
	if len(os.Getenv("KINU_BACKWARD_COMPATIBLE_MODE")) != 0 {
		BackwardCompatibleMode = true
		logger.Warn("running backward compaztible mode. this mode is deprecated.")
	}
}
