package services

import (
	"context"
	"errors"

	"github.com/habbazettt/nutrisnap-server/internal/models"
	"github.com/habbazettt/nutrisnap-server/internal/repositories"
	"github.com/habbazettt/nutrisnap-server/pkg/openfoodfacts"
)

type ProductService interface {
	GetProductByBarcode(ctx context.Context, barcode string) (*models.Product, error)
}

type productService struct {
	productRepo repositories.ProductRepository
	offClient   *openfoodfacts.Client
}

func NewProductService(productRepo repositories.ProductRepository, offClient *openfoodfacts.Client) ProductService {
	return &productService{
		productRepo: productRepo,
		offClient:   offClient,
	}
}

func (s *productService) GetProductByBarcode(ctx context.Context, barcode string) (*models.Product, error) {
	product, err := s.productRepo.FindByBarcode(barcode)
	if err == nil {
		return product, nil
	}

	if !errors.Is(err, repositories.ErrProductNotFound) {
		// Real DB error
		return nil, err
	}

	// 2. Fetch from OpenFoodFacts
	offProduct, err := s.offClient.GetProduct(barcode)
	if err != nil {
		return nil, err // External API error
	}

	if offProduct == nil {
		// Not found in OFF
		return nil, repositories.ErrProductNotFound
	}

	// 3. Save to local DB (Cache)
	if err := s.productRepo.Create(offProduct); err != nil {
		// Log error but return product anyway
		// logger.Error("Failed to cache OFF product", "err", err)
	}

	return offProduct, nil
}
