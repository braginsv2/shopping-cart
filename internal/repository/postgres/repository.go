package postgres

import (
	"shopping-cart/internal/domain"
	"shopping-cart/internal/repository"

	"gorm.io/gorm"
)

type cartRepository struct {
	db *gorm.DB
}

type cartItemRepository struct {
	db *gorm.DB
}

type orderRepository struct {
	db *gorm.DB
}

type productRepository struct {
	db *gorm.DB
}

func NewCartRepository(db *gorm.DB) repository.CartRepository {
	return &cartRepository{db: db}
}

func NewCartItemRepository(db *gorm.DB) repository.CartItemRepository {
	return &cartItemRepository{db: db}
}

func NewOrderRepository(db *gorm.DB) repository.OrderRepository {
	return &orderRepository{db: db}
}

func NewProductRepository(db *gorm.DB) repository.ProductRepository {
	return &productRepository{db: db}
}

// Cart Repository Implementation
func (r *cartRepository) Create(cart *domain.Cart) error {
	return r.db.Create(cart).Error
}

func (r *cartRepository) GetByID(id uint) (*domain.Cart, error) {
	var cart domain.Cart
	err := r.db.Preload("Items.Product").First(&cart, id).Error
	return &cart, err
}

func (r *cartRepository) GetByUserID(userID uint) (*domain.Cart, error) {
	var cart domain.Cart
	err := r.db.Preload("Items.Product").Where("user_id = ?", userID).First(&cart).Error
	return &cart, err
}

func (r *cartRepository) Update(cart *domain.Cart) error {
	return r.db.Save(cart).Error
}

func (r *cartRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Cart{}, id).Error
}

// CartItem Repository Implementation
func (r *cartItemRepository) Create(item *domain.CartItem) error {
	return r.db.Create(item).Error
}

func (r *cartItemRepository) GetByID(id uint) (*domain.CartItem, error) {
	var item domain.CartItem
	err := r.db.Preload("Product").First(&item, id).Error
	return &item, err
}

func (r *cartItemRepository) GetByCartID(cartID uint) ([]domain.CartItem, error) {
	var items []domain.CartItem
	err := r.db.Preload("Product").Where("cart_id = ?", cartID).Find(&items).Error
	return items, err
}

func (r *cartItemRepository) Update(item *domain.CartItem) error {
	return r.db.Save(item).Error
}

func (r *cartItemRepository) Delete(id uint) error {
	return r.db.Delete(&domain.CartItem{}, id).Error
}

// Order Repository Implementation
func (r *orderRepository) Create(order *domain.Order) error {
	return r.db.Create(order).Error
}

func (r *orderRepository) CreateOrderItem(item *domain.OrderItem) error {
	return r.db.Create(item).Error
}

func (r *orderRepository) GetByID(id uint) (*domain.Order, error) {
	var order domain.Order
	err := r.db.Preload("Items.Product").First(&order, id).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *orderRepository) GetByUserID(userID uint) ([]domain.Order, error) {
	var orders []domain.Order
	err := r.db.Preload("Items.Product").Where("user_id = ?", userID).Find(&orders).Error
	return orders, err
}

func (r *orderRepository) Update(order *domain.Order) error {
	return r.db.Save(order).Error
}

func (r *orderRepository) UpdateStatus(id uint, status string) error {
	return r.db.Model(&domain.Order{}).Where("id = ?", id).Update("status", status).Error
}

// Product Repository Implementation
func (r *productRepository) Create(product *domain.Product) error {
	return r.db.Create(product).Error
}

func (r *productRepository) GetByID(id uint) (*domain.Product, error) {
	var product domain.Product
	err := r.db.First(&product, id).Error
	return &product, err
}

func (r *productRepository) GetAll() ([]domain.Product, error) {
	var products []domain.Product
	err := r.db.Find(&products).Error
	return products, err
}

func (r *productRepository) Update(product *domain.Product) error {
	return r.db.Save(product).Error
}

func (r *productRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Product{}, id).Error
}
