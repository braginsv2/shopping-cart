package repository

import (
	"shopping-cart/internal/domain"
)

// CartRepository определяет методы для работы с корзинами
type CartRepository interface {
	Create(cart *domain.Cart) error
	GetByID(id uint) (*domain.Cart, error)
	GetByUserID(userID uint) (*domain.Cart, error)
	Update(cart *domain.Cart) error
	Delete(id uint) error
}

// CartItemRepository определяет методы для работы с элементами корзины
type CartItemRepository interface {
	Create(item *domain.CartItem) error
	GetByID(id uint) (*domain.CartItem, error)
	GetByCartID(cartID uint) ([]domain.CartItem, error)
	Update(item *domain.CartItem) error
	Delete(id uint) error
}

// OrderRepository определяет методы для работы с заказами
type OrderRepository interface {
	Create(order *domain.Order) error
	GetByID(id uint) (*domain.Order, error)
	GetByUserID(userID uint) ([]domain.Order, error)
	Update(order *domain.Order) error
	UpdateStatus(id uint, status string) error
	CreateOrderItem(item *domain.OrderItem) error
}

// ProductRepository определяет методы для работы с товарами
type ProductRepository interface {
	Create(product *domain.Product) error
	GetByID(id uint) (*domain.Product, error)
	GetAll() ([]domain.Product, error)
	Update(product *domain.Product) error
	Delete(id uint) error
}
