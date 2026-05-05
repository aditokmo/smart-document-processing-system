package repository

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"backend/internal/domain/model"
	"backend/pkg/utils"
)

type ExtractionServiceImpl struct{}

func NewExtractionService() *ExtractionServiceImpl {
	return &ExtractionServiceImpl{}
}

func (e *ExtractionServiceImpl) Extract(content string, docType model.DocumentType) (*model.Document, error) {
	doc := &model.Document{
		ID:   utils.GenerateID(),
		Type: docType,
	}

	lines := strings.Split(content, "\n")

	log.Printf("Extracting data from document content:\n%+v", doc)

	doc.SupplierName = e.extractSupplierName(lines)
	doc.DocumentNumber = e.extractDocumentNumber(lines)
	doc.IssueDate = e.extractDate(lines, "issue")
	doc.DueDate = e.extractDate(lines, "due")
	doc.Currency = e.extractCurrency(lines)
	doc.LineItems = e.extractLineItems(lines)
	doc.Subtotal = e.extractAmount(lines, "subtotal")
	doc.Tax = e.extractAmount(lines, "tax")
	doc.Total = e.extractAmount(lines, "total")

	if doc.Subtotal == 0 && doc.Total > 0 {
		doc.Subtotal = doc.Total
		doc.Total = doc.Subtotal + doc.Tax
	} else if doc.Subtotal > 0 && doc.Total == 0 {
		doc.Total = doc.Subtotal + doc.Tax
	}

	if len(doc.LineItems) == 0 && doc.Total > 0 {
		doc.LineItems = []model.LineItem{
			{
				Description: "",
				Quantity:    1,
				Price:       doc.Subtotal,
				Total:       doc.Subtotal,
			},
		}
	}

	return doc, nil
}

func (e *ExtractionServiceImpl) extractSupplierName(lines []string) string {
	for _, line := range lines {
		lower := strings.ToLower(line)
		if strings.Contains(lower, "supplier") || strings.Contains(lower, "vendor") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				return strings.TrimSpace(parts[1])
			}
		}
	}
	return ""
}

func (e *ExtractionServiceImpl) extractDocumentNumber(lines []string) string {
	re := regexp.MustCompile(`(?i)(?:invoice|order|number|num|no|po)[:\s#]*([A-Z0-9\-]+)`)
	for _, line := range lines {
		matches := re.FindStringSubmatch(line)
		if len(matches) > 1 {
			return strings.ToUpper(strings.TrimSpace(matches[1]))
		}
	}
	return ""
}

func (e *ExtractionServiceImpl) extractDate(lines []string, dateType string) *time.Time {
	var pattern string
	if dateType == "issue" {
		pattern = `(?i)(?:issue\s*date|date)[:\s]*(.+)`
	} else {
		pattern = `(?i)` + dateType + `\s*date[:\s]*(.+)`
	}

	re := regexp.MustCompile(pattern)
	datePatterns := []string{`\d{4}-\d{2}-\d{2}`, `\d{1,2}[/\-]\d{1,2}[/\-]\d{4}`, `\d{1,2}\s+(?:Jan|Feb|Mar|Apr|May|Jun|Jul|Aug|Sep|Oct|Nov|Dec)[a-z]*\s+\d{4}`}
	formats := []string{"2006-01-02", "02/01/2006", "02-01-2006", "01/02/2006", "2006/01/02"}

	for i, line := range lines {
		line = strings.TrimSpace(line)

		matches := re.FindStringSubmatch(line)
		if len(matches) > 1 {
			for _, datePattern := range datePatterns {
				dateRe := regexp.MustCompile(datePattern)
				dateStr := dateRe.FindString(matches[1])
				if dateStr != "" {
					for _, format := range formats {
						if t, err := time.Parse(format, dateStr); err == nil {
							return &t
						}
					}
				}
			}
		}

		labelRe := regexp.MustCompile(`(?i)` + dateType + `\s*date`)
		if labelRe.MatchString(line) {
			for j := i + 1; j < len(lines) && j < i+3; j++ {
				nextLine := strings.TrimSpace(lines[j])
				if nextLine == "" {
					continue
				}

				for _, datePattern := range datePatterns {
					dateRe := regexp.MustCompile(datePattern)
					dateStr := dateRe.FindString(nextLine)
					if dateStr != "" {
						for _, format := range formats {
							if t, err := time.Parse(format, dateStr); err == nil {
								return &t
							}
						}
					}
				}
			}
		}
	}
	return nil
}

