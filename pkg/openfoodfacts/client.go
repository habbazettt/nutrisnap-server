package openfoodfacts

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/habbazettt/nutrisnap-server/internal/models"
)

const (
	BaseURL = "https://world.openfoodfacts.org/api/v0/product"
)

type Client struct {
	httpClient *http.Client
	userAgent  string
}

func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		userAgent: "NutriSnap - Android - Version 1.0 - www.nutrisnap.app",
	}
}

// GetProduct fetches product details by barcode
func (c *Client) GetProduct(barcode string) (*models.Product, error) {
	// Create request
	url := fmt.Sprintf("%s/%s.json", BaseURL, barcode)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", c.userAgent)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		fmt.Printf("OFF Client: Network error fetching %s: %v\n", url, err)
		return nil, fmt.Errorf("failed to fetch product from OFF: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		fmt.Printf("OFF Client: Product not found (404) for %s\n", barcode)
		return nil, nil // Product not found
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("OFF Client: Unexpected status code %d for %s\n", resp.StatusCode, barcode)
		return nil, fmt.Errorf("OFF API returned status: %d", resp.StatusCode)
	}

	var result ProductResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		fmt.Printf("OFF Client: Decode error for %s: %v\n", barcode, err)
		return nil, fmt.Errorf("failed to decode OFF response: %w", err)
	}

	// Debug result status
	// fmt.Printf("OFF Client: Result status: %v (Type: %T)\n", result.Status, result.Status)

	// Check status
	found := false
	switch v := result.Status.(type) {
	case float64:
		found = v == 1
	case int:
		found = v == 1
	case string:
		found = v == "1"
	}

	if !found || result.Product == nil {
		fmt.Printf("OFF Client: Product logical not found (status not 1 or nil product) for %s. Status: %v\n", barcode, result.Status)
		return nil, nil // Not found
	}

	return c.mapToModel(barcode, result.Product), nil
}

func (c *Client) mapToModel(barcode string, offProduct *Product) *models.Product {
	nutrients := &models.Nutrients{
		EnergyKcal:    toFloat(offProduct.Nutriments.EnergyKcal),
		ProteinG:      toFloat(offProduct.Nutriments.Proteins),
		CarbohydrateG: toFloat(offProduct.Nutriments.Carbohydrates),
		SugarG:        toFloat(offProduct.Nutriments.Sugars),
		FatG:          toFloat(offProduct.Nutriments.Fat),
		SaturatedFatG: toFloat(offProduct.Nutriments.SaturatedFat),
		FiberG:        toFloat(offProduct.Nutriments.Fiber),
		SodiumMg:      toMg(offProduct.Nutriments.Sodium), // OFF usually returns grams for sodium too
		SaltG:         toFloat(offProduct.Nutriments.Salt),
		CholesterolMg: toMg(offProduct.Nutriments.Cholesterol),
		VitaminAIU:    toFloat(offProduct.Nutriments.VitaminA), // Need to check unit usually
		VitaminCMg:    toMg(offProduct.Nutriments.VitaminC),
		CalciumMg:     toMg(offProduct.Nutriments.Calcium),
		IronMg:        toMg(offProduct.Nutriments.Iron),
		PotassiumMg:   toMg(offProduct.Nutriments.Potassium),
	}

	// Basic validation: if no energy, assume incomplete data
	// But we store what we get

	nutrientsJSON, _ := json.Marshal(nutrients)

	// Determine NutriScore value (A=1, E=5 or score calculation?)
	// OFF gives "a", "b", "c", "d", "e"
	var nutriScoreVal int
	switch offProduct.NutriscoreGrade {
	case "a":
		nutriScoreVal = 1
	case "b":
		nutriScoreVal = 2
	case "c":
		nutriScoreVal = 3
	case "d":
		nutriScoreVal = 4
	case "e":
		nutriScoreVal = 5
	}

	return &models.Product{
		Barcode:         barcode,
		Name:            offProduct.ProductName,
		Brand:           &offProduct.Brands,
		ImageURL:        &offProduct.ImageURL,
		Source:          models.SourceOpenFoodFacts,
		NutrientsJSON:   nutrientsJSON,
		ServingSize:     &offProduct.ServingSize,
		NutriScore:      &offProduct.NutriscoreGrade,
		NutriScoreValue: &nutriScoreVal,
	}
}

// Helper to convert interface{} to *float64
func toFloat(v interface{}) *float64 {
	if v == nil {
		return nil
	}
	switch val := v.(type) {
	case float64:
		return &val
	case int:
		f := float64(val)
		return &f
	case string:
		if f, err := strconv.ParseFloat(val, 64); err == nil {
			return &f
		}
	}
	return nil
}

// Convert grams (default in OFF) to mg (expected in our model for some fields)
// Note: Only for fields where our model expects mg but OFF gives g
// Need to be careful. OFF fields like sodium_100g is in GRAMS.
// Model SodiumMg is in MILLIGRAMS.
func toMg(v interface{}) *float64 {
	val := toFloat(v)
	if val == nil {
		return nil
	}
	mg := *val * 1000
	return &mg
}
