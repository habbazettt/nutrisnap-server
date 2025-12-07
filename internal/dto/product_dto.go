package dto

import "github.com/habbazettt/nutrisnap-server/internal/models"

type ProductResponse struct {
	ID              string            `json:"id"`
	Barcode         string            `json:"barcode"`
	Name            string            `json:"name"`
	Brand           *string           `json:"brand,omitempty"`
	ImageURL        *string           `json:"image_url,omitempty"`
	Source          string            `json:"source"`
	Nutrients       *models.Nutrients `json:"nutrients,omitempty"`
	ServingSize     *string           `json:"serving_size,omitempty"`
	NutriScore      *string           `json:"nutri_score,omitempty"`
	NutriScoreValue *int              `json:"nutri_score_value,omitempty"`
}

func ToProductResponse(p *models.Product) ProductResponse {
	nutrients, _ := p.GetNutrients()

	return ProductResponse{
		ID:              p.ID.String(),
		Barcode:         p.Barcode,
		Name:            p.Name,
		Brand:           p.Brand,
		ImageURL:        p.ImageURL,
		Source:          string(p.Source),
		Nutrients:       nutrients,
		ServingSize:     p.ServingSize,
		NutriScore:      p.NutriScore,
		NutriScoreValue: p.NutriScoreValue,
	}
}
