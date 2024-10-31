// internal/data/product.go
package data

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/RayMC17/AWT_Test1/internal/validator"
)

type Product struct {
	ID            int64     `json:"id"`
	Name          string    `json:"name"`
	Description   string    `json:"description,omitempty"`
	Category      string    `json:"category"`
	ImageURL      string    `json:"image_url"`
	AverageRating float32   `json:"average_rating"`
	CreatedAt     time.Time `json:"-"`
	UpdatedAt     time.Time `json:"-"`
}

type ProductModel struct {
	DB *sql.DB
}

func ValidateProduct(v *validator.Validator, product *Product) {
	v.Check(product.Name != "", "name", "must be provided")
	v.Check(len(product.Name) <= 100, "name", "must not be more than 100 characters")
	v.Check(product.Category != "", "category", "must be provided")
	v.Check(product.ImageURL != "", "image_url", "must be a valid URL")
}

// Insert adds a new product to the database.
func (m ProductModel) Insert(product *Product) error {
	query := `
        INSERT INTO products (name, description, category, image_url)
        VALUES ($1, $2, $3, $4)
        RETURNING id, created_at, updated_at`

	args := []interface{}{product.Name, product.Description, product.Category, product.ImageURL}

	return m.DB.QueryRow(query, args...).Scan(&product.ID, &product.CreatedAt, &product.UpdatedAt)
}

// Get retrieves a specific product by ID.
func (m ProductModel) Get(id int64) (*Product, error) {
	query := `
        SELECT id, name, description, category, image_url, average_rating, created_at, updated_at
        FROM products
        WHERE id = $1`

	var product Product
	err := m.DB.QueryRow(query, id).Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.Category,
		&product.ImageURL,
		&product.AverageRating,
		&product.CreatedAt,
		&product.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("product not found")
	} else if err != nil {
		return nil, err
	}

	return &product, nil
}

// Update modifies an existing product's information in the database.
func (m ProductModel) Update(product *Product) error {
	query := `
        UPDATE products
        SET name = $1, description = $2, category = $3, image_url = $4, updated_at = NOW()
        WHERE id = $5`

	args := []interface{}{product.Name, product.Description, product.Category, product.ImageURL, product.ID}
	_, err := m.DB.Exec(query, args...)
	return err
}

// Delete removes a product by ID from the database.
func (m ProductModel) Delete(id int64) error {
	query := `
        DELETE FROM products
        WHERE id = $1`

	_, err := m.DB.Exec(query, id)
	return err
}

// GetAll retrieves all products with optional filtering, sorting, and pagination.

// internal/data/product.go

func (m ProductModel) GetAll(name string, category string, filters Filters) ([]*Product, error) {
    baseQuery := `
        SELECT id, name, description, category, image_url, average_rating, created_at, updated_at
        FROM products
        WHERE ($1 = '%%' OR LOWER(name) LIKE LOWER($1))
          AND ($2 = '' OR category = $2)
    `
    
    // Use the BuildQuery method to add sorting, limit, and offset to the query
    query := filters.BuildQuery(baseQuery)

    args := []interface{}{"%" + name + "%", category}

    rows, err := m.DB.Query(query, args...)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var products []*Product
    for rows.Next() {
        var product Product
        err := rows.Scan(
            &product.ID,
            &product.Name,
            &product.Description,
            &product.Category,
            &product.ImageURL,
            &product.AverageRating,
            &product.CreatedAt,
            &product.UpdatedAt,
        )
        if err != nil {
            return nil, err
        }
        products = append(products, &product)
    }

    return products, nil
}

// UpdateAverageRating recalculates the average rating for a product based on its reviews.
func (m ProductModel) UpdateAverageRating(productID int64) error {
	query := `
        UPDATE products
        SET average_rating = (
            SELECT COALESCE(AVG(rating), 0)
            FROM reviews
            WHERE product_id = $1
        )
        WHERE id = $1`

	_, err := m.DB.Exec(query, productID)
	return err
}
