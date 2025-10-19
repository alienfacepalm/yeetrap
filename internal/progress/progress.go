package progress

import (
	"fmt"
	"sync"
	"time"
)

// Progress represents the progress of an operation
type Progress struct {
	Total     int
	Completed int
	Failed    int
	Current   string
	StartTime time.Time
	mu        sync.RWMutex
}

// ProgressCallback is a function type for progress updates
type ProgressCallback func(p *Progress)

// ProgressTracker manages progress tracking for operations
type ProgressTracker struct {
	progress   *Progress
	callbacks  []ProgressCallback
	mu         sync.RWMutex
	updateChan chan *Progress
	done       chan bool
}

// NewProgressTracker creates a new progress tracker
func NewProgressTracker(total int) *ProgressTracker {
	pt := &ProgressTracker{
		progress: &Progress{
			Total:     total,
			Completed: 0,
			Failed:    0,
			StartTime: time.Now(),
		},
		callbacks:  make([]ProgressCallback, 0),
		updateChan: make(chan *Progress, 10),
		done:       make(chan bool),
	}
	
	// Start the progress updater
	go pt.updateLoop()
	
	return pt
}

// AddCallback adds a progress callback
func (pt *ProgressTracker) AddCallback(callback ProgressCallback) {
	pt.mu.Lock()
	defer pt.mu.Unlock()
	pt.callbacks = append(pt.callbacks, callback)
}

// UpdateProgress updates the progress
func (pt *ProgressTracker) UpdateProgress(completed, failed int, current string) {
	pt.mu.Lock()
	pt.progress.Completed = completed
	pt.progress.Failed = failed
	pt.progress.Current = current
	pt.mu.Unlock()
	
	// Send update to channel
	select {
	case pt.updateChan <- pt.getProgress():
	default:
		// Channel is full, skip this update
	}
}

// IncrementCompleted increments the completed count
func (pt *ProgressTracker) IncrementCompleted(current string) {
	pt.mu.Lock()
	pt.progress.Completed++
	pt.progress.Current = current
	pt.mu.Unlock()
	
	select {
	case pt.updateChan <- pt.getProgress():
	default:
	}
}

// IncrementFailed increments the failed count
func (pt *ProgressTracker) IncrementFailed(current string) {
	pt.mu.Lock()
	pt.progress.Failed++
	pt.progress.Current = current
	pt.mu.Unlock()
	
	select {
	case pt.updateChan <- pt.getProgress():
	default:
	}
}

// GetProgress returns a copy of the current progress
func (pt *ProgressTracker) GetProgress() *Progress {
	return pt.getProgress()
}

// getProgress returns a copy of the current progress (internal use)
func (pt *ProgressTracker) getProgress() *Progress {
	pt.mu.RLock()
	defer pt.mu.RUnlock()
	
	return &Progress{
		Total:     pt.progress.Total,
		Completed: pt.progress.Completed,
		Failed:    pt.progress.Failed,
		Current:   pt.progress.Current,
		StartTime: pt.progress.StartTime,
	}
}

// Stop stops the progress tracker
func (pt *ProgressTracker) Stop() {
	close(pt.done)
}

// updateLoop handles progress updates
func (pt *ProgressTracker) updateLoop() {
	for {
		select {
		case progress := <-pt.updateChan:
			pt.mu.RLock()
			callbacks := make([]ProgressCallback, len(pt.callbacks))
			copy(callbacks, pt.callbacks)
			pt.mu.RUnlock()
			
			// Call all callbacks
			for _, callback := range callbacks {
				callback(progress)
			}
		case <-pt.done:
			return
		}
	}
}

// GetPercentage returns the completion percentage
func (p *Progress) GetPercentage() float64 {
	if p.Total == 0 {
		return 0
	}
	return float64(p.Completed+p.Failed) / float64(p.Total) * 100
}

// GetElapsedTime returns the elapsed time
func (p *Progress) GetElapsedTime() time.Duration {
	return time.Since(p.StartTime)
}

// GetEstimatedTimeRemaining returns estimated time remaining
func (p *Progress) GetEstimatedTimeRemaining() time.Duration {
	completed := p.Completed + p.Failed
	if completed == 0 {
		return 0
	}
	
	elapsed := p.GetElapsedTime()
	avgTimePerItem := elapsed / time.Duration(completed)
	remaining := p.Total - completed
	
	return avgTimePerItem * time.Duration(remaining)
}

// GetRate returns the rate of completion (items per second)
func (p *Progress) GetRate() float64 {
	elapsed := p.GetElapsedTime()
	if elapsed.Seconds() == 0 {
		return 0
	}
	return float64(p.Completed+p.Failed) / elapsed.Seconds()
}

// String returns a string representation of the progress
func (p *Progress) String() string {
	percentage := p.GetPercentage()
	elapsed := p.GetElapsedTime()
	rate := p.GetRate()
	
	return fmt.Sprintf("[%d/%d] %.1f%% - %s - %.1f items/sec - %v elapsed",
		p.Completed+p.Failed, p.Total, percentage, p.Current, rate, elapsed.Round(time.Second))
}

// DefaultProgressCallback provides a default progress display
func DefaultProgressCallback(p *Progress) {
	fmt.Printf("\r%s", p.String())
}

// DetailedProgressCallback provides detailed progress information
func DetailedProgressCallback(p *Progress) {
	percentage := p.GetPercentage()
	elapsed := p.GetElapsedTime()
	remaining := p.GetEstimatedTimeRemaining()
	rate := p.GetRate()
	
	fmt.Printf("\r[%d/%d] %.1f%% | %s | %.1f/sec | %v elapsed | ~%v remaining",
		p.Completed+p.Failed, p.Total, percentage, p.Current, rate, 
		elapsed.Round(time.Second), remaining.Round(time.Second))
}

// SimpleProgressCallback provides simple progress information
func SimpleProgressCallback(p *Progress) {
	fmt.Printf("\r[%d/%d] %s", p.Completed+p.Failed, p.Total, p.Current)
}

// NewLineProgressCallback provides progress with newlines
func NewLineProgressCallback(p *Progress) {
	fmt.Printf("[%d/%d] %s\n", p.Completed+p.Failed, p.Total, p.Current)
}
