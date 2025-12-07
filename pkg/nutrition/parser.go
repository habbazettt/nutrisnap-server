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
	text = cleanupOCRText(text)
	text = strings.ReplaceAll(text, "|", " ")  // Pipe often read as separator
	text = strings.ReplaceAll(text, "\n", " ") // Treat newlines as space for regex simplicity

	nutrients := &models.Nutrients{}
	var servingSize string

	// Extract Serving Size
	// Pattern: "Serving Size ... 30g" or "Takaran Saji ... 30 g" or "Takaran Saji 30g"
	// Regex: (serving size|takaran saji).*?(\d+\s*[gk]?[gm]l?)
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
			// .{0,60}? : match any characters (non-greedy, max 60 chars) to skip separators like " / Total Fat :"
			// (\d+[.,]?\d*) : capture group for number
			pattern := fmt.Sprintf(`%s.{0,60}?(\d+[.,]?\d*)`, regexp.QuoteMeta(key))
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

	// Energy (Kcal priority)
	if val := extract([]string{"energi total", "energy", "kalori", "calories", "kcal"}); val != nil {
		nutrients.EnergyKcal = val
	}

	// Protein
	if val := extract([]string{"protein", "proteine"}); val != nil {
		nutrients.ProteinG = val
	}

	// Fat (Total Fat)
	// Prioritize "Total Fat" over generic "Fat"
	if val := extract([]string{"lemak total", "total fat", "lemak", "fat", "lipides"}); val != nil {
		nutrients.FatG = val
	}

	// Saturated Fat
	if val := extract([]string{"lemak jenuh", "saturated fat", "sat fat"}); val != nil {
		nutrients.SaturatedFatG = val
	}

	// Carbohydrate
	if val := extract([]string{"karbohidrat", "total carb", "carbohydrate", "carb", "glucides"}); val != nil {
		nutrients.CarbohydrateG = val
	}

	// Sugar
	if val := extract([]string{"gula", "sugars", "sugar", "total sugar"}); val != nil {
		nutrients.SugarG = val
	}

	// Fiber (Serat)
	if val := extract([]string{"serat pangan", "serat", "fiber", "dietary fiber"}); val != nil {
		nutrients.FiberG = val
	}

	// Sodium (Natrium)
	if val := extract([]string{"natrium", "sodium"}); val != nil {
		nutrients.SodiumMg = val
	}

	return nutrients, servingSize
}

// cleanupOCRText fixes common OCR typos and noise
func cleanupOCRText(text string) string {
	text = strings.ToLower(text)

	// Dictionary-based fix
	replacements := map[string]string{
		"lomak":       "lemak",
		"lornak":      "lemak",
		"garm":        "garam",
		"kabohidar":   "karbohidrat",
		"karbohidart": "karbohidrat",

		// Noise fixes (specific to observed output)
		"0'3": "0.3",
		"0'5": "0.5",
		"0'0": "0.0",
	}

	for old, new := range replacements {
		text = strings.ReplaceAll(text, old, new)
	}

	// Regex-based fixes

	// Fix quote as decimal separator: 0'5 -> 0.5
	// Use regex carefully inside number
	reQuoteDecimal := regexp.MustCompile(`(\d)'(\d)`)
	text = reQuoteDecimal.ReplaceAllString(text, "$1.$2")

	return text
}
