package jobs

import (
	"context"
	"log"
	"time"

	"github.com/habbazettt/nutrisnap-server/internal/repositories"
	"github.com/habbazettt/nutrisnap-server/pkg/storage"
)

// CleanupConfig holds cleanup job configuration
type CleanupConfig struct {
	// RetentionDays is the number of days to keep images
	RetentionDays int
	// Interval is how often to run the cleanup
	Interval time.Duration
	// BatchSize is the number of records to process per batch
	BatchSize int
}

// DefaultCleanupConfig returns default cleanup configuration
func DefaultCleanupConfig() CleanupConfig {
	return CleanupConfig{
		RetentionDays: 30,             // Keep images for 30 days
		Interval:      24 * time.Hour, // Run daily
		BatchSize:     100,            // Process 100 records per batch
	}
}

// CleanupJob handles periodic cleanup of old images
type CleanupJob struct {
	config        CleanupConfig
	scanRepo      repositories.ScanRepository
	storageClient *storage.Client
	stopChan      chan struct{}
	isRunning     bool
}

// NewCleanupJob creates a new cleanup job
func NewCleanupJob(config CleanupConfig, scanRepo repositories.ScanRepository, storageClient *storage.Client) *CleanupJob {
	return &CleanupJob{
		config:        config,
		scanRepo:      scanRepo,
		storageClient: storageClient,
		stopChan:      make(chan struct{}),
	}
}

// Start starts the cleanup job scheduler
func (j *CleanupJob) Start() {
	if j.isRunning {
		return
	}
	j.isRunning = true

	go func() {
		// Run immediately on start
		j.runCleanup()

		ticker := time.NewTicker(j.config.Interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				j.runCleanup()
			case <-j.stopChan:
				log.Println("Cleanup job stopped")
				return
			}
		}
	}()

	log.Printf("Cleanup job started (retention: %d days, interval: %s)", j.config.RetentionDays, j.config.Interval)
}

// Stop stops the cleanup job
func (j *CleanupJob) Stop() {
	if !j.isRunning {
		return
	}
	close(j.stopChan)
	j.isRunning = false
}

// runCleanup performs the actual cleanup
func (j *CleanupJob) runCleanup() {
	ctx := context.Background()
	cutoffDate := time.Now().AddDate(0, 0, -j.config.RetentionDays)

	log.Printf("Running cleanup for images older than %s", cutoffDate.Format("2006-01-02"))

	// Get old scans with stored images
	scans, err := j.scanRepo.FindOldScansWithImages(cutoffDate, j.config.BatchSize)
	if err != nil {
		log.Printf("Error finding old scans: %v", err)
		return
	}

	if len(scans) == 0 {
		log.Println("No old images to clean up")
		return
	}

	deletedCount := 0
	failedCount := 0

	for _, scan := range scans {
		if scan.ImageRef == nil || !scan.ImageStored {
			continue
		}

		// Delete image from MinIO
		if err := j.storageClient.Delete(ctx, *scan.ImageRef); err != nil {
			log.Printf("Failed to delete image %s: %v", *scan.ImageRef, err)
			failedCount++
			continue
		}

		// Update scan record - mark image as not stored
		scan.ImageStored = false
		scan.ImageRef = nil
		if err := j.scanRepo.Update(&scan); err != nil {
			log.Printf("Failed to update scan %s: %v", scan.ID, err)
			failedCount++
			continue
		}

		deletedCount++
	}

	log.Printf("Cleanup completed: %d deleted, %d failed", deletedCount, failedCount)
}

// RunNow runs cleanup immediately (for manual trigger)
func (j *CleanupJob) RunNow() {
	go j.runCleanup()
}
