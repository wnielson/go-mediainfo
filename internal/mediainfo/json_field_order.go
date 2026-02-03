package mediainfo

import "sort"

var jsonGeneralFieldOrder = map[string]int{
	"@type":                    0,
	"ID":                       1,
	"UniqueID":                 1,
	"VideoCount":               2,
	"AudioCount":               3,
	"TextCount":                4,
	"ImageCount":               5,
	"MenuCount":                6,
	"FileExtension":            7,
	"Format":                   8,
	"Format_Settings":          9,
	"Format_Version":           10,
	"Format_Profile":           11,
	"CodecID":                  12,
	"CodecID_Compatible":       13,
	"FileSize":                 14,
	"Duration":                 15,
	"OverallBitRate_Mode":      16,
	"OverallBitRate":           17,
	"FrameRate":                18,
	"FrameCount":               19,
	"StreamSize":               20,
	"HeaderSize":               21,
	"DataSize":                 22,
	"FooterSize":               23,
	"IsStreamable":             24,
	"File_Created_Date":        25,
	"File_Created_Date_Local":  26,
	"File_Modified_Date":       27,
	"File_Modified_Date_Local": 28,
	"Encoded_Application":      29,
	"Encoded_Library":          30,
	"extra":                    31,
}

var jsonVideoFieldOrder = map[string]int{
	"@type":                             0,
	"StreamOrder":                       1,
	"FirstPacketOrder":                  2,
	"ID":                                3,
	"MenuID":                            4,
	"UniqueID":                          5,
	"Format":                            6,
	"Format_Version":                    7,
	"Format_Profile":                    8,
	"Format_Level":                      9,
	"Format_Settings_CABAC":             10,
	"Format_Settings_RefFrames":         11,
	"Format_Settings_BVOP":              12,
	"Format_Settings_QPel":              12,
	"Format_Settings_GMC":               12,
	"Format_Settings_Matrix":            13,
	"Format_Settings_Matrix_Data":       14,
	"Format_Settings_GOP":               15,
	"CodecID":                           16,
	"Duration":                          17,
	"BitRate_Mode":                      18,
	"BitRate":                           19,
	"BitRate_Nominal":                   20,
	"BitRate_Maximum":                   21,
	"Width":                             22,
	"Height":                            23,
	"Sampled_Width":                     24,
	"Sampled_Height":                    25,
	"PixelAspectRatio":                  26,
	"DisplayAspectRatio":                27,
	"Rotation":                          28,
	"FrameRate_Mode":                    29,
	"FrameRate_Mode_Original":           30,
	"FrameRate":                         31,
	"FrameRate_Num":                     32,
	"FrameRate_Den":                     33,
	"FrameCount":                        34,
	"Standard":                          35,
	"ColorSpace":                        36,
	"ChromaSubsampling":                 37,
	"BitDepth":                          38,
	"ScanType":                          39,
	"Compression_Mode":                  40,
	"Delay":                             41,
	"Delay_Settings":                    42,
	"Delay_DropFrame":                   43,
	"Delay_Source":                      44,
	"Delay_Original":                    45,
	"Delay_Original_DropFrame":          46,
	"Delay_Original_Source":             47,
	"TimeCode_FirstFrame":               48,
	"TimeCode_Source":                   49,
	"Gop_OpenClosed":                    50,
	"Gop_OpenClosed_FirstFrame":         51,
	"StreamSize":                        52,
	"BufferSize":                        53,
	"Encoded_Library":                   54,
	"Encoded_Library_Name":              55,
	"Encoded_Library_Version":           56,
	"Encoded_Library_Settings":          57,
	"Default":                           58,
	"Forced":                            59,
	"colour_description_present":        60,
	"colour_description_present_Source": 61,
	"colour_range":                      62,
	"colour_range_Source":               63,
	"List_StreamKind":                   64,
	"List_StreamPos":                    65,
	"ServiceName":                       66,
	"ServiceProvider":                   67,
	"ServiceType":                       68,
	"extra":                             69,
}

var jsonAudioFieldOrder = map[string]int{
	"@type":                      0,
	"StreamOrder":                1,
	"FirstPacketOrder":           2,
	"ID":                         3,
	"MenuID":                     4,
	"UniqueID":                   5,
	"Format":                     6,
	"Format_Commercial_IfAny":    7,
	"Format_Settings_Endianness": 8,
	"Format_Version":             9,
	"Format_Settings_SBR":        10,
	"Format_AdditionalFeatures":  11,
	"MuxingMode":                 12,
	"CodecID":                    13,
	"Duration":                   14,
	"Source_Duration":            15,
	"Source_Duration_LastFrame":  16,
	"BitRate_Mode":               17,
	"BitRate":                    18,
	"BitRate_Maximum":            19,
	"Channels":                   20,
	"ChannelPositions":           21,
	"ChannelLayout":              22,
	"SamplesPerFrame":            23,
	"SamplingRate":               24,
	"SamplingCount":              25,
	"FrameRate":                  26,
	"FrameCount":                 27,
	"Source_FrameCount":          28,
	"Compression_Mode":           29,
	"Delay":                      30,
	"Delay_Source":               31,
	"Video_Delay":                32,
	"Encoded_Library":            33,
	"StreamSize":                 34,
	"Source_StreamSize":          35,
	"Default":                    36,
	"Forced":                     37,
	"ServiceKind":                38,
	"AlternateGroup":             39,
	"extra":                      40,
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
