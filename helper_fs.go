package main

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/rs/zerolog"
)

// buildDirCache
func buildDirCache(log zerolog.Logger, cacheMap *sync.Map, directory string) error {
	defer logDebugElapsed(log, time.Now(), fmt.Sprintf("buildDirCache for %s done in", directory))
	err := filepath.WalkDir(directory, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if (d.Type()&fs.ModeSymlink != 0 || d.IsDir()) && isValidGPhotosId(d.Name()) {
			cacheMap.Store(d.Name(), d)
		}
		return nil // Continue the walk
	})
	return err
}

// migrateYearMonth
func migrateYearMonth(log zerolog.Logger, downloadDir string, imageId string) error {
	imagIdDir := filepath.Join(downloadDir, imageId)
	infoImageIdDir, err := os.Stat(imagIdDir)
	if err == nil && infoImageIdDir.IsDir() {
		log.Debug().Msgf("migrating item to year month structure")
		defer logDebugElapsed(log, time.Now(), "migrateYearMonth done")
		entries, err := os.ReadDir(imagIdDir)
		if err != nil {
			return err
		}
		for _, imageFile := range entries {
			if !imageFile.Type().IsRegular() {
				return errors.New("file in imageId fodler is not a regular file")
			}
			imageFileInfo, err := imageFile.Info()
			if err != nil {
				return err
			}
			modTime := imageFileInfo.ModTime()
			year := modTime.Format("2006")
			month := modTime.Format("01")

			targetDirPath := filepath.Join(downloadDir, year, month)

			err = os.MkdirAll(targetDirPath, 0700)
			if err != nil {
				return err
			}
			err = os.Rename(filepath.Join(downloadDir, infoImageIdDir.Name(), imageFile.Name()), filepath.Join(targetDirPath, imageFile.Name()))
			if err != nil {
				return err
			}
			err = os.Symlink(filepath.Join(targetDirPath, imageFile.Name()), filepath.Join(targetDirPath, imageId))
			if err != nil {
				return err
			}
		}
		err = os.Remove(filepath.Join(downloadDir, infoImageIdDir.Name()))
		if err != nil {
			return err
		}
		return nil
	}
	if os.IsNotExist(err) {
		return nil
	}
	return err
}
