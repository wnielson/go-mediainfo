package mediainfo

var generalFieldOrder = map[string]int{
	"Complete name":    0,
	"Format":           1,
	"File size":        2,
	"Duration":         3,
	"Overall bit rate": 4,
}

var streamFieldOrder = map[string]int{
	"ID":            0,
	"Format":        1,
	"Duration":      2,
	"Bit rate":      3,
	"Width":         4,
	"Height":        5,
	"Frame rate":    6,
	"Channel(s)":    7,
	"Sampling rate": 8,
}