func (e *ExtractionServiceImpl) extractLineItems(lines []string) []model.LineItem {
	// Skip header lines and section markers
	skip := regexp.MustCompile(`(?i)^(description|subtotal|tax|total|qty|unit\s*price|quantity|price|invoice|number|date|supplier|vendor)`)
	// Match a complete line: description + qty + price + total
	re := regexp.MustCompile(`^(.+?)\s+(\d+)\s+([\d.,]+)\s+([\d.,]+)\s*(?:[A-Z]{3})?$`)
	reNumOnly := regexp.MustCompile(`^[\d.,]+$`)

	var items []model.LineItem
	var currentItem *model.LineItem

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if line == "" {
			continue
		}

		if skip.MatchString(line) {
			continue
		}

		if matches := re.FindStringSubmatch(line); len(matches) == 5 {
			qty, _ := strconv.Atoi(matches[2])
			price, _ := strconv.ParseFloat(strings.ReplaceAll(matches[3], ",", ""), 64)
			total, _ := strconv.ParseFloat(strings.ReplaceAll(matches[4], ",", ""), 64)
			items = append(items, model.LineItem{
				Description: strings.TrimSpace(matches[1]),
				Quantity:    qty,
				Price:       price,
				Total:       total,
			})
			currentItem = nil
			continue
		}

		fields := strings.Fields(line)
		if len(fields) >= 4 {
			maybeQty := fields[len(fields)-3]
			maybePrice := strings.ReplaceAll(fields[len(fields)-2], ",", "")
			maybeTotal := strings.ReplaceAll(fields[len(fields)-1], ",", "")
			if qty, err1 := strconv.Atoi(maybeQty); err1 == nil {
				if price, err2 := strconv.ParseFloat(maybePrice, 64); err2 == nil {
					if total, err3 := strconv.ParseFloat(maybeTotal, 64); err3 == nil {
						desc := strings.Join(fields[:len(fields)-3], " ")
						items = append(items, model.LineItem{
							Description: strings.TrimSpace(desc),
							Quantity:    qty,
							Price:       price,
							Total:       total,
						})
						currentItem = nil
						continue
					}
				}
			}
		}

		if !reNumOnly.MatchString(line) {
			if currentItem == nil || currentItem.Description != "" {
				if currentItem != nil && currentItem.Quantity > 0 && currentItem.Price > 0 && currentItem.Total > 0 {
					items = append(items, *currentItem)
				}
				currentItem = &model.LineItem{Description: line}
			}
		} else if currentItem != nil {
			val, _ := strconv.ParseFloat(strings.ReplaceAll(line, ",", ""), 64)

			if intVal, err := strconv.Atoi(line); err == nil && intVal > 0 && intVal < 1000 && currentItem.Quantity == 0 {
				currentItem.Quantity = intVal
			} else if val > 0 {
				if currentItem.Price == 0 {
					currentItem.Price = val
				} else if currentItem.Total == 0 {
					currentItem.Total = val
					items = append(items, *currentItem)
					currentItem = nil
				}
			}
		}
	}

	if currentItem != nil && currentItem.Quantity > 0 && currentItem.Price > 0 && currentItem.Total > 0 {
		items = append(items, *currentItem)
	}

	return items
}

func (e *ExtractionServiceImpl) extractAmount(lines []string, field string) float64 {
	// Join relevant lines to handle multi-line content better
	joinedText := strings.Join(lines, " ")

	// Build regex patterns that look for the field followed by a number
	// Handle cases like "Subtotal 645", "Tax (20%) 129.0", "Total 800.0"
	patterns := []string{
		fmt.Sprintf(`(?i)%s\s*\([^)]*\)\s*([\d.]+)`, field), // Handles "Tax (20%) 129.0"
		fmt.Sprintf(`(?i)%s\s*[:=]\s*([\d.]+)`, field),      // Handles "Subtotal: 645"
		fmt.Sprintf(`(?i)%s\s+([\d.]+)`, field),             // Handles "Subtotal 645"
	}

	var results []float64

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		allMatches := re.FindAllStringSubmatch(joinedText, -1)
		for _, matches := range allMatches {
			if len(matches) > 1 {
				if val, err := strconv.ParseFloat(matches[1], 64); err == nil && val > 0 {
					results = append(results, val)
				}
			}
		}
	}

	if strings.ToLower(field) == "total" && len(results) > 0 {
		max := results[0]
		for _, val := range results {
			if val > max {
				max = val
			}
		}
		return max
	}

	if len(results) > 0 {
		return results[0]
	}

	labelRe := regexp.MustCompile(`(?i)` + field)
	numRe := regexp.MustCompile(`^([\d.]+)$`)

	for i, line := range lines {
		line = strings.TrimSpace(line)

		if labelRe.MatchString(line) {
			allNums := regexp.MustCompile(`([\d.]+)`).FindAllString(line, -1)
			for _, num := range allNums {
				if val, err := strconv.ParseFloat(num, 64); err == nil && val > 1 {
					return val
				}
			}

			for j := i + 1; j < len(lines) && j < i+5; j++ {
				nextLine := strings.TrimSpace(lines[j])
				if nextLine == "" {
					continue
				}

				if numRe.MatchString(nextLine) {
					if val, err := strconv.ParseFloat(nextLine, 64); err == nil && val > 1 {
						return val
					}
				}

				if regexp.MustCompile(`(?i)^(subtotal|tax|total|description|qty|quantity|unit|price|invoice|number|date|supplier|vendor)`).MatchString(nextLine) {
					break
				}
			}
		}
	}
	return 0
}

func (e *ExtractionServiceImpl) extractCurrency(lines []string) string {
	re := regexp.MustCompile(`\$|€|BAM|USD|EUR|BAM`)
	for _, line := range lines {
		if match := re.FindString(line); match != "" {
			if match == "$" {
				return "USD"
			}
			if match == "€" {
				return "EUR"
			}
			if match == "BAM" {
				return "BAM"
			}
			return match
		}
	}
	return "USD"
}
