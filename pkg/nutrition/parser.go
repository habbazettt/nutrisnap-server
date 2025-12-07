package nutrition

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/habbazettt/nutrisnap-server/internal/models"
)

// ParseFromText extracts nutrient information and serving size from OCR text
func ParseFromText(text string) (*models.Nutrients, string) {
	// 1. Text Normalization
	text = strings.ToLower(text)
	text = strings.ReplaceAll(text, "|", " ")  // Pipe often read as separator
	text = strings.ReplaceAll(text, "\n", " ") // Treat newlines as space for regex simplicity

	nutrients := &models.Nutrients{}
	var servingSize string

	// Extract Serving Size
	// Pattern: "Serving Size ... 30g" or "Takaran Saji ... 30 g" or "Takaran Saji 30g"
	// Regex: (serving size|takaran saji).*?(\d+\s*[gk]?[gm]l?)
	// Capture group 2 is the value+unit
	reServing := regexp.MustCompile(`(serving size|takaran saji).*?(\d+\s*[gk]?[gm]l?)`)
	matchesServing := reServing.FindStringSubmatch(text)
	if len(matchesServing) > 2 {
		raw := matchesServing[2]
		// Clean up space: "30 g" -> "30g"
		servingSize = strings.ReplaceAll(raw, " ", "")
	}

	// Helper to extract value
	extract := func(keys []string) *float64 {
		for _, key := range keys {
			// Regex explanation:
			// %s : key word
			// \s* : optional whitespace
			// [:=-]? : optional separator (: or = or -)
			// \s* : optional whitespace
			// (\d+[.,]?\d*) : capture group for number
			pattern := fmt.Sprintf(`%s\s*[:=-]?\s*(\d+[.,]?\d*)`, regexp.QuoteMeta(key))
			re := regexp.MustCompile(pattern)

			matches := re.FindStringSubmatch(text)
			if len(matches) > 1 {
				numStr := matches[1]
				// Normalize comma to dot for parsing
				numStr = strings.Replace(numStr, ",", ".", 1)

				val, err := strconv.ParseFloat(numStr, 64)
				if err == nil {
					return &val
				}
			}
		}
		return nil
	}

	// 1. Energy (Kcal priority)
	if val := extract([]string{"energi total", "energy", "kalori", "calories", "kcal"}); val != nil {
		nutrients.EnergyKcal = val
	}

	// 2. Protein
	if val := extract([]string{"protein", "proteine"}); val != nil {
		nutrients.ProteinG = val
	}

	// 3. Fat (Total Fat)
	// Prioritize "Total Fat" over generic "Fat"
	if val := extract([]string{"lemak total", "total fat", "lemak", "fat", "lipides"}); val != nil {
		nutrients.FatG = val
	}

	// 4. Saturated Fat
	if val := extract([]string{"lemak jenuh", "saturated fat", "sat fat"}); val != nil {
		nutrients.SaturatedFatG = val
	}

	// 5. Carbohydrate
	if val := extract([]string{"karbohidrat", "total carb", "carbohydrate", "carb", "glucides"}); val != nil {
		nutrients.CarbohydrateG = val
	}

	// 6. Sugar
	if val := extract([]string{"gula", "sugars", "sugar", "total sugar"}); val != nil {
		nutrients.SugarG = val
	}

	// 7. Fiber (Serat)
	if val := extract([]string{"serat pangan", "serat", "fiber", "dietary fiber"}); val != nil {
		nutrients.FiberG = val
	}

	// 8. Sodium (Natrium)
	if val := extract([]string{"natrium", "sodium"}); val != nil {
		nutrients.SodiumMg = val
	}

	return nutrients, servingSize
}
