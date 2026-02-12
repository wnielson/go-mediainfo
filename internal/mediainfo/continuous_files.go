package mediainfo

import (
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

type continuousFileSet struct {
	LastPath  string
	TotalSize int64
	LastSize  int64
	Count     int
}

func detectContinuousFileSet(path string) (continuousFileSet, bool) {
	dir := filepath.Dir(path)
	base := filepath.Base(path)
	ext := filepath.Ext(base)
	name := strings.TrimSuffix(base, ext)

	// trailing digits
	i := len(name)
	for i > 0 {
		c := name[i-1]
		if c < '0' || c > '9' {
			break
		}
		i--
	}
	if i == len(name) {
		return continuousFileSet{}, false
	}
	prefix := name[:i]
	digits := name[i:]
	width := len(digits)
	start, err := strconv.Atoi(digits)
	if err != nil {
		return continuousFileSet{}, false
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return continuousFileSet{}, false
	}

	type numberedFile struct {
		index int
		path  string
		size  int64
	}
	matches := make([]numberedFile, 0, 32)
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		filename := entry.Name()
		if filepath.Ext(filename) != ext {
			continue
		}
		stem := strings.TrimSuffix(filename, ext)
		if !strings.HasPrefix(stem, prefix) {
			continue
		}
		suffix := stem[len(prefix):]
		if len(suffix) != width || suffix == "" {
			continue
		}
		validDigits := true
		for i := 0; i < len(suffix); i++ {
			if suffix[i] < '0' || suffix[i] > '9' {
				validDigits = false
				break
			}
		}
		if !validDigits {
			continue
		}
		index, err := strconv.Atoi(suffix)
		if err != nil || index < start {
			continue
		}
		info, err := entry.Info()
		if err != nil {
			continue
		}
		matches = append(matches, numberedFile{
			index: index,
			path:  filepath.Join(dir, filename),
			size:  info.Size(),
		})
	}
	if len(matches) < 2 {
		return continuousFileSet{}, false
	}
	sort.Slice(matches, func(i, j int) bool {
		if matches[i].index != matches[j].index {
			return matches[i].index < matches[j].index
		}
		return matches[i].path < matches[j].path
	})
	if matches[0].index != start {
		return continuousFileSet{}, false
	}
	if matches[0].path == path && len(matches) == 1 {
		return continuousFileSet{}, false
	}

	var total int64
	for _, file := range matches {
		total += file.size
	}
	last := matches[len(matches)-1]
	if last.path == path {
		return continuousFileSet{}, false
	}
	return continuousFileSet{
		LastPath:  last.path,
		TotalSize: total,
		LastSize:  last.size,
		Count:     len(matches),
	}, true
}
