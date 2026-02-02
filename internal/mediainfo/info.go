package mediainfo

func InfoParameters() string {
	return "General\nComplete name\nFormat\nFile size\nDuration\nOverall bit rate mode\nOverall bit rate\n\nVideo\nFormat\nWidth\nHeight\nFrame rate\nBit rate mode\nBit rate\n\nAudio\nFormat\nChannel(s)\nChannel layout\nSampling rate\nBit rate mode\nBit rate\n\nText\nFormat\nDuration\n"
}

func InfoCanHandleUrls() string {
	return "No"
}
