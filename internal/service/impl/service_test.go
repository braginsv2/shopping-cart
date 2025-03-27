package impl

import (
	"shopping-cart/internal/domain"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockCartRepository - мок репозитория корзины
type MockCartRepository struct {
	mock.Mock
}

func (m *MockCartRepository) Create(cart *domain.Cart) error {
	args := m.Called(cart)
	return args.Error(0)
}

func (m *MockCartRepository) GetByID(id uint) (*domain.Cart, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Cart), args.Error(1)
}

func (m *MockCartRepository) GetByUserID(userID uint) (*domain.Cart, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Cart), args.Error(1)
}

func (m *MockCartRepository) Update(cart *domain.Cart) error {
	args := m.Called(cart)
	return args.Error(0)
}

func (m *MockCartRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

// MockCartItemRepository - мок репозитория элементов корзины
type MockCartItemRepository struct {
	mock.Mock
}

func (m *MockCartItemRepository) Create(item *domain.CartItem) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *MockCartItemRepository) GetByID(id uint) (*domain.CartItem, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.CartItem), args.Error(1)
}

func (m *MockCartItemRepository) GetByCartID(cartID uint) ([]domain.CartItem, error) {
	args := m.Called(cartID)
	return args.Get(0).([]domain.CartItem), args.Error(1)
}

func (m *MockCartItemRepository) Update(item *domain.CartItem) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *MockCartItemRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

// MockProductRepository - мок репозитория товаров
type MockProductRepository struct {
	mock.Mock
}

func (m *MockProductRepository) Create(product *domain.Product) error {
	args := m.Called(product)
	return args.Error(0)
}

func (m *MockProductRepository) GetByID(id uint) (*domain.Product, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Product), args.Error(1)
}

func (m *MockProductRepository) GetAll() ([]domain.Product, error) {
	args := m.Called()
	return args.Get(0).([]domain.Product), args.Error(1)
}

func (m *MockProductRepository) Update(product *domain.Product) error {
	args := m.Called(product)
	return args.Error(0)
}

func (m *MockProductRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

// Тесты для CartService
func TestAddItem(t *testing.T) {
	mockCartRepo := new(MockCartRepository)
	mockCartItemRepo := new(MockCartItemRepository)
	mockProductRepo := new(MockProductRepository)

	service := NewCartService(mockCartRepo, mockCartItemRepo, mockProductRepo)

	tests := []struct {
		name          string
		userID        uint
		productID     uint
		quantity      int
		setupMocks    func()
		expectedError error
	}{
		{
			name:      "Успешное добавление товара",
			userID:    1,
			productID: 1,
			quantity:  2,
			setupMocks: func() {
				mockCartRepo.On("GetByUserID", uint(1)).Return(&domain.Cart{ID: 1, UserID: 1}, nil)
				mockProductRepo.On("GetByID", uint(1)).Return(&domain.Product{ID: 1, Price: 100}, nil)
				mockCartItemRepo.On("GetByCartID", uint(1)).Return([]domain.CartItem{}, nil)
				mockCartItemRepo.On("Create", mock.AnythingOfType("*domain.CartItem")).Return(nil)
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()
			err := service.AddItem(tt.userID, tt.productID, tt.quantity)
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
				return
			}
			assert.NoError(t, err)
			mockCartRepo.AssertExpectations(t)
			mockCartItemRepo.AssertExpectations(t)
			mockProductRepo.AssertExpectations(t)
		})
	}
}

// Тесты для ProductService
func TestCreateProduct(t *testing.T) {
	mockProductRepo := new(MockProductRepository)
	service := NewProductService(mockProductRepo)

	tests := []struct {
		name          string
		product       *domain.Product
		setupMocks    func()
		expectedError error
	}{
		{
			name: "Успешное создание товара",
			product: &domain.Product{
				Name:        "Test Product",
				Description: "Test Description",
				Price:       100,
			},
			setupMocks: func() {
				mockProductRepo.On("Create", mock.AnythingOfType("*domain.Product")).Return(nil)
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()
			err := service.CreateProduct(tt.product)
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
				return
			}
			assert.NoError(t, err)
			mockProductRepo.AssertExpectations(t)
		})
	}
}
