package services

import (
	"context"
	"fmt"
	"math"

	"github.com/habbazettt/nutrisnap-server/internal/dto"
	"github.com/habbazettt/nutrisnap-server/internal/models"
	"github.com/habbazettt/nutrisnap-server/internal/repositories"
)

type CompareService interface {
	CompareProducts(ctx context.Context, productAID, productBID string) (*dto.CompareResponse, error)
}

type compareService struct {
	productRepo repositories.ProductRepository
	scanRepo    repositories.ScanRepository
}

func NewCompareService(productRepo repositories.ProductRepository, scanRepo repositories.ScanRepository) CompareService {
	return &compareService{
		productRepo: productRepo,
		scanRepo:    scanRepo,
	}
}

func (s *compareService) CompareProducts(ctx context.Context, productAID, productBID string) (*dto.CompareResponse, error) {
	// Try to find products by barcode first, then by scan_id
	productA, err := s.findProduct(productAID)
	if err != nil {
		return nil, fmt.Errorf("product A not found: %w", err)
	}

	productB, err := s.findProduct(productBID)
	if err != nil {
		return nil, fmt.Errorf("product B not found: %w", err)
	}

	// Get nutrients
	nutrientsA, _ := productA.GetNutrients()
	nutrientsB, _ := productB.GetNutrients()

	// Build comparison
	comparisons := s.compareNutrients(nutrientsA, nutrientsB)

	// Determine winner and verdict
	winner, verdict := s.generateVerdict(productA, productB, comparisons)

	return &dto.CompareResponse{
		ProductA: dto.ProductSummary{
			ID:         productA.ID.String(),
			Name:       productA.Name,
			Barcode:    productA.Barcode,
			NutriScore: productA.NutriScore,
			ImageURL:   productA.ImageURL,
		},
		ProductB: dto.ProductSummary{
			ID:         productB.ID.String(),
			Name:       productB.Name,
			Barcode:    productB.Barcode,
			NutriScore: productB.NutriScore,
			ImageURL:   productB.ImageURL,
		},
		Comparisons: comparisons,
		Winner:      winner,
		Verdict:     verdict,
	}, nil
}

func (s *compareService) findProduct(identifier string) (*models.Product, error) {
	// Try barcode first
	product, err := s.productRepo.FindByBarcode(identifier)
	if err == nil && product != nil {
		return product, nil
	}

	// Try scan_id
	scan, err := s.scanRepo.FindByID(identifier)
	if err == nil && scan != nil && scan.Product != nil {
		return scan.Product, nil
	}

	return nil, fmt.Errorf("not found")
}

func (s *compareService) compareNutrients(a, b *models.Nutrients) []dto.NutrientComparison {
	comparisons := []dto.NutrientComparison{}

	// Define nutrients to compare (lower is better for most)
	type nutrientDef struct {
		name        string
		unit        string
		valueA      *float64
		valueB      *float64
		lowerBetter bool
	}

	if a == nil {
		a = &models.Nutrients{}
	}
	if b == nil {
		b = &models.Nutrients{}
	}

	nutrients := []nutrientDef{
		{"Calories", "kcal", a.EnergyKcal, b.EnergyKcal, true},
		{"Protein", "g", a.ProteinG, b.ProteinG, false}, // higher is better
		{"Fat", "g", a.FatG, b.FatG, true},
		{"Saturated Fat", "g", a.SaturatedFatG, b.SaturatedFatG, true},
		{"Carbohydrate", "g", a.CarbohydrateG, b.CarbohydrateG, true},
		{"Sugar", "g", a.SugarG, b.SugarG, true},
		{"Fiber", "g", a.FiberG, b.FiberG, false}, // higher is better
		{"Sodium", "mg", a.SodiumMg, b.SodiumMg, true},
	}

	for _, n := range nutrients {
		comp := dto.NutrientComparison{
			Name:   n.name,
			ValueA: n.valueA,
			ValueB: n.valueB,
			Unit:   n.unit,
		}

		// Calculate difference and winner
		if n.valueA != nil && n.valueB != nil {
			diff := *n.valueA - *n.valueB
			comp.Difference = &diff

			if n.lowerBetter {
				if *n.valueA < *n.valueB {
					comp.Winner = "a"
					pct := (*n.valueB - *n.valueA) / *n.valueB * 100
					comp.Note = fmt.Sprintf("%.0f%% less %s", math.Abs(pct), n.name)
				} else if *n.valueB < *n.valueA {
					comp.Winner = "b"
					pct := (*n.valueA - *n.valueB) / *n.valueA * 100
					comp.Note = fmt.Sprintf("%.0f%% less %s", math.Abs(pct), n.name)
				} else {
					comp.Winner = "tie"
				}
			} else {
				// Higher is better (protein, fiber)
				if *n.valueA > *n.valueB {
					comp.Winner = "a"
					pct := (*n.valueA - *n.valueB) / *n.valueA * 100
					comp.Note = fmt.Sprintf("%.0f%% more %s", math.Abs(pct), n.name)
				} else if *n.valueB > *n.valueA {
					comp.Winner = "b"
					pct := (*n.valueB - *n.valueA) / *n.valueB * 100
					comp.Note = fmt.Sprintf("%.0f%% more %s", math.Abs(pct), n.name)
				} else {
					comp.Winner = "tie"
				}
			}
		} else {
			comp.Winner = "unknown"
		}

		comparisons = append(comparisons, comp)
	}

	return comparisons
}

func (s *compareService) generateVerdict(a, b *models.Product, comparisons []dto.NutrientComparison) (string, string) {
	// Primary: Compare NutriScore
	if a.NutriScore != nil && b.NutriScore != nil {
		scoreA := *a.NutriScore
		scoreB := *b.NutriScore

		if scoreA < scoreB { // A is better (A < B < C < D < E)
			return "a", fmt.Sprintf("%s lebih sehat dengan NutriScore %s vs %s", a.Name, scoreA, scoreB)
		} else if scoreB < scoreA {
			return "b", fmt.Sprintf("%s lebih sehat dengan NutriScore %s vs %s", b.Name, scoreB, scoreA)
		}
	}

	// Secondary: Count wins
	winsA, winsB := 0, 0
	for _, c := range comparisons {
		if c.Winner == "a" {
			winsA++
		} else if c.Winner == "b" {
			winsB++
		}
	}

	if winsA > winsB {
		return "a", fmt.Sprintf("%s unggul di %d dari %d kategori nutrisi", a.Name, winsA, len(comparisons))
	} else if winsB > winsA {
		return "b", fmt.Sprintf("%s unggul di %d dari %d kategori nutrisi", b.Name, winsB, len(comparisons))
	}

	return "tie", "Kedua produk memiliki profil nutrisi yang serupa"
}
