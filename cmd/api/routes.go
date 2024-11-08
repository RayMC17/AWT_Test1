package main

import (
    "net/http"
    "github.com/julienschmidt/httprouter"
)

func (a *applicationDependencies) routes() http.Handler {
    router := httprouter.New()

    // // User routes
    // router.HandlerFunc(http.MethodPost, "/v1/users", a.createUserHandler)
    // router.HandlerFunc(http.MethodGet, "/v1/users/:id", a.showUserHandler)
    // router.HandlerFunc(http.MethodPatch, "/v1/users/:id", a.updateUserHandler)
    // router.HandlerFunc(http.MethodDelete, "/v1/users/:id", a.deleteUserHandler)

    // Product routes
    router.HandlerFunc(http.MethodPost, "/v1/products", a.createProductHandler)
    router.HandlerFunc(http.MethodGet, "/v1/products/:id", a.showProductHandler)
    router.HandlerFunc(http.MethodPatch, "/v1/products/:id", a.updateProductHandler)
    router.HandlerFunc(http.MethodDelete, "/v1/products/:id", a.deleteProductHandler)
    router.HandlerFunc(http.MethodGet, "/v1/products", a.listProductsHandler)

    // Review routes
    router.HandlerFunc(http.MethodPost, "/v1/products/:id/reviews", a.createReviewHandler)
    router.HandlerFunc(http.MethodGet, "/v1/products/:id/reviews/:review_id", a.showReviewHandler)
    router.HandlerFunc(http.MethodPatch, "/v1/products/:id/reviews/:review_id", a.updateReviewHandler)
    router.HandlerFunc(http.MethodDelete, "/v1/products/:id/reviews/:review_id", a.deleteReviewHandler)
    router.HandlerFunc(http.MethodGet, "/v1/reviews", a.listReviewsHandler)
    router.HandlerFunc(http.MethodGet, "/v1/products/:id/reviews", a.listReviewsHandler)

//     return a.recoverPanic(router)
return a.recoverPanic(a.rateLimit(router))
}
