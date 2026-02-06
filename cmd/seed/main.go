package main

import (
	"log"

	"e-wallet-go/internal/config"
	"e-wallet-go/internal/handlers"
	"e-wallet-go/internal/middleware"
	"e-wallet-go/internal/repository"
	"e-wallet-go/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	db := config.ConnectDB()

	userRepo := repository.NewUserRepository(db)
	walletRepo := repository.NewWalletRepository(db)

	authService := services.NewAuthService(userRepo)
	userService := services.NewUserService(userRepo)
	walletService := services.NewWalletService(walletRepo)

	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService)
	walletHandler := handlers.NewWalletHandler(walletService)

	r := gin.Default()

	// Public
	api := r.Group("/api")
	{
		api.POST("/register", authHandler.Register)
		api.POST("/login", authHandler.Login)
	}

	// Protected
	protected := r.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	{
		// Profile
		protected.GET("/profile", userHandler.GetMyProfile)
		protected.PATCH("/profile", userHandler.UpdateMyProfile)

		// Wallet User
		protected.GET("/wallet", walletHandler.GetMyBalance)
		protected.POST("/wallet/deposit", walletHandler.UserTopUp)
		protected.POST("/wallet/withdraw", walletHandler.UserWithdraw)
		protected.GET("/transactions", walletHandler.GetTransactionHistory)

		// Admin Features
		protected.GET("/users", userHandler.GetAllUsers)
		protected.GET("/users/:id", userHandler.GetUserByID)
		protected.PATCH("/users/:id", userHandler.UpdateUser)
		protected.DELETE("/users/:id", userHandler.DeleteUser)

		protected.POST("/admin/users/:userID/topup", walletHandler.AdminTopUpUser)
		protected.POST("/admin/users/:userID/deduct", walletHandler.AdminDeductUser)
		protected.GET("/admin/transactions", walletHandler.GetAdminGlobalTransactions)
	}

	log.Println("Server running on port 8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to run server:", err)
	}
}
