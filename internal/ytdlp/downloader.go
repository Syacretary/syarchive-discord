package ytdlp

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

type Downloader struct {
	maxConcurrent int
}

type DownloadOptions struct {
	URL      string
	Format   string
	Audio    bool
	NoCookie bool
}

type VideoInfo struct {
	Title       string `json:"title"`
	Duration    int    `json:"duration"`
	Uploader    string `json:"uploader"`
	ViewCount   int    `json:"view_count"`
	LikeCount   int    `json:"like_count"`
	Description string `json:"description"`
	Thumbnail   string `json:"thumbnail"`
}

func NewDownloader() *Downloader {
	return &Downloader{
		maxConcurrent: 3,
	}
}

func (d *Downloader) DownloadVideo(opts DownloadOptions) (string, error) {
	args := []string{"--no-check-certificate"}

	if opts.NoCookie {
		args = append(args, "--no-cookies")
	}

	if opts.Audio {
		args = append(args, "-x", "--audio-format", "mp3")
	} else if opts.Format != "" {
		args = append(args, "-f", opts.Format)
	}

	// Output to temporary file
	args = append(args, "-o", "/tmp/%(title)s.%(ext)s", opts.URL)

	cmd := exec.Command("yt-dlp", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("download failed: %v, output: %s", err, string(output))
	}

	// Extract filename from output
	lines := strings.Split(string(output), "\n")
	for i := len(lines) - 1; i >= 0; i-- {
		if strings.Contains(lines[i], "[Merger] Merging formats into") {
			// Extract filename from merger line
			parts := strings.Split(lines[i], " ")
			if len(parts) > 5 {
				filename := strings.Trim(parts[5], "\"")
				return filename, nil
			}
		}
		
		// Also check for destination line
		if strings.Contains(lines[i], "[download] Destination:") {
			parts := strings.Split(lines[i], " ")
			if len(parts) > 3 {
				filename := strings.TrimSpace(parts[3])
				return filename, nil
			}
		}
	}

	return "", fmt.Errorf("could not determine output filename")
}

func (d *Downloader) GetInfo(url string) (*VideoInfo, error) {
	cmd := exec.Command("yt-dlp", "--dump-json", url)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to get video info: %v, output: %s", err, string(output))
	}

	var info VideoInfo
	if err := json.Unmarshal(output, &info); err != nil {
		return nil, fmt.Errorf("failed to parse video info: %v", err)
	}

	return &info, nil
}

func (d *Downloader) GetFormats(url string) (string, error) {
	cmd := exec.Command("yt-dlp", "-F", url)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to get available formats: %v, output: %s", err, string(output))
	}

	return string(output), nil
}

func (d *Downloader) DownloadWithFormat(url, format string) (string, error) {
	opts := DownloadOptions{
		URL:    url,
		Format: format,
	}
	
	return d.DownloadVideo(opts)
}

func (d *Downloader) DownloadAudio(url string) (string, error) {
	opts := DownloadOptions{
		URL:   url,
		Audio: true,
	}
	
	return d.DownloadVideo(opts)
}