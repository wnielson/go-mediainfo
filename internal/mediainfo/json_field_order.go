package mediainfo

import "sort"

var jsonGeneralFieldOrder = map[string]int{
	"@type":                    0,
	"VideoCount":               1,
	"AudioCount":               2,
	"TextCount":                3,
	"ImageCount":               4,
	"MenuCount":                5,
	"FileExtension":            6,
	"Format":                   7,
	"Format_Profile":           8,
	"CodecID":                  9,
	"CodecID_Compatible":       10,
	"FileSize":                 11,
	"Duration":                 12,
	"OverallBitRate_Mode":      13,
	"OverallBitRate":           14,
	"FrameRate":                15,
	"FrameCount":               16,
	"StreamSize":               17,
	"HeaderSize":               18,
	"DataSize":                 19,
	"FooterSize":               20,
	"IsStreamable":             21,
	"File_Created_Date":        22,
	"File_Created_Date_Local":  23,
	"File_Modified_Date":       24,
	"File_Modified_Date_Local": 25,
	"Encoded_Application":      26,
}

var jsonVideoFieldOrder = map[string]int{
	"@type":                     0,
	"StreamOrder":               1,
	"ID":                        2,
	"Format":                    3,
	"Format_Profile":            4,
	"Format_Level":              5,
	"Format_Settings_CABAC":     6,
	"Format_Settings_RefFrames": 7,
	"CodecID":                   8,
	"Duration":                  9,
	"BitRate":                   10,
	"BitRate_Nominal":           11,
	"BitRate_Maximum":           12,
	"Width":                     13,
	"Height":                    14,
	"Sampled_Width":             15,
	"Sampled_Height":            16,
	"PixelAspectRatio":          17,
	"DisplayAspectRatio":        18,
	"Rotation":                  19,
	"FrameRate_Mode":            20,
	"FrameRate_Mode_Original":   21,
	"FrameRate":                 22,
	"FrameRate_Num":             23,
	"FrameRate_Den":             24,
	"FrameCount":                25,
	"ChromaSubsampling":         26,
	"BitDepth":                  27,
	"ScanType":                  28,
	"StreamSize":                29,
	"Encoded_Library":           30,
	"Encoded_Library_Name":      31,
	"Encoded_Library_Version":   32,
	"Encoded_Library_Settings":  33,
	"extra":                     34,
}

var jsonAudioFieldOrder = map[string]int{
	"@type":                     0,
	"StreamOrder":               1,
	"ID":                        2,
	"Format":                    3,
	"Format_Settings_SBR":       4,
	"Format_AdditionalFeatures": 5,
	"CodecID":                   6,
	"Duration":                  7,
	"Source_Duration":           8,
	"Source_Duration_LastFrame": 9,
	"BitRate_Mode":              10,
	"BitRate":                   11,
	"BitRate_Maximum":           12,
	"Channels":                  13,
	"ChannelPositions":          14,
	"ChannelLayout":             15,
	"SamplesPerFrame":           16,
	"SamplingRate":              17,
	"SamplingCount":             18,
	"FrameRate":                 19,
	"FrameCount":                20,
	"Source_FrameCount":         21,
	"Compression_Mode":          22,
	"StreamSize":                23,
	"Source_StreamSize":         24,
	"Default":                   25,
	"AlternateGroup":            26,
	"extra":                     27,
}

func sortJSONFields(kind StreamKind, fields []jsonKV) []jsonKV {
	order := jsonVideoFieldOrder
	switch kind {
	case StreamGeneral:
		order = jsonGeneralFieldOrder
	case StreamAudio:
		order = jsonAudioFieldOrder
	case StreamVideo:
		order = jsonVideoFieldOrder
	}
	positions := map[string]int{}
	for i, field := range fields {
		positions[field.Key] = i
	}
	sort.SliceStable(fields, func(i, j int) bool {
		ai, aok := order[fields[i].Key]
		aj, bok := order[fields[j].Key]
		switch {
		case aok && bok:
			return ai < aj
		case aok:
			return true
		case bok:
			return false
		default:
			return positions[fields[i].Key] < positions[fields[j].Key]
		}
	})
	return fields
}
