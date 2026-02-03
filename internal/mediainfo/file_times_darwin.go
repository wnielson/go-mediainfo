//go:build darwin

package mediainfo

import (
	"os"
	"syscall"
	"time"
)

func fileTimes(path string) (string, string, string, string, bool) {
	info, err := os.Stat(path)
	if err != nil {
		return "", "", "", "", false
	}
	stat, ok := info.Sys().(*syscall.Stat_t)
	if !ok {
		return "", "", "", "", false
	}
	created := time.Unix(stat.Birthtimespec.Sec, stat.Birthtimespec.Nsec)
	mod := time.Unix(stat.Mtimespec.Sec, stat.Mtimespec.Nsec)
	createdUTC := created.UTC().Format("2006-01-02 15:04:05 MST")
	createdLocal := created.Local().Format("2006-01-02 15:04:05")
	modUTC := mod.UTC().Format("2006-01-02 15:04:05 MST")
	modLocal := mod.Local().Format("2006-01-02 15:04:05")
	return createdUTC, createdLocal, modUTC, modLocal, true
}
