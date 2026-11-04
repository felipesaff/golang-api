package main

import (
	"log"
	"net/http"
	"time"

	"crud-api/config"
	"crud-api/handlers"
	"crud-api/repositories"
	"crud-api/services"
)

func main() {
	db := config.LoadDB()
	defer db.Close()

	// Ligação das camadas: repository -> service -> handler.
	userRepo := repositories.NewUserRepository(db)
	productRepo := repositories.NewProductRepository(db)
	orderRepo := repositories.NewOrderRepository(db)

	userService := services.NewUserService(userRepo)
	productService := services.NewProductService(productRepo)
	orderService := services.NewOrderService(orderRepo, userRepo, productRepo)

	mux := http.NewServeMux()
	handlers.RegisterHandlers(mux,
		handlers.NewUserHandler(userService),
		handlers.NewProductHandler(productService),
		handlers.NewOrderHandler(orderService),
	)

	srv := &http.Server{
		Addr:              ":8080",
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Println("Server is running on port 8080")
	log.Fatal(srv.ListenAndServe())
}
