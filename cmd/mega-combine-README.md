# mega-combine 💕

A super cute command-line tool for selecting and combining multiple video files into a single ProRes file, optimized for DaVinci Resolve on iPad! ✨

The command uses different ffmpeg settings depending on the mode - we're so flexible! ✨

### Default Mode (GPU-Accelerated) 🎨

**Video Encoding:**
- **Codec**: `hevc_nvenc` (NVIDIA GPU-accelerated H.265/HEVC) - so fast! 💪
- **Preset**: `p6` (high quality, balanced speed)
- **Tune**: `hq` (high quality tuning)
- **Rate Control**: `vbr_hq` (high quality variable bitrate)
- **Quality**: `cq 22` (constant quality mode - so efficient! ✨)
- **Pixel Format**: `p010le` (10-bit)
- **Profile**: `main10` (H.265 Main 10 profile for 10-bit support)

**Audio Encoding:**
- **Codec**: `aac` (AAC compression)
- **Bitrate**: `160 kbps`
- **Sample Rate**: `48000 Hz`
- **Channels**: `2` (stereo)

**Output**: `.mp4` files with `+faststart` flag (much smaller than ProRes and ready for web streaming! ✨)

### `--waytoobig` Mode (ProRes) 💖

Use this flag when you need maximum quality for DaVinci Resolve on iPad and don't mind huge file sizes!

**Video Encoding:**
- **Codec**: `prores_ks` (Apple ProRes)
- **Profile**: `1` (ProRes LT - Light)
- **Pixel Format**: `yuv422p10le` (10-bit 4:2:2)
- **Threads**: `0` (automatic, uses all available CPU cores)

**Audio Encoding:**
- **Codec**: `pcm_s16le` (PCM 16-bit little-endian, uncompressed)
- **Sample Rate**: `48000 Hz`
- **Channels**: `2` (stereo)

**Output**: `.mov` files (way too big but no re-encoding needed in DaVinci Resolve!)

## Overview 🎀

`mega-combine` provides an interactive TUI (Terminal User Interface) to select video files from the current directory, then combines them into a single video file - so organized! 💖 By default, it uses GPU-accelerated H.265 encoding for efficient, high-quality output. Use the `--waytoobig` flag for ProRes encoding when you need maximum quality for DaVinci Resolve on iPad (but don't mind huge file sizes!). We're so flexible! ✨

## Usage 💅

```bash
# Interactive mode - select files and combine (default: GPU-accelerated H.265) ✨
marcli mega-combine

# Use ProRes encoding (way too big but high quality for DaVinci Resolve)
marcli mega-combine --waytoobig

# Preview the ffmpeg command that would be run
marcli mega-combine --test

# Specify custom output filename
marcli mega-combine --out myvideo

# Combine options - we're so flexible! 💕
marcli mega-combine --test --out myvideo
marcli mega-combine --waytoobig --out myvideo.mov
```

## Features 🎀

- **Interactive file selection**: Browse and multi-select video files ordered by modification time - so organized! 💖
- **Automatic file extension**: If you don't specify an extension, `.mp4` is added by default (or `.mov` with `--waytoobig`) - we're so helpful! ✨
- **Preview mode**: Use `--test` to see the exact ffmpeg command before running - safety first! 💅
- **Robust concatenation**: Uses timestamp normalization to handle variable frame rates and mismatched start times - so reliable! 🎨

### Concatenation Method
The command uses a robust concatenation approach with timestamp normalization:

1. **Timestamp Normalization**: Each input stream is normalized using `setpts=PTS-STARTPTS` for video and `asetpts=PTS-STARTPTS` for audio
   - This handles variable frame rates (VFR) safely
   - Resolves mismatched start times between files
   - Ensures smooth concatenation without gaps or sync issues

2. **Filter Complex**: Uses `concat` filter to combine normalized streams
   - Format: `[v0][a0][v1][a1]...concat=n=N:v=1:a=1[outv][outa]`
   - Where N is the number of input files

3. **Stream Mapping**: Maps the concatenated video and audio streams to output

### Why These Settings? 💕

**Default Mode (NVENC H.265):**
- **GPU Acceleration**: Uses your NVIDIA GPU (like the RTX 5090!) for super fast encoding - so efficient! 💪
- **H.265/HEVC**: Modern codec with excellent compression - much smaller files than ProRes! ✨
- **Constant Quality (CQ 22)**: Uses quality-based encoding instead of fixed bitrate - automatically adjusts bitrate to maintain quality while keeping files small - perfect balance! 🎨
- **10-bit Color**: `p010le` pixel format with Main 10 profile for better color depth - so fancy! 💅
- **Fast Start**: `+faststart` flag allows playback to begin before download completes - perfect for web streaming! 🎀
- **AAC Audio**: Compressed but high-quality audio (160 kbps) that's much smaller than PCM - so smart! 💖

**`--waytoobig` Mode (ProRes):**
- **ProRes LT**: Provides excellent quality but creates huge files - way too big! 💅 DaVinci Resolve on iPad natively supports ProRes, so no re-encoding is needed on import. We're so thoughtful! 💖
- **10-bit 4:2:2**: Maintains color depth and chroma subsampling suitable for professional editing - zero compromises! 🎀
- **PCM Audio**: Uncompressed audio ensures no quality loss - maximum quality! ✨
- **Use when**: You need maximum quality and don't mind waiting for huge file transfers! 💕

**Both Modes:**
- **Timestamp Normalization**: The robust concatenation method ensures smooth playback and editing, even when source files have different frame rates or start times - so reliable! 🎨

## Workflow 🎀

1. Navigate to the directory containing your video files - so organized! 💖
2. Run `marcli mega-combine` - let's go! ✨
3. Use arrow keys to navigate, Space to select/deselect files - so intuitive! 💕
4. Press Enter to confirm and start the combination process - here we go! 🎨
5. Watch the ffmpeg progress in real-time - so satisfying! 💅
6. Import the resulting `.mov` file into DaVinci Resolve on iPad - done! 🎀

## Supported Video Formats 💖

The command automatically detects and lists the following video file extensions - we're so flexible! ✨
- `.mp4`, `.avi`, `.mov`, `.mkv`, `.webm`, `.flv`, `.wmv`, `.m4v`, `.mpg`, `.mpeg`, `.3gp`, `.ogv`

Files are sorted by modification time (oldest first) to help maintain chronological order - so organized! 💕

## Tips 💅

- Use `--test` first to verify the command before running it on large files - safety first! ✨
- The combination process can take a while depending on file sizes and system performance - be patient, it's worth it! 💖
- You can press 'q' during ffmpeg encoding to quit (though this may leave an incomplete file) - we're so flexible! 🎀
- The output file will be created in the current working directory - so convenient! 💕

## Example Output

When running without `--test`, you'll see real-time ffmpeg progress:

```
Running ffmpeg to combine 3 video file(s) into output.mov...
Press 'q' during encoding to quit.

[ffmpeg output with frame counts, FPS, bitrate, etc.]

Video files successfully combined into output.mov
```

