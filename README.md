# Shopping Cart API

REST API для управления корзиной товаров в интернет-магазине, реализованное на Go с использованием Gin и PostgreSQL.

## Требования

- Go 1.21 или выше
- Docker и Docker Compose
- PostgreSQL (если запускаете без Docker)

## Установка и запуск

1. Клонируйте репозиторий:
```bash
git clone https://github.com/braginsv2/shopping-cart
cd shopping-cart
```

2. Создайте файл `.env` в корневой директории проекта со следующим содержимым:
```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=shopping_cart
SERVER_PORT=8081
```

3. Запустите PostgreSQL через Docker Compose:
```bash
docker-compose up -d
```

4. Запустите приложение:
```bash
go run cmd/main.go
```

## API Endpoints

### Товары
- `GET /api/products` - получить список всех товаров
- `GET /api/products/:id` - получить информацию о товаре
- `POST /api/products` - создать новый товар
- `PUT /api/products/:id` - обновить информацию о товаре
- `DELETE /api/products/:id` - удалить товар

### Корзина
- `GET /api/cart` - получить содержимое корзины
- `POST /api/cart/items` - добавить товар в корзину
- `DELETE /api/cart/items/:id` - удалить товар из корзины
- `DELETE /api/cart` - очистить корзину

### Заказы
- `GET /api/orders` - получить список заказов пользователя
- `GET /api/orders/:id` - получить информацию о заказе
- `POST /api/orders` - создать новый заказ
- `PATCH /api/orders/:id/status` - обновить статус заказа

## Swagger документация

Swagger UI доступен по адресу: http://localhost:8081/swagger/index.html

## Очистка базы данных

Для очистки базы данных выполните:
```bash
clean_db.bat
```

## Структура проекта

```
shopping-cart/
├── cmd/
│   └── main.go
├── internal/
│   ├── delivery/
│   │   └── http/
│   │       └── handler.go
│   ├── domain/
│   │   └── models.go
│   ├── repository/
│   │   ├── postgres/
│   │   │   └── repository.go
│   │   └── repository.go
│   └── service/
│       ├── impl/
│       │   └── service.go
│       └── service.go
├── docs/
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── docker-compose.yml
├── Dockerfile
├── go.mod
├── go.sum
└── README.md
```

## Технологии

- Go 1.21
- Gin (веб-фреймворк)
- PostgreSQL (база данных)
- GORM (ORM)
- Swagger (документация API)
- Docker и Docker Compose

## TODO

- [ ] Добавить аутентификацию и авторизацию
- [ ] Улучшить обработку ошибок
- [ ] Добавить валидацию входных данных
- [ ] Добавить логирование
- [ ] Написать интеграционные тесты 
