package main

import (
	"fmt"
	"net/http"
	//"strconv"

	"github.com/RayMC17/AWT_Test1/internal/data"
	"github.com/RayMC17/AWT_Test1/internal/validator"
)

func (a *applicationDependencies) createProductHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Category    string `json:"category"`
		ImageURL    string `json:"image_url"`
	}

	err := a.readJSON(w, r, &input)
	if err != nil {
		a.badRequestResponse(w, r, err)
		return
	}

	product := &data.Product{
		Name:        input.Name,
		Description: input.Description,
		Category:    input.Category,
		ImageURL:    input.ImageURL,
	}

	v := validator.New()
	data.ValidateProduct(v, product)
	if !v.IsEmpty() {
		a.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = a.productModel.Insert(product)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/products/%d", product.ID))
	err = a.writeJSON(w, http.StatusCreated, envelope{"product": product}, headers)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}

func (a *applicationDependencies) listProductsHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	category := r.URL.Query().Get("category")

	// Initialize filters from query parameters
	filters := data.Filters{
		Sort:   r.URL.Query().Get("sort"),
		Limit:  parseInt(r.URL.Query().Get("limit"), 10),
		Offset: parseInt(r.URL.Query().Get("offset"), 0),
	}

	// Check if the sort parameter is valid
	if err := filters.ValidateSort(); err != nil {
		a.badRequestResponse(w, r, err)
		return
	}

	// Retrieve products based on filters
	products, err := a.productModel.GetAll(name, category, filters)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

	// Return a 404 response if no products were found
	if len(products) == 0 {
		a.notFoundResponse(w, r, "No products found matching the specified filters.")
		return
	}

	// Return the list of products
	err = a.writeJSON(w, http.StatusOK, envelope{"products": products}, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}





func (a *applicationDependencies) showProductHandler(w http.ResponseWriter, r *http.Request) {
	id, err := a.readIDParam(r)
	if err != nil {
		a.notFoundResponse(w, r, "")
		return
	}

	product, err := a.productModel.Get(id)
	if err != nil {
		if err.Error() == "product not found" {
			a.notFoundResponse(w, r, "")
		} else {
			a.serverErrorResponse(w, r, err)
		}
		return
	}

	err = a.writeJSON(w, http.StatusOK, envelope{"product": product}, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}

func (a *applicationDependencies) updateProductHandler(w http.ResponseWriter, r *http.Request) {
	id, err := a.readIDParam(r)
	if err != nil {
		a.notFoundResponse(w, r, "")
		return
	}

	product, err := a.productModel.Get(id)
	if err != nil {
		a.notFoundResponse(w, r, "")
		return
	}

	var input struct {
		Name        *string `json:"name"`
		Description *string `json:"description"`
		Category    *string `json:"category"`
		ImageURL    *string `json:"image_url"`
	}

	err = a.readJSON(w, r, &input)
	if err != nil {
		a.badRequestResponse(w, r, err)
		return
	}

	if input.Name != nil {
		product.Name = *input.Name
	}
	if input.Description != nil {
		product.Description = *input.Description
	}
	if input.Category != nil {
		product.Category = *input.Category
	}
	if input.ImageURL != nil {
		product.ImageURL = *input.ImageURL
	}

	v := validator.New()
	data.ValidateProduct(v, product)
	if !v.IsEmpty() {
		a.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = a.productModel.Update(product)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

	err = a.writeJSON(w, http.StatusOK, envelope{"product": product}, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}

func (a *applicationDependencies) deleteProductHandler(w http.ResponseWriter, r *http.Request) {
	id, err := a.readIDParam(r)
	if err != nil {
		a.notFoundResponse(w, r, "")
		return
	}

	err = a.productModel.Delete(id)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

	err = a.writeJSON(w, http.StatusOK, envelope{"message": "product successfully deleted"}, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}
