// filters.go

package data

import (
	"errors"
	"fmt"
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

// ValidateSort checks if the sort parameter is valid.
func (f *Filters) ValidateSort() error {
	validSorts := map[string]bool{
		"rating": true,
		"date":   true,
	}

	if _, valid := validSorts[f.Sort]; !valid {
		return errors.New(fmt.Sprintf("invalid sort parameter: %s", f.Sort))
	}
	return nil
}

// BuildQuery appends sorting, limit, and offset to a base query.
func (f *Filters) BuildQuery(baseQuery string) string {
	// Apply default values to filter fields
	f.ValidateFilter()

	// Append sorting, limit, and offset clauses to the base query
	return fmt.Sprintf("%s ORDER BY %s DESC LIMIT %d OFFSET %d", baseQuery, f.SortColumn(), f.Limit, f.Offset)
}
