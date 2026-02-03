//go:build !darwin

package mediainfo

import "os"

func fileTimes(path string) (string, string, string, string, bool) {
	info, err := os.Stat(path)
	if err != nil {
		return "", "", "", "", false
	}
	mod := info.ModTime()
	createdUTC := mod.UTC().Format("2006-01-02 15:04:05 MST")
	createdLocal := mod.Local().Format("2006-01-02 15:04:05")
	modUTC := createdUTC
	modLocal := createdLocal
	return createdUTC, createdLocal, modUTC, modLocal, true
}
