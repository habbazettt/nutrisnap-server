package models

import (
	"encoding/json"
)

type ProductSource string

const (
	SourceOpenFoodFacts ProductSource = "openfoodfacts"
	SourceOCRScan       ProductSource = "ocr_scan"
	SourceManual        ProductSource = "manual"
)

type Nutrients struct {
	EnergyKcal    *float64 `json:"energy_kcal,omitempty"`
	ProteinG      *float64 `json:"protein_g,omitempty"`
	CarbohydrateG *float64 `json:"carbohydrate_g,omitempty"`
	SugarG        *float64 `json:"sugar_g,omitempty"`
	FatG          *float64 `json:"fat_g,omitempty"`
	SaturatedFatG *float64 `json:"saturated_fat_g,omitempty"`
	FiberG        *float64 `json:"fiber_g,omitempty"`
	SodiumMg      *float64 `json:"sodium_mg,omitempty"`
	SaltG         *float64 `json:"salt_g,omitempty"`
	CholesterolMg *float64 `json:"cholesterol_mg,omitempty"`
	TransFatG     *float64 `json:"trans_fat_g,omitempty"`
	VitaminAIU    *float64 `json:"vitamin_a_iu,omitempty"`
	VitaminCMg    *float64 `json:"vitamin_c_mg,omitempty"`
	CalciumMg     *float64 `json:"calcium_mg,omitempty"`
	IronMg        *float64 `json:"iron_mg,omitempty"`
	PotassiumMg   *float64 `json:"potassium_mg,omitempty"`
}

type Product struct {
	BaseWithoutSoftDelete
	Barcode              string        `gorm:"uniqueIndex;size:50" json:"barcode"`
	Name                 string        `gorm:"size:500;not null" json:"name"`
	Brand                *string       `gorm:"size:255" json:"brand,omitempty"`
	ImageURL             *string       `gorm:"size:1000" json:"image_url,omitempty"`
	Source               ProductSource `gorm:"type:varchar(50);default:manual" json:"source"`
	NutrientsJSON        JSON          `gorm:"type:jsonb" json:"nutrients"`
	ServingSize          *string       `gorm:"size:100" json:"serving_size,omitempty"`
	ServingNutrientsJSON JSON          `gorm:"type:jsonb" json:"serving_nutrients,omitempty"`
	NutriScore           *string       `gorm:"size:1" json:"nutri_score,omitempty"`
	NutriScoreValue      *int          `json:"nutri_score_value,omitempty"`
	HighlightsJSON       JSON          `gorm:"type:jsonb" json:"highlights,omitempty"`
	InsightsJSON         JSON          `gorm:"type:jsonb" json:"insights,omitempty"`

	// Relations
	Scans []Scan `gorm:"foreignKey:ProductID" json:"scans,omitempty"`
}

func (Product) TableName() string {
	return "products"
}

func (p *Product) GetNutrients() (*Nutrients, error) {
	if p.NutrientsJSON == nil {
		return nil, nil
	}

	var nutrients Nutrients
	if err := json.Unmarshal(p.NutrientsJSON, &nutrients); err != nil {
		return nil, err
	}
	return &nutrients, nil
}

func (p *Product) SetNutrients(nutrients *Nutrients) error {
	if nutrients == nil {
		p.NutrientsJSON = nil
		return nil
	}

	data, err := json.Marshal(nutrients)
	if err != nil {
		return err
	}
	p.NutrientsJSON = data
	return nil
}

type JSON []byte

func (j JSON) Value() (interface{}, error) {
	if len(j) == 0 {
		return nil, nil
	}
	return []byte(j), nil
}

func (j *JSON) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}

	switch v := value.(type) {
	case []byte:
		*j = append((*j)[0:0], v...)
	case string:
		*j = []byte(v)
	}
	return nil
}

func (j JSON) MarshalJSON() ([]byte, error) {
	if len(j) == 0 {
		return []byte("null"), nil
	}
	return j, nil
}

func (j *JSON) UnmarshalJSON(data []byte) error {
	if data == nil || string(data) == "null" {
		*j = nil
		return nil
	}
	*j = append((*j)[0:0], data...)
	return nil
}
