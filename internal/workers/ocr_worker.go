package workers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/habbazettt/nutrisnap-server/internal/models"
	"github.com/habbazettt/nutrisnap-server/internal/repositories"
	"github.com/habbazettt/nutrisnap-server/internal/services"
	"github.com/habbazettt/nutrisnap-server/pkg/nutrition"
)

type OCRWorker struct {
	scanRepo    repositories.ScanRepository
	productRepo repositories.ProductRepository
	ocrService  services.OCRService
	scanQueue   chan string // Channel receiving ScanIDs
	quit        chan bool
}

func NewOCRWorker(scanRepo repositories.ScanRepository, productRepo repositories.ProductRepository, ocrService services.OCRService, bufferSize int) *OCRWorker {
	return &OCRWorker{
		scanRepo:    scanRepo,
		productRepo: productRepo,
		ocrService:  ocrService,
		scanQueue:   make(chan string, bufferSize),
		quit:        make(chan bool),
	}
}

// Start launches 'workers' number of goroutines
func (w *OCRWorker) Start(workers int) {
	for i := 0; i < workers; i++ {
		go w.run(i)
	}
	log.Printf("OCR Worker: Started %d worker(s)", workers)
}

// Stop signals all workers to stop
func (w *OCRWorker) Stop() {
	go func() {
		w.quit <- true
	}()
}

// EnqueueScan adds a scan ID to the processing queue
func (w *OCRWorker) EnqueueScan(scanID string) {
	select {
	case w.scanQueue <- scanID:
		// Successfully queued
	default:
		log.Printf("OCR Worker: Queue full, dropping scan %s", scanID)
	}
}

func (w *OCRWorker) run(id int) {
	for {
		select {
		case scanID := <-w.scanQueue:
			// Process scan
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
			if err := w.processScan(ctx, scanID); err != nil {
				log.Printf("OCR Worker [%d]: Failed to process scan %s: %v", id, scanID, err)
			} else {
				log.Printf("OCR Worker [%d]: Successfully processed scan %s", id, scanID)
			}
			cancel()
		case <-w.quit:
			log.Printf("OCR Worker [%d]: Stopping", id)
			return
		}
	}
}

func (w *OCRWorker) processScan(ctx context.Context, scanID string) error {
	// 1. Get Scan
	scan, err := w.scanRepo.FindByID(scanID)
	if err != nil {
		return fmt.Errorf("scan not found: %w", err)
	}

	if scan.ImageRef == nil || !scan.ImageStored {
		return fmt.Errorf("no image to process")
	}

	// Update status to processing
	scan.Status = "processing"
	w.scanRepo.Update(scan)

	// 2. Run OCR
	nutrients, servingSize, err := w.ocrService.ProcessImageFromStorage(ctx, *scan.ImageRef)
	if err != nil {
		scan.Status = "failed"
		// Append error?
		w.scanRepo.Update(scan)
		return err
	}

	// 3. Process Nutrition Data (Analysis, Score, etc.)
	// Calculate NutriScore
	grade, score := nutrition.CalculateNutriScore(nutrients)

	// Analyze Highlights & Insights
	highlights, insights := nutrition.Analyze(nutrients)

	// Marshal JSONs
	nutrientsJSON, _ := json.Marshal(nutrients)
	highlightsJSON, _ := json.Marshal(highlights)
	insightsJSON, _ := json.Marshal(insights)

	ocrBarcode := fmt.Sprintf("ocr-%s", scanID) // Unique pseudo-barcode

	var servingSizePtr *string
	if servingSize != "" {
		servingSizePtr = &servingSize
	}

	product := &models.Product{
		Barcode:         ocrBarcode,
		Name:            "Scanned Product " + time.Now().Format("02-Jan 15:04"),
		Source:          models.SourceOCRScan,
		NutrientsJSON:   nutrientsJSON,
		ServingSize:     servingSizePtr,
		NutriScore:      &grade,
		NutriScoreValue: &score,
		HighlightsJSON:  highlightsJSON,
		InsightsJSON:    insightsJSON,
	}

	if err := w.productRepo.Create(product); err != nil {
		// Possibly duplicate if re-scanning?
		scan.Status = "failed"
		w.scanRepo.Update(scan)
		return fmt.Errorf("failed to create ocr product: %w", err)
	}

	// 4. Link Product to Scan and Complete
	scan.ProductID = &product.ID

	// Update redundant Scan fields (optional but good for consistency if queries use Scan table)
	scan.NutriScore = &grade
	scan.NutriScoreValue = &score
	scan.HighlightsJSON = highlightsJSON
	scan.InsightsJSON = insightsJSON

	scan.Status = models.ScanStatusCompleted

	if err := w.scanRepo.Update(scan); err != nil {
		return fmt.Errorf("failed to update scan status: %w", err)
	}

	return nil
}
