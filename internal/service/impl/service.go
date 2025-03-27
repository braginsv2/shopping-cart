package impl

import (
	"errors"
	"shopping-cart/internal/domain"
	"shopping-cart/internal/repository"
	"shopping-cart/internal/service"
)

// cartService реализует интерфейс CartService
type cartService struct {
	cartRepo     repository.CartRepository
	cartItemRepo repository.CartItemRepository
	productRepo  repository.ProductRepository
}

// orderService реализует интерфейс OrderService
type orderService struct {
	orderRepo    repository.OrderRepository
	cartRepo     repository.CartRepository
	cartItemRepo repository.CartItemRepository
	productRepo  repository.ProductRepository
}

// productService реализует интерфейс ProductService
type productService struct {
	productRepo repository.ProductRepository
}

// NewCartService создает новый экземпляр CartService
func NewCartService(cartRepo repository.CartRepository, cartItemRepo repository.CartItemRepository, productRepo repository.ProductRepository) service.CartService {
	return &cartService{
		cartRepo:     cartRepo,
		cartItemRepo: cartItemRepo,
		productRepo:  productRepo,
	}
}

// NewOrderService создает новый экземпляр OrderService
func NewOrderService(orderRepo repository.OrderRepository, cartRepo repository.CartRepository, cartItemRepo repository.CartItemRepository, productRepo repository.ProductRepository) service.OrderService {
	return &orderService{
		orderRepo:    orderRepo,
		cartRepo:     cartRepo,
		cartItemRepo: cartItemRepo,
		productRepo:  productRepo,
	}
}

// NewProductService создает новый экземпляр ProductService
func NewProductService(productRepo repository.ProductRepository) service.ProductService {
	return &productService{
		productRepo: productRepo,
	}
}

// AddItem добавляет товар в корзину пользователя
// Если товар уже есть в корзине, увеличивает его количество
func (s *cartService) AddItem(userID uint, productID uint, quantity int) error {
	// Get or create cart
	cart, err := s.cartRepo.GetByUserID(userID)
	if err != nil {
		cart = &domain.Cart{UserID: userID}
		if err := s.cartRepo.Create(cart); err != nil {
			return err
		}
	}

	// Check if product exists
	if _, err := s.productRepo.GetByID(productID); err != nil {
		return err
	}

	// Check if item already exists in cart
	items, err := s.cartItemRepo.GetByCartID(cart.ID)
	if err != nil {
		return err
	}

	for _, item := range items {
		if item.ProductID == productID {
			// Update quantity of existing item
			item.Quantity += quantity
			return s.cartItemRepo.Update(&item)
		}
	}

	// Create new cart item if not exists
	cartItem := &domain.CartItem{
		CartID:    cart.ID,
		ProductID: productID,
		Quantity:  quantity,
	}

	return s.cartItemRepo.Create(cartItem)
}

// RemoveItem удаляет товар из корзины пользователя
func (s *cartService) RemoveItem(userID uint, itemID uint) error {
	cart, err := s.cartRepo.GetByUserID(userID)
	if err != nil {
		return err
	}

	item, err := s.cartItemRepo.GetByID(itemID)
	if err != nil {
		return err
	}

	if item.CartID != cart.ID {
		return errors.New("item does not belong to user's cart")
	}

	return s.cartItemRepo.Delete(itemID)
}

// GetCart возвращает корзину пользователя
// Если корзина не существует, создает новую
func (s *cartService) GetCart(userID uint) (*domain.Cart, error) {
	cart, err := s.cartRepo.GetByUserID(userID)
	if err != nil {
		// Если корзина не найдена, создаем новую
		cart = &domain.Cart{UserID: userID}
		if err := s.cartRepo.Create(cart); err != nil {
			return nil, err
		}
	}
	return cart, nil
}

// ClearCart очищает корзину пользователя
func (s *cartService) ClearCart(userID uint) error {
	cart, err := s.cartRepo.GetByUserID(userID)
	if err != nil {
		return err
	}

	items, err := s.cartItemRepo.GetByCartID(cart.ID)
	if err != nil {
		return err
	}

	for _, item := range items {
		if err := s.cartItemRepo.Delete(item.ID); err != nil {
			return err
		}
	}

	return nil
}

// CreateOrder создает новый заказ из корзины пользователя
// После создания заказа корзина очищается
func (s *orderService) CreateOrder(userID uint) (*domain.Order, error) {
	cart, err := s.cartRepo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	cartItems, err := s.cartItemRepo.GetByCartID(cart.ID)
	if err != nil {
		return nil, err
	}

	if len(cartItems) == 0 {
		return nil, errors.New("cart is empty")
	}

	// Calculate total
	var total float64
	for _, item := range cartItems {
		product, err := s.productRepo.GetByID(item.ProductID)
		if err != nil {
			return nil, err
		}
		total += float64(item.Quantity) * product.Price
	}

	// Create order
	order := &domain.Order{
		UserID: userID,
		Status: "pending",
		Total:  total,
	}

	if err := s.orderRepo.Create(order); err != nil {
		return nil, err
	}

	// Create order items
	for _, item := range cartItems {
		product, err := s.productRepo.GetByID(item.ProductID)
		if err != nil {
			return nil, err
		}

		orderItem := &domain.OrderItem{
			OrderID:   order.ID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     product.Price,
		}

		if err := s.orderRepo.CreateOrderItem(orderItem); err != nil {
			return nil, err
		}
	}

	// Clear cart after order creation
	if err := s.cartRepo.Delete(cart.ID); err != nil {
		return nil, err
	}

	// Get the complete order with items
	return s.orderRepo.GetByID(order.ID)
}

// GetOrder возвращает заказ по его ID
func (s *orderService) GetOrder(orderID uint) (*domain.Order, error) {
	return s.orderRepo.GetByID(orderID)
}

// GetUserOrders возвращает все заказы пользователя
func (s *orderService) GetUserOrders(userID uint) ([]domain.Order, error) {
	return s.orderRepo.GetByUserID(userID)
}

// UpdateOrderStatus обновляет статус заказа
func (s *orderService) UpdateOrderStatus(orderID uint, status string) error {
	return s.orderRepo.UpdateStatus(orderID, status)
}

// CreateProduct создает новый товар
func (s *productService) CreateProduct(product *domain.Product) error {
	return s.productRepo.Create(product)
}

// GetProduct возвращает товар по его ID
func (s *productService) GetProduct(id uint) (*domain.Product, error) {
	return s.productRepo.GetByID(id)
}

// GetAllProducts возвращает все товары
func (s *productService) GetAllProducts() ([]domain.Product, error) {
	return s.productRepo.GetAll()
}

// UpdateProduct обновляет информацию о товаре
func (s *productService) UpdateProduct(product *domain.Product) error {
	return s.productRepo.Update(product)
}

// DeleteProduct удаляет товар
func (s *productService) DeleteProduct(id uint) error {
	return s.productRepo.Delete(id)
}
