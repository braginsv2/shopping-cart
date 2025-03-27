-- Очистка таблиц в правильном порядке (с учетом внешних ключей)
DELETE FROM order_items;
DELETE FROM orders;
DELETE FROM cart_items;
DELETE FROM carts;
DELETE FROM products;

-- Сброс автоинкрементных счетчиков
ALTER SEQUENCE order_items_id_seq RESTART WITH 1;
ALTER SEQUENCE orders_id_seq RESTART WITH 1;
ALTER SEQUENCE cart_items_id_seq RESTART WITH 1;
ALTER SEQUENCE carts_id_seq RESTART WITH 1;
ALTER SEQUENCE products_id_seq RESTART WITH 1; 