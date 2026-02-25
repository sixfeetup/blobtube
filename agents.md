# Agents

## Video Conversion

All videos in this repository should be created as AV1 format for optimal compression and quality.

### Converting to AV1

To convert a GIF or other video format to AV1, use ffmpeg with the libaom-av1 codec:

```bash
ffmpeg -i input.gif -c:v libaom-av1 -crf 30 -b:v 0 output.mp4
```

### Scaling Videos

To create a scaled version (e.g., 4x4 pixels):

```bash
ffmpeg -i input.gif -vf "scale=4:4:flags=lanczos" -c:v libaom-av1 -crf 30 -b:v 0 output_4x4.mp4
```

### Parameters Explained

- `-c:v libaom-av1`: Use the AV1 codec
- `-crf 30`: Constant Rate Factor (quality), lower = better quality (0-63)
- `-b:v 0`: Use constant quality mode
- `-vf "scale=W:H:flags=lanczos"`: Scale video to width W and height H using Lanczos filter
