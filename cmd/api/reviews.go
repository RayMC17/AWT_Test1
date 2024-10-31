package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/RayMC17/AWT_Test1/internal/data"
	"github.com/RayMC17/AWT_Test1/internal/validator"
)

func (a *applicationDependencies) createReviewHandler(w http.ResponseWriter, r *http.Request) {
	productID, err := a.readIDParam(r)
	if err != nil {
		a.notFoundResponse(w, r, "")
		return
	}

	var input struct {
		Content string `json:"content"`
		Author  string `json:"author"`
		Rating  int    `json:"rating"`
	}

	err = a.readJSON(w, r, &input)
	if err != nil {
		a.badRequestResponse(w, r, err)
		return
	}

	review := &data.Review{
		ProductID: productID,
		Content:   input.Content,
		Author:    input.Author,
		Rating:    input.Rating,
	}

	v := validator.New()
	data.ValidateReview(v, review)
	if !v.IsEmpty() {
		a.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = a.reviewModel.Insert(review)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

	// Update product's average rating after creating a new review
	err = a.productModel.UpdateAverageRating(productID)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/products/%d/reviews/%d", review.ProductID, review.ID))
	err = a.writeJSON(w, http.StatusCreated, envelope{"review": review}, headers)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}

func (a *applicationDependencies) showReviewHandler(w http.ResponseWriter, r *http.Request) {
	id, err := a.readIDParam(r)
	if err != nil {
		a.notFoundResponse(w, r, "")
		return
	}

	review, err := a.reviewModel.Get(id)
	if err != nil {
		if err.Error() == "review not found" {
			a.notFoundResponse(w, r, "")
		} else {
			a.serverErrorResponse(w, r, err)
		}
		return
	}

	err = a.writeJSON(w, http.StatusOK, envelope{"review": review}, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}

func (a *applicationDependencies) updateReviewHandler(w http.ResponseWriter, r *http.Request) {
	id, err := a.readIDParam(r)
	if err != nil {
		a.notFoundResponse(w, r, "")
		return
	}

	review, err := a.reviewModel.Get(id)
	if err != nil {
		a.notFoundResponse(w, r, "")
		return
	}

	var input struct {
		Content *string `json:"content"`
		Author  *string `json:"author"`
		Rating  *int    `json:"rating"`
	}

	err = a.readJSON(w, r, &input)
	if err != nil {
		a.badRequestResponse(w, r, err)
		return
	}

	if input.Content != nil {
		review.Content = *input.Content
	}
	if input.Author != nil {
		review.Author = *input.Author
	}
	if input.Rating != nil {
		review.Rating = *input.Rating
	}

	v := validator.New()
	data.ValidateReview(v, review)
	if !v.IsEmpty() {
		a.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = a.reviewModel.Update(review)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

	// Update product's average rating after updating a review
	err = a.productModel.UpdateAverageRating(review.ProductID)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

	err = a.writeJSON(w, http.StatusOK, envelope{"review": review}, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}

func (a *applicationDependencies) deleteReviewHandler(w http.ResponseWriter, r *http.Request) {
	id, err := a.readIDParam(r)
	if err != nil {
		a.notFoundResponse(w, r, "")
		return
	}

	// Fetch the review before deleting to get ProductID for average rating update
	review, err := a.reviewModel.Get(id)
	if err != nil {
		a.notFoundResponse(w, r, "")
		return
	}

	err = a.reviewModel.Delete(id)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

	// Update product's average rating after deleting a review
	err = a.productModel.UpdateAverageRating(review.ProductID)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

	err = a.writeJSON(w, http.StatusOK, envelope{"message": "review successfully deleted"}, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}

func (a *applicationDependencies) listReviewsHandler(w http.ResponseWriter, r *http.Request) {
	productID, _ := strconv.ParseInt(r.URL.Query().Get("product_id"), 10, 64)

	// Initialize filters from query parameters
	filters := data.Filters{
		Sort:   r.URL.Query().Get("sort"),
		Limit:  parseInt(r.URL.Query().Get("limit"), 10),
		Offset: parseInt(r.URL.Query().Get("offset"), 0),
	}

	// Validate sort parameter
	if err := filters.ValidateSort(); err != nil {
		a.badRequestResponse(w, r, err)
		return
	}

	// Pass individual parameters instead of `filters`
	reviews, err := a.reviewModel.GetAll(productID, filters.Sort, filters.Limit, filters.Offset)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

	// Return a 404 response if no reviews were found
	if len(reviews) == 0 {
		a.notFoundResponse(w, r, "No reviews found matching the specified filters.")
		return
	}

	err = a.writeJSON(w, http.StatusOK, envelope{"reviews": reviews}, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}

// Helper function to parse integers with a default fallback
func parseInt(s string, defaultValue int) int {
	if i, err := strconv.Atoi(s); err == nil {
		return i
	}
	return defaultValue
}
