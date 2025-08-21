package service

import (
	"shopping-cart/internal/domain"
)

type CartService interface {
	AddItem(userID uint, productID uint, quantity int) error
	RemoveItem(userID uint, itemID uint) error
	GetCart(userID uint) (*domain.Cart, error)
	ClearCart(userID uint) error
}

type OrderService interface {
	CreateOrder(userID uint) (*domain.Order, error)
	GetOrder(orderID uint) (*domain.Order, error)
	GetUserOrders(userID uint) ([]domain.Order, error)
	UpdateOrderStatus(orderID uint, status string) error
}

type ProductService interface {
	CreateProduct(product *domain.Product) error
	GetProduct(id uint) (*domain.Product, error)
	GetAllProducts() ([]domain.Product, error)
	UpdateProduct(product *domain.Product) error
	DeleteProduct(id uint) error
} 


