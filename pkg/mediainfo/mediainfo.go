package mediainfo

import (
	"github.com/autobrr/go-mediainfo/internal/mediainfo"
)

// Types
type StreamKind = mediainfo.StreamKind
type Field = mediainfo.Field
type Stream = mediainfo.Stream
type Report = mediainfo.Report
type AnalyzeOptions = mediainfo.AnalyzeOptions

// Constants
const (
	StreamGeneral = mediainfo.StreamGeneral
	StreamVideo   = mediainfo.StreamVideo
	StreamAudio   = mediainfo.StreamAudio
	StreamText    = mediainfo.StreamText
	StreamImage   = mediainfo.StreamImage
	StreamMenu    = mediainfo.StreamMenu
)

// Functions
func AnalyzeFile(path string) (Report, error) {
	return mediainfo.AnalyzeFile(path)
}

func AnalyzeFileWithOptions(path string, opts AnalyzeOptions) (Report, error) {
	return mediainfo.AnalyzeFileWithOptions(path, opts)
}

func AnalyzeFilesWithOptions(paths []string, opts AnalyzeOptions) ([]Report, int, error) {
	return mediainfo.AnalyzeFilesWithOptions(paths, opts)
}

// Rendering
func RenderText(reports []Report) string {
	return mediainfo.RenderText(reports)
}

func RenderJSON(reports []Report) string {
	return mediainfo.RenderJSON(reports)
}

func RenderXML(reports []Report) string {
	return mediainfo.RenderXML(reports)
}

func RenderCSV(reports []Report) string {
	return mediainfo.RenderCSV(reports)
}

func RenderHTML(reports []Report) string {
	return mediainfo.RenderHTML(reports)
}

func RenderEBUCore(reports []Report) string {
	return mediainfo.RenderEBUCore(reports)
}

func RenderPBCore(reports []Report) string {
	return mediainfo.RenderPBCore(reports)
}

func RenderGraphSVG(reports []Report) string {
	return mediainfo.RenderGraphSVG(reports)
}

func RenderGraphDOT(reports []Report) string {
	return mediainfo.RenderGraphDOT(reports)
}

func InfoParameters() string {
	return mediainfo.InfoParameters()
}

func FormatVersion(version string) string {
	return mediainfo.FormatVersion(version)
}

func SetAppVersion(version string) {
	mediainfo.SetAppVersion(version)
}
