package data

import (
	//"errors"
	"fmt"
	"github.com/RayMC17/AWT_Test1/internal/validator"
)

// Filters struct holds the filter, sort, and pagination parameters.
type Filters struct {
	Sort   string
	Limit  int
	Offset int
}

// ValidateFilter ensures that the provided filter values are within acceptable bounds.
func (f *Filters) ValidateFilter() {
	// Set a default limit if none provided or if it's out of bounds
	if f.Limit <= 0 || f.Limit > 100 {
		f.Limit = 10 // Default limit
	}

	// Ensure offset is non-negative
	if f.Offset < 0 {
		f.Offset = 0
	}

	// Set a default sort if none is provided
	if f.Sort == "" {
		f.Sort = "created_at"
	}
}

// SortColumn returns a safe SQL column name for sorting based on the input sort parameter.
func (f *Filters) SortColumn() string {
	switch f.Sort {
	case "rating":
		return "average_rating"
	case "date":
		return "created_at"
	default:
		return "created_at"
	}
}

// ValidateSort checks if the sort parameter is valid using the Validator.
func (f *Filters) ValidateSort(v *validator.Validator) {
	validSorts := map[string]bool{
		"rating": true,
		"date":   true,
	}

	// Check if the Sort field is empty or invalid
	if f.Sort != "" && !validSorts[f.Sort] {
		v.AddError("sort", fmt.Sprintf("invalid sort parameter: %s", f.Sort))
	}
}

// BuildQuery appends sorting, limit, and offset to a base query.
func (f *Filters) BuildQuery(baseQuery string) string {
	// Apply default values to filter fields
	f.ValidateFilter()

	// Append sorting, limit, and offset clauses to the base query
	return fmt.Sprintf("%s ORDER BY %s DESC LIMIT %d OFFSET %d", baseQuery, f.SortColumn(), f.Limit, f.Offset)
}
