package mediainfo

import (
	"encoding/json"
	"path/filepath"
	"testing"
)

func TestMpegVideoJSONFields(t *testing.T) {
	path := filepath.Join("..", "..", "samples", "sample.mpg")
	report, err := AnalyzeFile(path)
	if err != nil {
		t.Fatalf("analyze sample: %v", err)
	}

	output := RenderJSON([]Report{report})
	var root map[string]any
	if err := json.Unmarshal([]byte(output), &root); err != nil {
		t.Fatalf("parse json: %v", err)
	}

	media, ok := root["media"].(map[string]any)
	if !ok {
		t.Fatalf("missing media object")
	}
	tracks, ok := media["track"].([]any)
	if !ok {
		t.Fatalf("missing track list")
	}

	var general map[string]any
	var video map[string]any
	for _, item := range tracks {
		track, ok := item.(map[string]any)
		if !ok {
			continue
		}
		switch track["@type"] {
		case "General":
			general = track
		case "Video":
			video = track
		}
	}
	if general == nil || video == nil {
		t.Fatalf("missing general or video track")
	}

	if got, _ := general["OverallBitRate"].(string); got != "328915" {
		t.Fatalf("unexpected overall bitrate: %v", got)
	}
	if got, _ := general["FrameCount"].(string); got != "400" {
		t.Fatalf("unexpected general framecount: %v", got)
	}
	if got, _ := general["StreamSize"].(string); got != "0" {
		t.Fatalf("unexpected general stream size: %v", got)
	}

	if got, _ := video["BitRate"].(string); got != "328915" {
		t.Fatalf("unexpected video bitrate: %v", got)
	}
	if got, _ := video["StreamSize"].(string); got != "548754" {
		t.Fatalf("unexpected video stream size: %v", got)
	}
	if got, _ := video["Delay_Settings"].(string); got != "drop_frame_flag=0 / closed_gop=1 / broken_link=0" {
		t.Fatalf("unexpected delay settings: %v", got)
	}
	if got, _ := video["BufferSize"].(string); got != "6144" {
		t.Fatalf("unexpected buffer size: %v", got)
	}
	if extra, ok := video["extra"].(map[string]any); !ok {
		t.Fatalf("missing extra")
	} else if got, _ := extra["intra_dc_precision"].(string); got != "8" {
		t.Fatalf("unexpected intra_dc_precision: %v", got)
	}
}
