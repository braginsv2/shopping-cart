package main

import (
	"fmt"
	"log"
	"os"
	"shopping-cart/internal/delivery/http"
	"shopping-cart/internal/domain"
	repo "shopping-cart/internal/repository/postgres"
	"shopping-cart/internal/service/impl"

	_ "shopping-cart/docs" // Импортируем сгенерированную документацию

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// @title Shopping Cart API
// @version 1.0
// @description REST API для управления корзиной товаров в интернет-магазине
func main() {
	// Загрузка переменных окружения из .env файла
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Инициализация подключения к базе данных
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Автоматическая миграция схемы базы данных
	if err := db.AutoMigrate(
		&domain.Product{},
		&domain.Cart{},
		&domain.CartItem{},
		&domain.Order{},
		&domain.OrderItem{},
	); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Инициализация репозиториев
	cartRepo := repo.NewCartRepository(db)
	cartItemRepo := repo.NewCartItemRepository(db)
	orderRepo := repo.NewOrderRepository(db)
	productRepo := repo.NewProductRepository(db)

	// Инициализация сервисов
	cartService := impl.NewCartService(cartRepo, cartItemRepo, productRepo)
	orderService := impl.NewOrderService(orderRepo, cartRepo, cartItemRepo, productRepo)
	productService := impl.NewProductService(productRepo)

	// Инициализация HTTP-обработчика
	handler := http.NewHandler(cartService, orderService, productService)

	// Инициализация маршрутизатора Gin
	router := gin.Default()

	// Добавляем Swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Регистрация маршрутов API
	handler.RegisterRoutes(router)

	// Определение порта сервера
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	// Запуск HTTP-сервера
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
