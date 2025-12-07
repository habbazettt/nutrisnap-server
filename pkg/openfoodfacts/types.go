package openfoodfacts

// ProductResponse represents the API response from OpenFoodFacts
type ProductResponse struct {
	Code          string      `json:"code"`
	Status        interface{} `json:"status"` // Can be int or string
	StatusVerbose string      `json:"status_verbose"`
	Product       *Product    `json:"product"`
}

// Product represents the product details from OFF
type Product struct {
	ProductName     string     `json:"product_name"`
	Brands          string     `json:"brands"`
	ImageURL        string     `json:"image_url"`
	Nutriments      Nutriments `json:"nutriments"`
	NutriscoreGrade string     `json:"nutriscore_grade"`
	ServingSize     string     `json:"serving_size"`
}

// Nutriments represents nutrition facts from OFF (all _100g)
type Nutriments struct {
	EnergyKcal    interface{} `json:"energy-kcal_100g"` // Can be string or float
	Proteins      interface{} `json:"proteins_100g"`
	Carbohydrates interface{} `json:"carbohydrates_100g"`
	Sugars        interface{} `json:"sugars_100g"`
	Fat           interface{} `json:"fat_100g"`
	SaturatedFat  interface{} `json:"saturated-fat_100g"`
	Fiber         interface{} `json:"fiber_100g"`
	Sodium        interface{} `json:"sodium_100g"`
	Salt          interface{} `json:"salt_100g"`
	Cholesterol   interface{} `json:"cholesterol_100g"` // Often missing
	VitaminA      interface{} `json:"vitamin-a_100g"`
	VitaminC      interface{} `json:"vitamin-c_100g"`
	Calcium       interface{} `json:"calcium_100g"`
	Iron          interface{} `json:"iron_100g"`
	Potassium     interface{} `json:"potassium_100g"`
}
