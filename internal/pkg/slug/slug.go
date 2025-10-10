package slug

import (
	"regexp"
	"strings"
)

// GenerateSlug creates a URL-friendly slug from a string.
func GenerateSlug(text string) string {
	// Convert to lowercase
	slug := strings.ToLower(text)

	// Replace non-alphanumeric characters with a hyphen
	// This regex will replace any sequence of non-alphanumeric chars with a single hyphen
	reg := regexp.MustCompile(`[^a-z0-9]+`)
	slug = reg.ReplaceAllString(slug, "-")

	// Trim leading/trailing hyphens that might have been created
	slug = strings.Trim(slug, "-")

	return slug
}
