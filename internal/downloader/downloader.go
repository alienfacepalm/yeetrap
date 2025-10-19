package downloader

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"

	"github.com/AlienFacepalm/YeeTrap/internal/youtube"
)

// Downloader handles video downloads
type Downloader struct {
	outputDir  string
	quality    string
	concurrent int
}

// NewDownloader creates a new downloader
func NewDownloader(outputDir, quality string, concurrent int) *Downloader {
	return &Downloader{
		outputDir:  outputDir,
		quality:    quality,
		concurrent: concurrent,
	}
}

// DownloadVideos downloads multiple videos with concurrency control
func (d *Downloader) DownloadVideos(videos []youtube.Video) error {
	// Check if yt-dlp is available
	if err := d.checkYtDlp(); err != nil {
		return err
	}

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, d.concurrent)
	errors := make(chan error, len(videos))

	for i, video := range videos {
		wg.Add(1)
		go func(idx int, v youtube.Video) {
			defer wg.Done()

			semaphore <- struct{}{}        // Acquire
			defer func() { <-semaphore }() // Release

			fmt.Printf("[%d/%d] Downloading: %s\n", idx+1, len(videos), v.Title)
			if err := d.downloadVideo(v); err != nil {
				errors <- fmt.Errorf("failed to download %s: %w", v.Title, err)
			} else {
				fmt.Printf("[%d/%d] ✓ Completed: %s\n", idx+1, len(videos), v.Title)
			}
		}(i, video)
	}

	wg.Wait()
	close(errors)

	// Collect all errors
	var downloadErrors []error
	for err := range errors {
		downloadErrors = append(downloadErrors, err)
	}

	if len(downloadErrors) > 0 {
		fmt.Println("\nSome downloads failed:")
		for _, err := range downloadErrors {
			fmt.Printf("  - %v\n", err)
		}
		return fmt.Errorf("%d download(s) failed", len(downloadErrors))
	}

	return nil
}

// downloadVideo downloads a single video using yt-dlp
func (d *Downloader) downloadVideo(video youtube.Video) error {
	url := fmt.Sprintf("https://www.youtube.com/watch?v=%s", video.ID)
	
	// Sanitize filename
	filename := sanitizeFilename(video.Title)
	outputPath := filepath.Join(d.outputDir, filename+".%(ext)s")

	args := []string{
		"-f", d.getFormatString(),
		"-o", outputPath,
		"--no-playlist",
		"--write-description",
		"--write-info-json",
		"--write-thumbnail",
		url,
	}

	cmd := exec.Command("yt-dlp", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// checkYtDlp checks if yt-dlp is installed
func (d *Downloader) checkYtDlp() error {
	cmd := exec.Command("yt-dlp", "--version")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("yt-dlp is not installed or not in PATH. Please install it from https://github.com/yt-dlp/yt-dlp")
	}
	return nil
}

// getFormatString returns the yt-dlp format string based on quality setting
func (d *Downloader) getFormatString() string {
	switch d.quality {
	case "1080p":
		return "bestvideo[height<=1080]+bestaudio/best[height<=1080]"
	case "720p":
		return "bestvideo[height<=720]+bestaudio/best[height<=720]"
	case "480p":
		return "bestvideo[height<=480]+bestaudio/best[height<=480]"
	default:
		return "bestvideo+bestaudio/best"
	}
}

// sanitizeFilename removes invalid characters from filenames
func sanitizeFilename(filename string) string {
	invalid := []string{"/", "\\", ":", "*", "?", "\"", "<", ">", "|"}
	result := filename

	for _, char := range invalid {
		result = replaceAll(result, char, "_")
	}

	// Limit filename length
	if len(result) > 200 {
		result = result[:200]
	}

	return result
}

// replaceAll replaces all occurrences of old with new in s
func replaceAll(s, old, new string) string {
	result := ""
	for _, c := range s {
		if string(c) == old {
			result += new
		} else {
			result += string(c)
		}
	}
	return result
}


