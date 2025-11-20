<h1>
Практическое задание №8<br><br>
Ремешевский В.А.<br>
ПИМО-01-25
</h1>

<h2><b>Тема</b><br>
Работа с MongoDB: подключение, создание коллекции, CRUD-операции</h2><br>

# PZ8-MONGO

## Краткое описание проекта

**PZ8-MONGO** — это HTTP API‑сервис на Go, реализующий CRUD-операции для заметок с хранением в MongoDB. Проект демонстрирует работу с коллекцией заметок, поиск, обновление, удаление, а также поддержку TTL (автоматическое удаление по времени) и статистику по коллекции.

## Требования

- Go ≥ 1.21
- Docker **или** установленный MongoDB Server (и MongoDB Shell)
- curl/Postman для тестирования API

## Структура проекта

```
pz8-mongo/
├── assets/
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── db/
│   │   └── mongo.go
│   └── notes/
│       ├── handler.go
│       ├── model.go
│       └── repo.go
├── .env.example
├── go.mod
└── README.md
```

---

## Настройка переменных окружения

Создайте файл `.env` на основе `.env.example` и укажите параметры подключения к MongoDB:

```env
MONGO_URI=mongodb://localhost:27017
MONGO_DB=pz8
MONGO_USER=root
MONGO_PASS=secret
HTTP_ADDR=:8080

# Если требуется подключение с аутентификацией, файл .env должен выглядеть так:
# MONGO_URI=mongodb://root:secret@localhost:27017/?authSource=admin
# MONGO_DB=pz8
# MONGO_USER=root
# MONGO_PASS=secret
# HTTP_ADDR=:8080
```

- `MONGO_URI` — строка подключения к MongoDB (обычно `mongodb://localhost:27017` для локального запуска или через Docker).
- `MONGO_DB` — имя базы данных для хранения заметок.
- `MONGO_USER` — имя пользователя MongoDB (если требуется).
- `MONGO_PASS` — пароль пользователя MongoDB (если требуется).
- `HTTP_ADDR` — порт, на котором заупщено приложение.

---

## Команды запуска

### Запуск MongoDB в Docker

```sh
docker compose up -d
```

### Копирование файла переменных окружения

```sh
cp .env.example .env
```
**PowerShell:**
```powershell
Copy-Item .env.example .env
```

### Запуск приложения

```sh
go run ./cmd/api
```

---

## Скриншоты

### Создание заметки
```sh
curl -Method POST http://localhost:8080/api/v1/notes `
  -ContentType "application/json" `
  -Body '{"title":"First note","content":"Hello Mongo!"}'
```
![create_note](assets/create_note.png)

### Создание заметки с TTL
```sh
curl -Method POST http://localhost:8080/api/v1/notes `
  -ContentType "application/json" `
  -Body '{"title":"Expiring note","content":"I am dying!","ttl_seconds":60}'
```
![create_note_ttl](assets/create_note_ttl.png)

### Получение списка заметок с фильтрацией и пагинацией
```sh
curl "http://localhost:8080/api/v1/notes?limit=5&skip=0&q=first"
```
![list_notes](assets/list_notes.png)

### Поиск заметок по тексту
```sh
curl 'http://localhost:8080/api/v1/notes?q=hello'
```
![search_notes](assets/search_notes.png)

### Получение заметки по id
```sh
curl http://localhost:8080/api/v1/notes/<id>
```
![get_note](assets/get_note.png)

### Обновление заметки
```sh
curl -Method PATCH http://localhost:8080/api/v1/notes/<id> `
  -ContentType "application/json" `
  -Body '{"content":"Hello Mongo! got updated"}'
```
![update_note](assets/update_note.png)

### Удаление заметки
```sh
curl -Method DELETE http://localhost:8080/api/v1/notes/<id>
```
![delete_note](assets/delete_note.png)


### Получение статистики по коллекции
```sh
curl 'http://localhost:8080/api/v1/notes/stats'
```
![notes_stats](assets/notes_stats.png)

---

## Отладка и типовые ошибки

- **connection refused / i/o timeout** — проверьте, что контейнер с MongoDB запущен и порт 27017 проброшен.
- **Аутентификация** — используйте параметр подключения `?authSource=admin` и правильные значения `MONGO_USER`/`MONGO_PASS`.
- **duplicate key error** — сработал уникальный индекс по `title`. В случае такой ошибки API возвращает HTTP 409 при создании заметки.
- **invalid ObjectID** — при неверном id возвращается HTTP 404 (ошибка мапится на ErrNotFound).
- **CORS** (если тестируете из браузера) — добавьте middleware CORS или используйте curl/Postman.

---