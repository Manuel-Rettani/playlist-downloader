# yt-playlist-downloader

### Playlist downloader written in GO

playlist-downloader downloads a playlist given an id and uploads a zip with its songs to an s3 bucket, so that it can be downloaded from any device.

### Usage

1. Install python
2. Install yt-dlp (pip install yt-dlp)

    ```bash
    pip install yt-dlp
    ```
3. brew install ffmpeg
4. Check installation with

    ```bash
    which ffmpeg
    which ffprobe
    ```

If receiving errors from yt-dlp, like 403, try updating it

```bash
python3 -m pip install -U yt-dlp
```