# TODO

## Parity (1:1 MediaInfo CLI)
- Goal: 1:1 MediaInfo CLI parity (fields + values + ordering)
- Expand format coverage + field parity across CLI outputs
- Sample parity complete: MP4/MKV/TS/AVI/MPEG-PS (VOB)/MPEG Video (MPG) (text output)
- Next target: Audio in MPEG-PS (AC-3/AAC)
- MPEG-PS: match video bitrate/stream size when GOP-only duration (DVD menu/still frames), AC-3 extended metadata (dialnorm/compr), Menu stream output

## Post-parity
- Investigate MediaInfo issue #760: DVD IFO language/runtime regression (23.07 vs 23.06). https://github.com/MediaArea/MediaInfo/issues/760
