package http

import (
	"net/http"
	"shopping-cart/internal/domain"
	"shopping-cart/internal/service"

	"github.com/gin-gonic/gin"
)

// Handler обрабатывает HTTP-запросы
type Handler struct {
	cartService    service.CartService
	orderService   service.OrderService
	productService service.ProductService
}

// NewHandler создает новый экземпляр HTTP-обработчика
func NewHandler(cartService service.CartService, orderService service.OrderService, productService service.ProductService) *Handler {
	return &Handler{
		cartService:    cartService,
		orderService:   orderService,
		productService: productService,
	}
}

// RegisterRoutes регистрирует маршруты API
func (h *Handler) RegisterRoutes(router *gin.Engine) {
	// Cart routes
	cart := router.Group("/api/cart")
	{
		cart.GET("/", h.GetCart)
		cart.POST("/items", h.AddItem)
		cart.DELETE("/items/:id", h.RemoveItem)
		cart.DELETE("/", h.ClearCart)
	}

	// Order routes
	orders := router.Group("/api/orders")
	{
		orders.POST("/", h.CreateOrder)
		orders.GET("/:id", h.GetOrder)
		orders.GET("/", h.GetUserOrders)
		orders.PATCH("/:id/status", h.UpdateOrderStatus)
	}

	// Product routes
	products := router.Group("/api/products")
	{
		products.POST("/", h.CreateProduct)
		products.GET("/:id", h.GetProduct)
		products.GET("/", h.GetAllProducts)
		products.PUT("/:id", h.UpdateProduct)
		products.DELETE("/:id", h.DeleteProduct)
	}
}

// @Summary Получить корзину пользователя
// @Description Возвращает содержимое корзины пользователя
// @Tags cart
// @Accept json
// @Produce json
// @Success 200 {object} domain.Cart
// @Failure 500 {object} map[string]string
// @Router /cart [get]
func (h *Handler) GetCart(c *gin.Context) {
	userID := uint(1) // TODO: Get from auth middleware
	cart, err := h.cartService.GetCart(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, cart)
}

// @Summary Добавить товар в корзину
// @Description Добавляет указанный товар в корзину пользователя
// @Tags cart
// @Accept json
// @Produce json
// @Param item body domain.CartItem true "Товар для добавления"
// @Success 200 {object} domain.Cart
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /cart/items [post]
func (h *Handler) AddItem(c *gin.Context) {
	var request struct {
		ProductID uint `json:"product_id" binding:"required"`
		Quantity  int  `json:"quantity" binding:"required,min=1"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := uint(1) // TODO: Get from auth middleware
	if err := h.cartService.AddItem(userID, request.ProductID, request.Quantity); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusCreated)
}

// @Summary Удалить товар из корзины
// @Description Удаляет указанный товар из корзины пользователя
// @Tags cart
// @Accept json
// @Produce json
// @Param item_id path int true "ID товара"
// @Success 200 {object} domain.Cart
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /cart/items/{item_id} [delete]
func (h *Handler) RemoveItem(c *gin.Context) {
	itemID := uint(1) // TODO: Parse from URL
	userID := uint(1) // TODO: Get from auth middleware
	if err := h.cartService.RemoveItem(userID, itemID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// ClearCart очищает корзину
func (h *Handler) ClearCart(c *gin.Context) {
	userID := uint(1) // TODO: Get from auth middleware
	if err := h.cartService.ClearCart(userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// @Summary Создать заказ
// @Description Создает новый заказ на основе содержимого корзины
// @Tags order
// @Accept json
// @Produce json
// @Success 201 {object} domain.Order
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /orders [post]
func (h *Handler) CreateOrder(c *gin.Context) {
	userID := uint(1) // TODO: Get from auth middleware
	order, err := h.orderService.CreateOrder(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, order)
}

// @Summary Получить заказ
// @Description Возвращает информацию о заказе по его ID
// @Tags order
// @Accept json
// @Produce json
// @Param id path int true "ID заказа"
// @Success 200 {object} domain.Order
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /orders/{id} [get]
func (h *Handler) GetOrder(c *gin.Context) {
	orderID := uint(1) // TODO: Parse from URL
	order, err := h.orderService.GetOrder(orderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, order)
}

// GetUserOrders возвращает список заказов пользователя
func (h *Handler) GetUserOrders(c *gin.Context) {
	userID := uint(1) // TODO: Get from auth middleware
	orders, err := h.orderService.GetUserOrders(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, orders)
}

// UpdateOrderStatus обновляет статус заказа
func (h *Handler) UpdateOrderStatus(c *gin.Context) {
	var request struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	orderID := uint(1) // TODO: Parse from URL
	if err := h.orderService.UpdateOrderStatus(orderID, request.Status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

// @Summary Создать товар
// @Description Создает новый товар в магазине
// @Tags product
// @Accept json
// @Produce json
// @Param product body domain.Product true "Товар для создания"
// @Success 201 {object} domain.Product
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /products [post]
func (h *Handler) CreateProduct(c *gin.Context) {
	var product domain.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.productService.CreateProduct(&product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, product)
}

// @Summary Получить товар
// @Description Возвращает информацию о товаре по его ID
// @Tags product
// @Accept json
// @Produce json
// @Param id path int true "ID товара"
// @Success 200 {object} domain.Product
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /products/{id} [get]
func (h *Handler) GetProduct(c *gin.Context) {
	productID := uint(1) // TODO: Parse from URL
	product, err := h.productService.GetProduct(productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, product)
}

// GetAllProducts возвращает список всех товаров
func (h *Handler) GetAllProducts(c *gin.Context) {
	products, err := h.productService.GetAllProducts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, products)
}

// UpdateProduct обновляет информацию о товаре
func (h *Handler) UpdateProduct(c *gin.Context) {
	var product domain.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.productService.UpdateProduct(&product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, product)
}

// DeleteProduct удаляет товар
func (h *Handler) DeleteProduct(c *gin.Context) {
	productID := uint(1) // TODO: Parse from URL
	if err := h.productService.DeleteProduct(productID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
