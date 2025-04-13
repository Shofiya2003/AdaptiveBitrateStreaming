#!/bin/bash

# Input parameters
INPUT_FILE="$1"
WORK_DIR="$2"
OUTPUT_DIR="$3"

mkdir -p "$WORK_DIR" "$OUTPUT_DIR"

# Define resolutions and bitrates
declare -A RESOLUTIONS=(
    ["1080p"]="1920x1080:5000k"
    ["720p"]="1280x720:2500k"
    ["480p"]="854x480:1500k"
    ["360p"]="640x360:800k"
    ["240p"]="426x240:400k"
)

# Step 1: Calculate VMAF for each resolution
declare -A VMAF_SCORES
for res in "${!RESOLUTIONS[@]}"; do
    IFS=":" read -r scale bitrate <<< "${RESOLUTIONS[$res]}"
    
    # Transcode to current resolution (without VMAF for now)
    ffmpeg -y -i "$INPUT_FILE" -vf "scale=${scale}" -c:v mpeg4 -b:v "$bitrate" "$WORK_DIR/$res.mp4"

    # Compute VMAF score using the original video
    VMAF_SCORE=$(ffmpeg -i "$INPUT_FILE" -i "$WORK_DIR/$res.mp4" \
        -lavfi "[0:v]scale=1920x1080[ref];[1:v]scale=1920x1080[dist];[ref][dist]libvmaf=log_path=$WORK_DIR/vmaf_${res}.json:log_fmt=json" \
        -f null - 2>&1 | grep -oP 'VMAF score: \K[\d.]+')

    VMAF_SCORES["$res"]=$VMAF_SCORE
    echo "🎯 VMAF Score for $res: $VMAF_SCORE"
done

# Step 2: Filter based on VMAF scores (threshold = 75)
declare -A SELECTED_RES

for res in "${!VMAF_SCORES[@]}"; do
    if (( $(echo "${VMAF_SCORES[$res]} > 90" | bc -l) )); then
        SELECTED_RES["$res"]=1
    fi
done

# Step 3: Transcode and segment only selected resolutions
for res in "${!SELECTED_RES[@]}"; do
    IFS=":" read -r scale bitrate <<< "${RESOLUTIONS[$res]}"
    
    # Transcode the selected resolution with mpeg4 codec
    ffmpeg -y -i "$INPUT_FILE" -vf "scale=${scale},setsar=1" -c:v mpeg4 -b:v "$bitrate" "$WORK_DIR/$res.mp4"
    
    # Segment to HLS
    ffmpeg -y -i "$WORK_DIR/$res.mp4" -c:v copy -c:a copy \
        -hls_time 5 -hls_flags independent_segments \
        -hls_segment_filename "$OUTPUT_DIR/${res}_%03d.ts" \
        -f hls "$OUTPUT_DIR/${res}.m3u8"
done

# Step 4: Create master playlist
echo "#EXTM3U" > "$OUTPUT_DIR/index.m3u8"
for res in "${!SELECTED_RES[@]}"; do
    case $res in
        1080p) BANDWIDTH=5000000; RESOLUTION=1920x1080 ;;
        720p) BANDWIDTH=2500000; RESOLUTION=1280x720 ;;
        480p) BANDWIDTH=1500000; RESOLUTION=854x480 ;;
        360p) BANDWIDTH=800000;  RESOLUTION=640x360 ;;
        240p) BANDWIDTH=400000;  RESOLUTION=426x240 ;;
    esac
    echo "#EXT-X-STREAM-INF:BANDWIDTH=$BANDWIDTH,RESOLUTION=$RESOLUTION" >> "$OUTPUT_DIR/index.m3u8"
    echo "${res}.m3u8" >> "$OUTPUT_DIR/index.m3u8"
done

# 🖨️ Print selected resolutions
echo "🧩 Selected Resolutions to Transcode:"
for res in "${!SELECTED_RES[@]}"; do
    echo "- $res"
done

echo "✅ Done! VMAF-based Resolution Ladder created. HLS outputs in $OUTPUT_DIR"

rm -rf "$WORK_DIR"