package mediainfo

func appendChannelFields(fields []Field, channels uint64) []Field {
	if channels == 0 {
		return fields
	}
	fields = append(fields, Field{Name: "Channel(s)", Value: formatChannels(channels)})
	if layout := channelLayout(channels); layout != "" {
		fields = append(fields, Field{Name: "Channel layout", Value: layout})
	}
	return fields
}

func appendSampleRateField(fields []Field, sampleRate float64) []Field {
	if sampleRate <= 0 {
		return fields
	}
	return append(fields, Field{Name: "Sampling rate", Value: formatSampleRate(sampleRate)})
}
