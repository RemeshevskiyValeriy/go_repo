<h1>
Практическое задание №14<br><br>
Ремешевский В.А.<br>
ПИМО-01-25
</h1>

<h2><b>Тема</b><br>
Оптимизация запросов к БД. Использование connection pool</h2>


# Цели работы

1. Научиться находить «узкие места» в SQL-запросах и устранять их (индексы, переписывание запросов, пагинация, батчинг).
2. Освоить настройку пула подключений (**connection pool**) в Go и параметры его тюнинга.
3. Научиться использовать `EXPLAIN/ANALYZE`, базовые метрики (`pg_stat_statements`), подготовленные запросы и транзакции.
4. Применить техники уменьшения **N+1** запросов и сокращения аллокаций на горячем пути.

---

## Используемый стек

- Go 1.25.1
- PostgreSQL (драйвер `pgx`)
- Конфиг пула: `MaxOpenConns=5`, `MaxIdleConns=2`, `ConnMaxLifetime=30m`

---

# NOTES-API-POSTGRES

## Краткое описание проекта

**NOTES-API-POSTGRES** — это учебный HTTP API‑сервис на Go с хранением заметок в PostgreSQL, реализующий оптимизированные методы выборки, пагинации и батчинга, а также профилирование и нагрузочное тестирование.

## Структура проекта

```
notes-api-postgres/
├── api/
│   └── openapi.yaml
├── assets/
├── cmd/
│   └── api/
│       └── main.go
├── docs/
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── internal/
│   ├── core/
│   │   ├── dto.go
│   │   └── note.go
│   ├── http/
│   │   ├── router.go
│   │   └── handlers/
│   │       ├── error.go
│   │       └── notes.go
│   └── repo/
│       ├── note_mem.go
│       ├── note_repo_pg.go
│       └── note_repo.go
├── go.mod
├── Makefile
├── README.md
```

---

## Как начать работу

### Инициализация и установка зависимостей

```sh
cd notes-api-postgres
go mod tidy
go get github.com/jackc/pgx/v5
go get github.com/joho/godotenv
go get github.com/go-chi/chi/v5
go get github.com/jackc/pgx/v5/stdlib
```

### Создание `.env`

1. В корне проекта создайте файл `.env`.
2. Добавьте в него переменную `DATABASE_URL` с параметрами подключения к PostgreSQL:

```env
DATABASE_URL=postgres://USERNAME:PASSWORD@HOST:PORT/DBNAME?sslmode=disable
```

- `USERNAME` — имя пользователя базы данных.
- `PASSWORD` — пароль пользователя.
- `HOST` — адрес сервера базы данных (`localhost` для локальной базы или IP/домен удалённого сервера).
- `PORT` — порт PostgreSQL (обычно `5432`).
- `DBNAME` — имя базы данных.
- `sslmode` — режим SSL (`disable` для локальной разработки, `require` для продакшена при необходимости).

### Запуск приложения

```powershell
go run ./cmd/api
```

Или через `Makefile`:
```sh
make run
```

---

### **Обоснование настроек пула (параметры для удалённого сервера)**

Хост: удалённый сервер Ubuntu 24.04 с характеристиками: 1 CPU, 2 GB RAM, 40 GB диск.  
Исходя из этого:
- Ограничить число одновременно открытых соединений, чтобы не исчерпать ресурсы процессора и памяти. Рекомендуемая конфигурация для такого сервера: `MaxOpenConns = 3–6`, `MaxIdleConns = 1–2`. В проекте установлены `db.SetMaxOpenConns(5)` и `db.SetMaxIdleConns(2)` — это компромисс между производительностью и безопасностью использования ресурсов.
- `ConnMaxLifetime = 30m` — позволяет периодически обновлять соединения и избегать долгоживущих проблем с сетью/аутентификацией.

Таким образом выбранные параметры позволяют обслуживать небольшой объём параллельных запросов без перегрузки БД и ОС.

