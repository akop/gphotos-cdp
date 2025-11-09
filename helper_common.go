package main

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"slices"
	"strings"
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


func imageIdFromUrl(location string) (string, error) {
	// Parse the URL
	u, err := url.Parse(location)
	if err != nil {
		return "", fmt.Errorf("invalid URL %v: %w", location, err)
	}

	// Split the path into segments
	parts := strings.Split(strings.Trim(u.Path, "/"), "/")

	// Look for "photo" segment and ensure there's a following segment
	for i := 0; i < len(parts)-1; i++ {
		if parts[i] == "photo" {
			return parts[i+1], nil
		}
	}
	return "", fmt.Errorf("could not find /photo/{imageId} pattern in URL: %v", location)
}

func absInt(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// Compare two file names, where s2 is sometimes mangled by gphotos
// there will be some false positives, this is ok
func compareMangled(_s1, _s2 string) bool {
	doCompare := func(s1, s2 string) bool {
		sr1 := []rune(s1)
		sr2 := []rune(s2)

		l1 := len(sr1)
		for i := range slices.Backward(sr1) {
			if sr1[i] == '.' {
				l1 = i + 1
				break
			}
		}

		l2 := len(sr2)
		for i := range slices.Backward(sr2) {
			if sr2[i] == '.' {
				l2 = i + 1
				break
			}
		}

		i1 := 0
		for i1 < len(sr1) && sr1[i1] == '.' && sr2[i1] != '.' {
			i1++
		}

		for i2 := range l2 {
			if i1 >= len(sr1) {
				return i2 == l2-1 && sr2[i2] == '.'
			}
			if sr1[i1] != sr2[i2] && sr2[i2] != '_' {
				return false
			}
			i1++
		}

		return i1 >= l1
	}

	if doCompare(_s1, _s2) {
		return true
	}

	// URL-decoding s1 since Google Photos may return URL-encoded filenames
	if decoded, err := url.QueryUnescape(_s1); err == nil {
		return doCompare(decoded, _s2)
	}

	return false
}
