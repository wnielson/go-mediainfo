package mediainfo

import (
	"fmt"
	"io"
	"os"
)

func AnalyzeFile(path string) (Report, error) {
	stat, err := os.Stat(path)
	if err != nil {
		return Report{}, err
	}

	header := make([]byte, maxSniffBytes)
	file, err := os.Open(path)
	if err != nil {
		return Report{}, err
	}
	defer file.Close()

	n, _ := io.ReadFull(file, header)
	header = header[:n]

	format := DetectFormat(header, path)

	general := Stream{Kind: StreamGeneral}
	general.Fields = append(general.Fields,
		Field{Name: "Complete name", Value: path},
		Field{Name: "Format", Value: format},
		Field{Name: "File size", Value: formatBytes(stat.Size())},
	)

	return Report{
		Ref:     path,
		General: general,
	}, nil
}

func AnalyzeFiles(paths []string) ([]Report, int, error) {
	reports := make([]Report, 0, len(paths))
	for _, path := range paths {
		report, err := AnalyzeFile(path)
		if err != nil {
			return nil, 0, fmt.Errorf("%s: %w", path, err)
		}
		reports = append(reports, report)
	}
	return reports, len(reports), nil
}