---

## Скриншоты

### Проблемные запросы до оптимизации (EXPLAIN/ANALYZE в pgAdmin 4)

**OFFSET pagination (большой offset):**
```sql
EXPLAIN (ANALYZE, BUFFERS)
SELECT id, title, content, created_at
FROM notes
ORDER BY created_at DESC, id DESC
OFFSET $1 LIMIT $2;
```
![explain_offset](assets/offset.png)

**Single SELECT по id:**
```sql
EXPLAIN (ANALYZE, BUFFERS)
SELECT id, title FROM notes WHERE id = $1;
```
![explain_single](assets/n_plus_one.png)

---

### Нагрузочное тестирование (hey)

**OFFSET pagination:**
```sh
hey -z 30s -c 50 "http://localhost:8080/api/v1/notes?offset=40000&limit=20"
```
![hey_offset](assets/offset_pagination.png)

**Single SELECT по id:**
```sh
hey -z 30s -c 50 "http://localhost:8080/api/v1/notes/40000"
```
![hey_single](assets/single_select.png)

**KEYSET pagination:**
```sh
hey -z 30s -c 50 "http://localhost:8080/api/v1/notes"
```
![hey_keyset](assets/keyset_pagination.png)

**BATCH SELECT:**
```sh
hey -z 30s -c 50 "http://localhost:8080/api/v1/notes/batch?ids=10001,10002,10003,10004,10005,10006,10007,10008,10009,10010"
```
![hey_batch](assets/batch_any.png)

---

## Какие индексы/переписывания применили (SQL-фрагменты)

**Индексы:**
```sql
-- Частичный индекс для поиска по заголовку, если часто ищем по prefix
CREATE INDEX IF NOT EXISTS idx_notes_title_gin
  ON notes USING GIN (to_tsvector('simple', title));

-- Индекс для keyset-пагинации
CREATE INDEX IF NOT EXISTS idx_notes_created_id ON notes (created_at, id);
```

**Оптимизированные запросы:**
```sql
-- последующие страницы (используйте последнее seen_created_at, seen_id)
SELECT id, title, content, created_at
FROM notes
WHERE (created_at, id) < ($1, $2)
ORDER BY created_at DESC, id DESC
LIMIT $3;

-- один батч:
SELECT id, title FROM notes WHERE id = ANY($1); -- массив bigint[]
```

---

## Сводная таблица производительности (до / после оптимизации)

| Эндпоинт                        | Реализация         | RPS     | p95      | p99      | Ошибки         |
|----------------------------------|--------------------|---------|----------|----------|----------------|
| GET /api/v1/notes?offset=40000   | OFFSET pagination  | ~88     | ~1.52 s  | ~2.19 s  | ~30% (500)     |
| GET /api/v1/notes                | Keyset pagination  | ~90     | ~1.51 s  | ~2.25 s  | ~30% (500)     |
| GET /api/v1/notes/{id}           | Single SELECT      | ~90     | ~1.49 s  | ~2.39 s  | ~30% (500)     |
| GET /api/v1/notes/batch          | Batch (ANY)        | ~27 955 | ~8 ms    | ~26 ms   | 400*           |

* 400 ошибки связаны с валидацией входных параметров, а не с производительностью БД.

---

## Выводы

- **OFFSET пагинация** плохо масштабируется при больших offset.
- **Keyset pagination** требует составного индекса, иначе выигрыш минимален.
- **Batch-запросы** (`WHERE id = ANY($1)`) радикально снижают latency и повышают RPS.
- Основной выигрыш в системе дала ликвидация N+1, а не пагинация.
- Наибольший эффект дали оптимизация выборок по id и батчинг, а также правильная настройка connection pool.
- В дальнейшем стоит рассмотреть оптимизацию других "горячих" путей, расширение индексов и анализ `pg_stat_statements` для выявления новых узких мест.