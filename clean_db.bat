@echo off
echo Очистка базы данных...
psql -h localhost -p 5432 -U postgres -d shopping_cart -f clean_db.sql
echo База данных очищена.
pause 