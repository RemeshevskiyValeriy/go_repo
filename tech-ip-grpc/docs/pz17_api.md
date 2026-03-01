# API PZ17: Auth и Tasks

## Переменные окружения и порты

- **AUTH_PORT** — порт сервиса авторизации (по умолчанию 8081)
- **TASKS_PORT** — порт сервиса задач (по умолчанию 8082)
- **AUTH_BASE_URL** — базовый URL сервиса авторизации (например, http://localhost:8081)

---

## Auth endpoints

### POST /v1/auth/login
Авторизация пользователя. Возвращает токен.

**Request:**
```json
{
  "username": "student",
  "password": "student"
}
```

**Response:**
- 200 OK: `{ "access_token": "demo-token", "token_type": "Bearer" }`
- 400 Bad Request: неверный JSON
- 401 Unauthorized: неверные данные

**Пример curl (PowerShell):**
```sh
curl -Method POST http://localhost:8081/v1/auth/login `
  -Headers @{
    "Content-Type" = "application/json"
    "X-Request-ID" = "req-001"
  } `
  -Body '{"username":"student","password":"student"}'
```

---

### GET /v1/auth/verify
Проверка токена авторизации.

**Headers:**
- Authorization: Bearer demo-token

**Response:**
- 200 OK: `{ "valid": true, "subject": "student" }`
- 401 Unauthorized: `{ "valid": false, "error": "unauthorized" }`

**Пример curl (PowerShell):**
```sh
curl http://localhost:8081/v1/auth/verify `
  -Headers @{
    "Authorization" = "Bearer demo-token"
    "X-Request-ID"  = "req-002"
  }
```

---

## Tasks endpoints

### GET /v1/tasks
Список задач (требует авторизации).

**Response:**
- 200 OK: `[ ... ]` — массив задач
- 401 Unauthorized: нет токена или неверный токен

**Пример curl (без авторизации):**
```sh
curl http://localhost:8082/v1/tasks `
  -Headers @{
    "X-Request-ID" = "req-no-auth"
  }
```

**Пример curl (с авторизацией):**
```sh
curl http://localhost:8082/v1/tasks `
  -Headers @{
    "Authorization" = "Bearer demo-token"
    "X-Request-ID"  = "req-list"
  }
```

---

### POST /v1/tasks
Создать задачу.

**Request:**
```json
{
  "title": "Do PZ17",
  "description": "make screenshots",
  "due_date": "2026-03-10"
}
```

**Response:**
- 201 Created: объект задачи
- 400 Bad Request: нет title или неверный JSON
- 401 Unauthorized: нет токена или неверный токен

**Пример curl:**
```sh
curl -Method POST http://localhost:8082/v1/tasks `
  -Headers @{
    "Content-Type"  = "application/json"
    "Authorization" = "Bearer demo-token"
    "X-Request-ID"  = "req-create-task"
  } `
  -Body '{
    "title": "Do PZ17",
    "description": "make screenshots",
    "due_date": "2026-03-10"
  }'
```

---

### GET /v1/tasks/{id}
Получить задачу по id.

**Response:**
- 200 OK: объект задачи
- 401 Unauthorized: нет токена или неверный токен
- 404 Not Found: задача не найдена

**Пример curl:**
```sh
curl http://localhost:8082/v1/tasks/t_001 `
  -Headers @{
    "Authorization" = "Bearer demo-token"
    "X-Request-ID"  = "req-get"
  }
```

---

### PATCH /v1/tasks/{id}
Обновить задачу.

**Request:**
```json
{
  "title": "Do PZ17 (updated)",
  "done": true
}
```

**Response:**
- 200 OK: обновлённая задача
- 400 Bad Request: неверный JSON
- 401 Unauthorized: нет токена или неверный токен
- 404 Not Found: задача не найдена

**Пример curl:**
```sh
curl -Method PATCH http://localhost:8082/v1/tasks/t_001 `
  -Headers @{
    "Content-Type"  = "application/json"
    "Authorization" = "Bearer demo-token"
    "X-Request-ID"  = "req-update"
  } `
  -Body '{
    "title": "Do PZ17 (updated)",
    "done": true
  }'
```

---

### DELETE /v1/tasks/{id}
Удалить задачу.

**Response:**
- 204 No Content: успешно удалено
- 401 Unauthorized: нет токена или неверный токен
- 404 Not Found: задача не найдена

**Пример curl:**
```sh
curl -Method DELETE http://localhost:8082/v1/tasks/t_001 `
  -Headers @{
    "Authorization" = "Bearer demo-token"
    "X-Request-ID"  = "req-delete"
  }
```

---

### GET /v1/tasks/{id} (после удаления)
**Response:**
- 404 Not Found: задача не найдена

**Пример curl:**
```sh
curl http://localhost:8082/v1/tasks/t_001 `
  -Headers @{
    "Authorization" = "Bearer demo-token"
    "X-Request-ID"  = "req-after-delete"
  }
```

---

## Коды ответов
- 200 OK — успешный запрос
- 201 Created — создано
- 204 No Content — удалено
- 400 Bad Request — ошибка запроса
- 401 Unauthorized — нет авторизации
- 404 Not Found — не найдено
- 502 Bad Gateway — ошибка сервиса авторизации

## Примечания
- Для всех запросов рекомендуется указывать заголовок `X-Request-ID` для трекинга.
- Для задач требуется авторизация через Bearer-токен.
