package nutrition

import (
	"github.com/habbazettt/nutrisnap-server/internal/models"
)

// CalculateNutriScore computes the Nutri-Score (A-E) and numeric score
// Based on approximate algorithm: https://en.wikipedia.org/wiki/Nutri-Score
func CalculateNutriScore(n *models.Nutrients) (string, int) {
	if n == nil {
		return "", 0
	}

	// 1. Calculate Negative Points (N)
	pointsN := 0
	pointsN += pointsEnergy(n.EnergyKcal)
	pointsN += pointsSugar(n.SugarG)
	pointsN += pointsSaturatedFat(n.SaturatedFatG)
	pointsN += pointsSodium(n.SodiumMg)

	// 2. Calculate Positive Points (P)
	pointsP := 0
	pointsP += pointsFiber(n.FiberG)
	pointsP += pointsProtein(n.ProteinG)
	// Fruit/Veg points usually require ingredients list analysis, assume 0 for OCR

	// 3. Final Score
	// Logic: Final = N - P
	// Exception: If N >= 11 and Fruit points < 5, Protein is not counted (unless Fruit points >= 5)
	// Simplified implementation:
	score := pointsN - pointsP

	// 4. Determine Grade
	grade := scoreToGrade(score)

	return grade, score
}

func scoreToGrade(score int) string {
	if score <= -1 {
		return "A"
	}
	if score <= 2 {
		return "B"
	}
	if score <= 10 {
		return "C"
	}
	if score <= 18 {
		return "D"
	}
	return "E"
}

// Points Calculation Helpers (Based on standard solid foods table)

func pointsEnergy(kcal *float64) int {
	if kcal == nil {
		return 0
	}
	kJ := *kcal * 4.184
	if kJ <= 335 {
		return 0
	}
	if kJ <= 670 {
		return 1
	}
	if kJ <= 1005 {
		return 2
	}
	if kJ <= 1340 {
		return 3
	}
	if kJ <= 1675 {
		return 4
	}
	if kJ <= 2010 {
		return 5
	}
	if kJ <= 2345 {
		return 6
	}
	if kJ <= 2680 {
		return 7
	}
	if kJ <= 3015 {
		return 8
	}
	if kJ <= 3350 {
		return 9
	}
	return 10
}

func pointsSugar(g *float64) int {
	if g == nil {
		return 0
	}
	v := *g
	if v <= 4.5 {
		return 0
	}
	if v <= 9 {
		return 1
	}
	if v <= 13.5 {
		return 2
	}
	if v <= 18 {
		return 3
	}
	if v <= 22.5 {
		return 4
	}
	if v <= 27 {
		return 5
	}
	if v <= 31 {
		return 6
	}
	if v <= 36 {
		return 7
	}
	if v <= 40 {
		return 8
	}
	if v <= 45 {
		return 9
	}
	return 10
}

func pointsSaturatedFat(g *float64) int {
	if g == nil {
		return 0
	}
	v := *g
	if v <= 1 {
		return 0
	}
	if v <= 2 {
		return 1
	}
	if v <= 3 {
		return 2
	}
	if v <= 4 {
		return 3
	}
	if v <= 5 {
		return 4
	}
	if v <= 6 {
		return 5
	}
	if v <= 7 {
		return 6
	}
	if v <= 8 {
		return 7
	}
	if v <= 9 {
		return 8
	}
	if v <= 10 {
		return 9
	}
	return 10
}

func pointsSodium(mg *float64) int {
	if mg == nil {
		return 0
	}
	v := *mg
	if v <= 90 {
		return 0
	}
	if v <= 180 {
		return 1
	}
	if v <= 270 {
		return 2
	}
	if v <= 360 {
		return 3
	}
	if v <= 450 {
		return 4
	}
	if v <= 540 {
		return 5
	}
	if v <= 630 {
		return 6
	}
	if v <= 720 {
		return 7
	}
	if v <= 810 {
		return 8
	}
	if v <= 900 {
		return 9
	}
	return 10
}

func pointsFiber(g *float64) int {
	if g == nil {
		return 0
	}
	v := *g
	if v <= 0.9 {
		return 0
	}
	if v <= 1.9 {
		return 1
	}
	if v <= 2.8 {
		return 2
	}
	if v <= 3.7 {
		return 3
	}
	if v <= 4.7 {
		return 4
	}
	return 5
}

func pointsProtein(g *float64) int {
	if g == nil {
		return 0
	}
	v := *g
	if v <= 1.6 {
		return 0
	}
	if v <= 3.2 {
		return 1
	}
	if v <= 4.8 {
		return 2
	}
	if v <= 6.4 {
		return 3
	}
	if v <= 8.0 {
		return 4
	}
	return 5
}
