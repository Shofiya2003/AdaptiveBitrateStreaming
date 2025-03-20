#!/bin/sh

INPUT_FILE=$1
OUTPUT_DIR=$2

mkdir -p "$OUTPUT_DIR"

ffmpeg -i "$INPUT_FILE" \
    -filter:v:0 "scale=-2:1080" -c:v:0 libx264 -b:v:0 5000k -preset fast -c:a:0 aac -b:a:0 192k -hls_time 5 -hls_playlist_type vod -hls_segment_filename "$OUTPUT_DIR/1080p_%03d.ts" -hls_flags independent_segments -f hls "$OUTPUT_DIR/1080p.m3u8" \
    -filter:v:1 "scale=-2:720"  -c:v:1 libx264 -b:v:1 2500k -preset fast -c:a:1 aac -b:a:1 128k -hls_time 5 -hls_playlist_type vod -hls_segment_filename "$OUTPUT_DIR/720p_%03d.ts" -hls_flags independent_segments -f hls "$OUTPUT_DIR/720p.m3u8" \
    -filter:v:2 "scale=-2:360"  -c:v:2 libx264 -b:v:2 800k  -preset fast -c:a:2 aac -b:a:2 96k  -hls_time 5 -hls_playlist_type vod -hls_segment_filename "$OUTPUT_DIR/360p_%03d.ts" -hls_flags independent_segments -f hls "$OUTPUT_DIR/360p.m3u8" \
    -filter:v:3 "scale=-2:240"  -c:v:3 libx264 -b:v:3 400k  -preset fast -c:a:3 aac -b:a:3 64k  -hls_time 5 -hls_playlist_type vod -hls_segment_filename "$OUTPUT_DIR/240p_%03d.ts" -hls_flags independent_segments -f hls "$OUTPUT_DIR/240p.m3u8" \

echo "#EXTM3U" > "$OUTPUT_DIR/index.m3u8"
echo "#EXT-X-STREAM-INF:BANDWIDTH=5000000,RESOLUTION=1920x1080" >> "$OUTPUT_DIR/index.m3u8"
echo "1080p.m3u8" >> "$OUTPUT_DIR/index.m3u8"
echo "#EXT-X-STREAM-INF:BANDWIDTH=2500000,RESOLUTION=1280x720" >> "$OUTPUT_DIR/index.m3u8"
echo "720p.m3u8" >> "$OUTPUT_DIR/index.m3u8"
echo "#EXT-X-STREAM-INF:BANDWIDTH=800000,RESOLUTION=640x360" >> "$OUTPUT_DIR/index.m3u8"
echo "360p.m3u8" >> "$OUTPUT_DIR/index.m3u8"
echo "#EXT-X-STREAM-INF:BANDWIDTH=400000,RESOLUTION=426x240" >> "$OUTPUT_DIR/index.m3u8"
echo "240p.m3u8" >> "$OUTPUT_DIR/index.m3u8"

echo "✅ Transcoding Complete! HLS files are in: $OUTPUT_DIR"
