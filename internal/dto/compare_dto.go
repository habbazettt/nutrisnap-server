package dto

// CompareRequest represents a product comparison request
// @Description Product comparison request
type CompareRequest struct {
	ProductA string `json:"product_a" validate:"required" example:"8992761136000"`
	ProductB string `json:"product_b" validate:"required" example:"8992388163138"`
}

// NutrientComparison represents a single nutrient comparison
// @Description Single nutrient comparison result
type NutrientComparison struct {
	Name       string   `json:"name" example:"Sugar"`
	ValueA     *float64 `json:"value_a,omitempty" example:"10.5"`
	ValueB     *float64 `json:"value_b,omitempty" example:"25.0"`
	Unit       string   `json:"unit" example:"g"`
	Difference *float64 `json:"difference,omitempty" example:"-14.5"`
	Winner     string   `json:"winner" example:"a"` // "a", "b", or "tie"
	Note       string   `json:"note,omitempty" example:"Product A has 58% less sugar"`
}

// ProductSummary represents a product summary for comparison
// @Description Product summary for comparison
type ProductSummary struct {
	ID         string  `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name       string  `json:"name" example:"Teh Botol Sosro"`
	Barcode    string  `json:"barcode,omitempty" example:"8992761136000"`
	NutriScore *string `json:"nutri_score,omitempty" example:"B"`
	ImageURL   *string `json:"image_url,omitempty"`
}

// CompareResponse represents the comparison result
// @Description Product comparison result
type CompareResponse struct {
	ProductA    ProductSummary       `json:"product_a"`
	ProductB    ProductSummary       `json:"product_b"`
	Comparisons []NutrientComparison `json:"comparisons"`
	Winner      string               `json:"winner" example:"a"` // "a", "b", or "tie"
	Verdict     string               `json:"verdict" example:"Product A lebih sehat karena memiliki NutriScore lebih baik (B vs D)"`
}
