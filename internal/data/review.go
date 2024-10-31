// internal/data/review.go
package data

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/RayMC17/AWT_Test1/internal/validator"
)

type Review struct {
	ID           int64     `json:"id"`
	ProductID    int64     `json:"product_id"`
	Content      string    `json:"content"`
	Author       string    `json:"author"`
	Rating       int       `json:"rating"`
	HelpfulCount int       `json:"helpful_count"`
	CreatedAt    time.Time `json:"-"`
	UpdatedAt    time.Time `json:"-"`
}

type ReviewModel struct {
	DB *sql.DB
}

func ValidateReview(v *validator.Validator, review *Review) {
	v.Check(review.Rating >= 1 && review.Rating <= 5, "rating", "must be between 1 and 5")
	v.Check(review.Content != "", "content", "must be provided")
	v.Check(review.Author != "", "author", "must be provided")
}

// Insert adds a new review to the database.
func (m ReviewModel) Insert(review *Review) error {
	query := `
        INSERT INTO reviews (product_id, content, author, rating)
        VALUES ($1, $2, $3, $4)
        RETURNING id, created_at, updated_at`

	args := []interface{}{review.ProductID, review.Content, review.Author, review.Rating}

	return m.DB.QueryRow(query, args...).Scan(&review.ID, &review.CreatedAt, &review.UpdatedAt)
}

// Get retrieves a specific review by ID.
func (m ReviewModel) Get(id int64) (*Review, error) {
	query := `
        SELECT id, product_id, content, author, rating, helpful_count, created_at, updated_at
        FROM reviews
        WHERE id = $1`

	var review Review
	err := m.DB.QueryRow(query, id).Scan(
		&review.ID,
		&review.ProductID,
		&review.Content,
		&review.Author,
		&review.Rating,
		&review.HelpfulCount,
		&review.CreatedAt,
		&review.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("review not found")
	} else if err != nil {
		return nil, err
	}

	return &review, nil
}

// Update modifies an existing review in the database.
func (m ReviewModel) Update(review *Review) error {
	query := `
        UPDATE reviews
        SET content = $1, author = $2, rating = $3, updated_at = NOW()
        WHERE id = $4`

	args := []interface{}{review.Content, review.Author, review.Rating, review.ID}
	_, err := m.DB.Exec(query, args...)
	return err
}

// Delete removes a review by ID from the database.
func (m ReviewModel) Delete(id int64) error {
	query := `
        DELETE FROM reviews
        WHERE id = $1`

	_, err := m.DB.Exec(query, id)
	return err
}

// GetAll retrieves all reviews with optional filtering, sorting, and pagination.
func (m ReviewModel) GetAll(productID int64, sort string, limit int, offset int) ([]*Review, error) {
	query := `
        SELECT id, product_id, content, author, rating, helpful_count, created_at, updated_at
        FROM reviews
        WHERE (product_id = $1 OR $1 = 0)
        ORDER BY CASE WHEN $2 = 'helpful' THEN helpful_count END DESC,
                 CASE WHEN $2 = 'date' THEN created_at END DESC
        LIMIT $3 OFFSET $4`

	args := []interface{}{productID, sort, limit, offset}

	rows, err := m.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviews []*Review
	for rows.Next() {
		var review Review
		err := rows.Scan(
			&review.ID,
			&review.ProductID,
			&review.Content,
			&review.Author,
			&review.Rating,
			&review.HelpfulCount,
			&review.CreatedAt,
			&review.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		reviews = append(reviews, &review)
	}

	return reviews, nil
}
