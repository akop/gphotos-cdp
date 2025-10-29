package main

import (
	"encoding/base64"
	"time"
)

// elapsedSince
func elapsedSince(start time.Time) string {
	elapsed := time.Since(start)
	return elapsed.String()
}

// isValidGPhotosId
func isValidGPhotosId(id string) bool {
	if len(id) < 20 {
		return false
	}
	_, err := base64.URLEncoding.DecodeString(id)
	return err == nil
}