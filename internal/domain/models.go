package domain

import (
	"time"

	"gorm.io/gorm"
)

// Product представляет товар в магазине
type Product struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Price       float64        `json:"price"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// Cart представляет корзину пользователя
type Cart struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	UserID    uint           `json:"user_id"`
	Items     []CartItem     `gorm:"foreignKey:CartID;constraint:OnDelete:CASCADE;" json:"items"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// CartItem представляет элемент корзины
type CartItem struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CartID    uint           `json:"cart_id"`
	ProductID uint           `json:"product_id"`
	Product   Product        `gorm:"foreignKey:ProductID" json:"product"`
	Quantity  int            `json:"quantity"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// OrderItem представляет элемент заказа
type OrderItem struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	OrderID   uint           `json:"order_id"`
	ProductID uint           `json:"product_id"`
	Product   Product        `gorm:"foreignKey:ProductID" json:"product"`
	Quantity  int            `json:"quantity"`
	Price     float64        `json:"price"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// Order представляет заказ пользователя
type Order struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	UserID    uint           `json:"user_id"`
	Status    string         `json:"status"`
	Items     []OrderItem    `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE;" json:"items"`
	Total     float64        `json:"total"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
