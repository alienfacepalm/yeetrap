package downloader

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"

	"github.com/AlienFacepalm/YeeTrap/internal/constants"
	"github.com/AlienFacepalm/YeeTrap/internal/errors"
	"github.com/AlienFacepalm/YeeTrap/internal/logger"
	"github.com/AlienFacepalm/YeeTrap/internal/progress"
	"github.com/AlienFacepalm/YeeTrap/internal/retry"
	"github.com/AlienFacepalm/YeeTrap/internal/validation"
	"github.com/AlienFacepalm/YeeTrap/internal/youtube"
)

// Downloader handles video downloads
type Downloader struct {
	outputDir  string
	quality    string
	concurrent int
	progress   *progress.ProgressTracker
}

// NewDownloader creates a new downloader
func NewDownloader(outputDir, quality string, concurrent int) (*Downloader, error) {
	// Validate inputs
	if err := validation.ValidateOutputDir(outputDir); err != nil {
		return nil, err
	}
	
	if err := validation.ValidateQuality(quality); err != nil {
		return nil, err
	}
	
	if err := validation.ValidateConcurrency(concurrent); err != nil {
		return nil, err
	}
	
	logger.Info("Creating downloader with output: %s, quality: %s, concurrent: %d", outputDir, quality, concurrent)
	
	return &Downloader{
		outputDir:  outputDir,
		quality:    quality,
		concurrent: concurrent,
	}, nil
}

// DownloadVideos downloads multiple videos with concurrency control
func (d *Downloader) DownloadVideos(videos []youtube.Video) error {
	logger.Info("Starting download of %d videos", len(videos))
	
	// Check if yt-dlp is available
	if err := d.checkYtDlp(); err != nil {
		return err
	}

	// Create output directory
	if err := os.MkdirAll(d.outputDir, 0755); err != nil {
		return errors.WrapFile(err, "failed to create output directory")
	}

	// Initialize progress tracker
	d.progress = progress.NewProgressTracker(len(videos))
	d.progress.AddCallback(progress.DefaultProgressCallback)
	defer d.progress.Stop()

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, d.concurrent)
	downloadErrors := make(chan error, len(videos))

	for i, video := range videos {
		wg.Add(1)
		go func(idx int, v youtube.Video) {
			defer wg.Done()

			semaphore <- struct{}{}        // Acquire
			defer func() { <-semaphore }() // Release

			// Update progress
			d.progress.UpdateProgress(idx, 0, fmt.Sprintf("Downloading: %s", v.Title))
			
			// Download with retry logic
			err := retry.RetryDownloadOperation(func() error {
				return d.downloadVideo(v)
			})
			
			if err != nil {
				logger.Error("Failed to download %s: %v", v.Title, err)
				d.progress.IncrementFailed(v.Title)
				downloadErrors <- errors.WrapExternal(err, fmt.Sprintf("failed to download %s", v.Title))
			} else {
				logger.Info("Successfully downloaded: %s", v.Title)
				d.progress.IncrementCompleted(v.Title)
			}
		}(i, video)
	}

	wg.Wait()
	close(downloadErrors)

	// Collect all errors
	var downloadErrorsList []error
	for err := range downloadErrors {
		downloadErrorsList = append(downloadErrorsList, err)
	}

	if len(downloadErrorsList) > 0 {
		logger.Warn("%d downloads failed", len(downloadErrorsList))
		fmt.Println("\nSome downloads failed:")
		for _, err := range downloadErrorsList {
			fmt.Printf("  - %v\n", err)
		}
		return errors.NewExternalError(fmt.Sprintf("%d download(s) failed", len(downloadErrorsList)))
	}

	logger.Info("All downloads completed successfully")
	return nil
}

// downloadVideo downloads a single video using yt-dlp
func (d *Downloader) downloadVideo(video youtube.Video) error {
	logger.Debug("Downloading video: %s (%s)", video.Title, video.ID)
	
	url := fmt.Sprintf("https://www.youtube.com/watch?v=%s", video.ID)
	
	// Sanitize filename
	filename := validation.SanitizeFilename(video.Title)
	outputPath := filepath.Join(d.outputDir, filename+".%(ext)s")

	args := []string{
		"-f", d.getFormatString(),
		"-o", outputPath,
		"--no-playlist",
		"--write-description",
		"--write-info-json",
		"--write-thumbnail",
		"--no-warnings",
		url,
	}

	cmd := exec.Command("yt-dlp", args...)
	// Capture output for logging instead of printing to stdout
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return errors.WrapExternal(err, fmt.Sprintf("yt-dlp failed for video %s", video.ID))
	}
	
	return nil
}

// checkYtDlp checks if yt-dlp is installed
func (d *Downloader) checkYtDlp() error {
	logger.Debug("Checking if yt-dlp is available")
	
	cmd := exec.Command("yt-dlp", "--version")
	if err := cmd.Run(); err != nil {
		return errors.NewExternalError("yt-dlp is not installed or not in PATH").
			WithDetails("Please install it from https://github.com/yt-dlp/yt-dlp")
	}
	
	logger.Debug("yt-dlp is available")
	return nil
}

// getFormatString returns the yt-dlp format string based on quality setting
func (d *Downloader) getFormatString() string {
	switch d.quality {
	case constants.Quality1080p:
		return "bestvideo[height<=1080]+bestaudio/best[height<=1080]"
	case constants.Quality720p:
		return "bestvideo[height<=720]+bestaudio/best[height<=720]"
	case constants.Quality480p:
		return "bestvideo[height<=480]+bestaudio/best[height<=480]"
	default:
		return "bestvideo+bestaudio/best"
	}
}



