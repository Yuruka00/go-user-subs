# go-user-subs
Сервис для управления подписками пользователей.

## Стек
- **Go 1.25**
- **Chi**
- **PostgreSQL** + **GORM**
- **Slog**
- **Swagger**

## Запуск

1. **Настройка окружения:**
    Создайте файл `.env` в корне проекта. Пример содержимого:
    ```env
    POSTGRES_USER=postgres
    POSTGRES_PASSWORD=123456
    POSTGRES_DB=user_subscriptions
    LOG_LEVEL=debug
    ```

2. **Запуск контейнеров:**
    ```bash
    docker-compose up --build -d
    ```

### Документация доступна по адресу: http://localhost:8080/swagger/index.html. Файлы спецификации хранятся в папке docs/