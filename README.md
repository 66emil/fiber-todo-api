# fiber-todo-api

REST API для управления задачами на Go: регистрация, JWT-аутентификация и CRUD по личному списку задач. Каждый пользователь видит и меняет только свои записи.

## Стек

Go 1.24, [Fiber v2](https://gofiber.io/), GORM с драйвером PostgreSQL, JWT (`golang-jwt/v5`), `godotenv`.

## Эндпоинты

| Метод | Путь | Доступ | Описание |
|---|---|---|---|
| `POST` | `/auth/register` | открытый | Регистрация по email и паролю |
| `POST` | `/auth/login` | открытый | Вход, возвращает JWT |
| `POST` | `/api/todos` | JWT | Создать задачу |
| `GET` | `/api/todos` | JWT | Список задач текущего пользователя |
| `PATCH` | `/api/todos/:id` | JWT | Обновить задачу |
| `DELETE` | `/api/todos/:id` | JWT | Удалить задачу |

Группа `/api` закрыта middleware `AuthRequired`: он проверяет токен и кладёт `user_id` в контекст запроса, а обработчики фильтруют выборки по этому идентификатору.

## Модели

```go
type User struct {
    gorm.Model
    Email        string  // уникальный индекс
    PasswordHash string
    Todos        []Todo
}

type Todo struct {
    gorm.Model
    UserID  uint
    Title   string
    IsDone  bool       // по умолчанию false
    Duedate *time.Time // опционально
}
```

Схема применяется автомиграцией GORM при старте.

## Запуск

```bash
git clone https://github.com/66emil/fiber-todo-api.git
cd fiber-todo-api
go mod download
```

Создайте `.env` в корне:

```env
DB_HOST=localhost
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=todos
DB_PORT=5432
DB_SSLMODE=disable
```

```bash
go run main.go
```

Сервер поднимется на `:3000`.

## Пример использования

```bash
# регистрация
curl -X POST http://localhost:3000/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"me@example.com","password":"secret"}'

# вход — сохраните токен из ответа
curl -X POST http://localhost:3000/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"me@example.com","password":"secret"}'

# создать задачу
curl -X POST http://localhost:3000/api/todos \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"title":"Написать README"}'

# список задач
curl http://localhost:3000/api/todos -H "Authorization: Bearer $TOKEN"
```

## Структура

```
main.go              точка входа и регистрация маршрутов
pkg/
├── database/        подключение к PostgreSQL и автомиграция
├── handlers/        auth.go — регистрация и вход
│                    todos.go — CRUD по задачам
├── middleware/      auth.go — проверка JWT
├── models/          User и Todo
└── utils/           вспомогательные функции
```
