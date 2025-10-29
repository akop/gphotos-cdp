package main

import (
	"encoding/base64"
	"time"

	"github.com/rs/zerolog"
)

// logInfoElapsed
func logInfoElapsed(log zerolog.Logger, start time.Time, name string) {
	elapsed := time.Since(start)
	log.Info().Msgf("%s in %s", name, elapsed)
}

// logDebugElapsed
func logDebugElapsed(log zerolog.Logger, start time.Time, name string) {
	elapsed := time.Since(start)
	log.Debug().Msgf("%s in %s", name, elapsed)
}

// isValidGPhotosId
func isValidGPhotosId(id string) bool {
	if len(id) < 20 {
		return false
	}
	_, err := base64.URLEncoding.DecodeString(id)
	return err == nil
}
