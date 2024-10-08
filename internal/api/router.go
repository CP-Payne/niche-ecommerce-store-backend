package api

import (
	"net/http"

	"github.com/CP-Payne/ecomstore/internal/api/handlers"
	"github.com/CP-Payne/ecomstore/internal/config"
	"github.com/CP-Payne/ecomstore/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth"
	"go.uber.org/zap"

	cmid "github.com/CP-Payne/ecomstore/internal/api/middleware"
)

func SetupRouter(cfg *config.Config) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Use(cmid.CorsMiddleware)

	paypalProcessor, err := service.NewPayPalProcessor(cfg.PaymentProcessor)
	if err != nil {
		cfg.Logger.Fatal("failed to setup router", zap.Error(err))
	}

	userSrv := service.NewUserService(cfg.DB)
	productSrv := service.NewProductService(cfg.DB)
	reviewSrv := service.NewReviewService(cfg.DB)
	cartSrv := service.NewCartService(cfg.DB)
	orderSrv := service.NewOrderService(cfg.DB, cfg.SqlDB)
	paymentSrv := service.NewPaymentService(cfg.DB, paypalProcessor, orderSrv, productSrv, cartSrv)

	authHandler := handlers.NewAuthHandler(userSrv)
	productHandler := handlers.NewProductHandler(productSrv)
	reviewHander := handlers.NewReviewHandler(reviewSrv, productSrv)
	cartHandler := handlers.NewCartHandler(cartSrv, productSrv)
	userHandler := handlers.NewUserHandler(userSrv)
	paymentHandler := handlers.NewPaymentHandler(productSrv, paymentSrv, cartSrv, orderSrv)
	orderHandler := handlers.NewOrderHandler(orderSrv)

	r.Group(func(r chi.Router) {
		r.Post("/register", authHandler.RegisterUser)
		r.Post("/login", authHandler.LoginUser)

		r.Get("/products", productHandler.GetAllProducts)
		r.Get("/products/{id}", productHandler.GetProduct)

		r.Get("/payment/capture-order", paymentHandler.CaptureOrder)

		r.Get("/products/categories", productHandler.GetProductCategories)
		r.Get("/products/categories/{id}", productHandler.GetProductsByCategory)

	})

	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(config.GetTokenAuth()))
		r.Use(jwtauth.Authenticator)

		r.Get("/user/profile", userHandler.GetUserDetails)
		r.Get("/user/orders", orderHandler.GetUserOrders)

		r.Post("/payment/create-order/product", paymentHandler.CreateOrderProduct)
		r.Post("/payment/create-order/cart", paymentHandler.CreateOrderCart)

		r.Get("/cart", cartHandler.GetCart)
		r.Post("/cart/add", cartHandler.AddToCart)
		r.Post("/cart/remove", cartHandler.RemoveFromCart)
		r.Post("/cart/reduce", cartHandler.ReduceFromCart)

	})

	r.Group(func(r chi.Router) {
		r.Use(cmid.ProductMiddleware(productSrv, cfg.Logger))
		r.Get("/products/{id}/reviews", reviewHander.GetProductReviews)
	})

	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(config.GetTokenAuth()))
		r.Use(jwtauth.Authenticator)
		r.Use(cmid.ProductMiddleware(productSrv, cfg.Logger))

		r.Get("/products/{id}/reviews/user", reviewHander.GetUserReviewForProduct)
		r.Patch("/products/{id}/reviews", reviewHander.UpdateUserReview)
		r.Delete("/products/{id}/reviews", reviewHander.DeleteReview)
		r.Post("/products/{id}/reviews", reviewHander.AddReview)
	})
	r.Get("/home", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("This will be the home page"))
	}))

	return r
}
