package nutrition

import (
	"github.com/habbazettt/nutrisnap-server/internal/models"
)

// Analyze generates highlights and insights based on nutrient values
func Analyze(n *models.Nutrients) ([]models.NutrientHighlight, []models.Insight) {
	if n == nil {
		return nil, nil
	}

	highlights := make([]models.NutrientHighlight, 0)
	insights := make([]models.Insight, 0)

	// --- 1. Sugar Analysis ---
	if n.SugarG != nil {
		val := *n.SugarG
		if val > 22.5 {
			highlights = append(highlights, models.NutrientHighlight{
				Nutrient: "Sugar",
				Level:    "high",
				Value:    val,
				Unit:     "g",
				Message:  "High Sugar",
			})
			insights = append(insights, models.Insight{
				Type:     "health",
				Title:    "Limit Intake",
				Message:  "Content contains high level of sugar.",
				Severity: "warning",
			})
		} else if val < 5 {
			highlights = append(highlights, models.NutrientHighlight{
				Nutrient: "Sugar",
				Level:    "low",
				Value:    val,
				Unit:     "g",
				Message:  "Low Sugar",
			})
		} else {
			highlights = append(highlights, models.NutrientHighlight{
				Nutrient: "Sugar",
				Level:    "medium",
				Value:    val,
				Unit:     "g",
				Message:  "Moderate Sugar",
			})
		}
	}

	// --- 2. Fat Analysis ---
	if n.FatG != nil {
		val := *n.FatG
		if val > 17.5 {
			highlights = append(highlights, models.NutrientHighlight{
				Nutrient: "Fat",
				Level:    "high",
				Value:    val,
				Unit:     "g",
				Message:  "High Fat",
			})
		} else if val < 3 {
			highlights = append(highlights, models.NutrientHighlight{
				Nutrient: "Fat",
				Level:    "low",
				Value:    val,
				Unit:     "g",
				Message:  "Low Fat",
			})
		}
	}

	// --- 3. Saturated Fat Analysis ---
	if n.SaturatedFatG != nil {
		val := *n.SaturatedFatG
		if val > 5 {
			highlights = append(highlights, models.NutrientHighlight{
				Nutrient: "Saturated Fat",
				Level:    "high",
				Value:    val,
				Unit:     "g",
				Message:  "High Saturated Fat",
			})
			insights = append(insights, models.Insight{
				Type:     "health",
				Title:    "Warning",
				Message:  "High in saturated fats/trans fats.",
				Severity: "warning",
			})
		}
	}

	// --- 4. Protein Analysis ---
	if n.ProteinG != nil {
		val := *n.ProteinG
		if val > 10 { // > 10g per 100g is generally good source (solid food)
			highlights = append(highlights, models.NutrientHighlight{
				Nutrient: "Protein",
				Level:    "high",
				Value:    val,
				Unit:     "g",
				Message:  "High Protein",
			})
		}
	}

	// --- 5. Sodium/Salt Analysis ---
	// Approx: Sodium > 600mg (1.5g salt) is high
	if n.SodiumMg != nil {
		val := *n.SodiumMg
		if val > 600 {
			highlights = append(highlights, models.NutrientHighlight{
				Nutrient: "Sodium",
				Level:    "high",
				Value:    val,
				Unit:     "mg",
				Message:  "High Sodium",
			})
		}
	}

	return highlights, insights
}
