#!/usr/bin/env bash
set -euo pipefail

need() {
  command -v "$1" >/dev/null 2>&1 || {
    echo "missing dependency: $1" >&2
    exit 1
  }
}

need ffmpeg

here="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
out="$here"

dur="${DUR:-4}"
rate="${RATE:-30000/1001}"
sr="${SR:-48000}"

tmp="$(mktemp -d)"
cleanup() { rm -rf "$tmp"; }
trap cleanup EXIT

echo "ffmpeg: $(ffmpeg -version | head -n1)"
echo "tmp: $tmp"

ffmpeg -hide_banner -loglevel error \
  -f lavfi -i "testsrc2=size=640x360:rate=${rate}:duration=${dur}" \
  -f lavfi -i "sine=frequency=1000:sample_rate=${sr}:duration=${dur}" \
  -map_metadata -1 -movflags +faststart \
  -c:v libx264 -preset veryfast -crf 28 -pix_fmt yuv420p -g 60 -keyint_min 60 -sc_threshold 0 \
  -c:a aac -b:a 96k \
  -shortest "$tmp/sample.mp4"

ffmpeg -hide_banner -loglevel error \
  -f lavfi -i "testsrc2=size=640x360:rate=${rate}:duration=${dur}" \
  -f lavfi -i "sine=frequency=1000:sample_rate=${sr}:duration=${dur}" \
  -map_metadata -1 \
  -c:v libx264 -preset veryfast -crf 28 -pix_fmt yuv420p -g 60 -keyint_min 60 -sc_threshold 0 \
  -c:a aac -b:a 96k \
  -shortest "$tmp/sample.mkv"

ffmpeg -hide_banner -loglevel error \
  -f lavfi -i "testsrc2=size=640x360:rate=${rate}:duration=${dur}" \
  -f lavfi -i "sine=frequency=1000:sample_rate=${sr}:duration=${dur}" \
  -map_metadata -1 \
  -c:v libx264 -preset veryfast -crf 28 -pix_fmt yuv420p -g 60 -keyint_min 60 -sc_threshold 0 \
  -c:a aac -b:a 96k \
  -f mpegts -mpegts_flags +resend_headers \
  -shortest "$tmp/sample.ts"

ffmpeg -hide_banner -loglevel error \
  -f lavfi -i "testsrc2=size=640x360:rate=${rate}:duration=${dur}" \
  -f lavfi -i "sine=frequency=1000:sample_rate=${sr}:duration=${dur}" \
  -map_metadata -1 \
  -c:v mpeg4 -q:v 5 -vtag XVID \
  -c:a libmp3lame -b:a 128k \
  -shortest "$tmp/sample.avi"

ffmpeg -hide_banner -loglevel error \
  -f lavfi -i "testsrc2=size=640x360:rate=${rate}:duration=${dur}" \
  -map_metadata -1 \
  -c:v mpeg2video -q:v 8 -g 12 -pix_fmt yuv420p \
  -f mpeg2video "$tmp/sample.mpg"

ffmpeg -hide_banner -loglevel error \
  -f lavfi -i "testsrc2=size=720x480:rate=${rate}:duration=${dur}" \
  -f lavfi -i "sine=frequency=1000:sample_rate=${sr}:duration=${dur}" \
  -map_metadata -1 \
  -c:v mpeg2video -b:v 2000k -maxrate 8000k -bufsize 1835k -g 15 -pix_fmt yuv420p \
  -c:a mp2 -b:a 192k \
  -f vob -shortest "$tmp/sample.vob"

ffmpeg -hide_banner -loglevel error \
  -f lavfi -i "testsrc2=size=720x480:rate=${rate}:duration=${dur}" \
  -f lavfi -i "sine=frequency=1000:sample_rate=${sr}:duration=${dur}" \
  -map_metadata -1 \
  -c:v mpeg2video -b:v 2000k -maxrate 8000k -bufsize 1835k -g 15 -pix_fmt yuv420p \
  -c:a ac3 -b:a 192k -ac 2 \
  -f vob -shortest "$tmp/sample_ac3.vob"

ffmpeg -hide_banner -loglevel error \
  -f lavfi -i "sine=frequency=1000:sample_rate=${sr}:duration=${dur}" \
  -map_metadata -1 \
  -c:a libmp3lame -b:a 128k \
  "$tmp/sample.mp3"

ffmpeg -hide_banner -loglevel error \
  -f lavfi -i "sine=frequency=1000:sample_rate=${sr}:duration=${dur}" \
  -map_metadata -1 \
  -c:a flac \
  "$tmp/sample.flac"

ffmpeg -hide_banner -loglevel error \
  -f lavfi -i "sine=frequency=1000:sample_rate=${sr}:duration=${dur}" \
  -map_metadata -1 \
  -c:a pcm_s16le \
  "$tmp/sample.wav"

ffmpeg -hide_banner -loglevel error \
  -f lavfi -i "sine=frequency=1000:sample_rate=${sr}:duration=${dur}" \
  -map_metadata -1 \
  -c:a libopus -b:a 96k \
  "$tmp/sample.ogg"

for f in sample.mp4 sample.mkv sample.ts sample.avi sample.mpg sample.vob sample_ac3.vob sample.mp3 sample.flac sample.wav sample.ogg; do
  mv -f "$tmp/$f" "$out/$f"
done

echo "sha256:"
shasum -a 256 "$out"/sample*.*
