<h1>
Практическое задание №5<br><br>
Ремешевский В.А.<br>
ПИМО-01-25
</h1>

<h2><b>Тема</b><br>
Реализация HTTPS (TLS-сертификаты). Защита от SQL-инъекций
</h2>

# pz5-security

## Краткое описание
Проект демонстрирует запуск Go-приложения по HTTPS с использованием самоподписанного TLS-сертификата, а также защиту от SQL-инъекций при работе с базой данных. В приложении реализованы безопасные SQL-запросы с параметризацией и настройка HTTPS-сервера.

## Структура проекта
```
pz5-security/
├── assets/                      
├── certs/
│   ├── server.crt
│   └── server.key
├── cmd/
│   └── server/
│       └── main.go              
├── internal/
│   ├── config/
│   │   └── config.go
│   ├── httpapi/
│   │   └── handler.go
│   └── student/
│       ├── model.go
│       └── repo.go
├── sql/
│   └── init.sql
├── .gitignore
├── go.mod
└── README.md
```

## Как начать работу

### Инициализация и установка зависимостей

```sh
cd pz5-security/
go mod tidy
go get github.com/lib/pq
go get github.com/joho/godotenv
```

### Настройка базы данных

1. Убедитесь, что у вас запущена база данных PostgreSQL.
2. Создайте базу данных для проекта (например, `pz5_security`).
3. Выполните SQL-запрос из файла `sql/init.sql` в вашей базе данных:

```sql
CREATE TABLE IF NOT EXISTS students (
    id BIGSERIAL PRIMARY KEY,
    full_name TEXT NOT NULL,
    study_group TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE
);

INSERT INTO students (full_name, study_group, email)
VALUES
    ('Иванов Иван Иванович', 'ИВБО-01-25', 'ivanov@example.com'),
    ('Петрова Мария Сергеевна', 'ИВБО-02-25', 'petrova@example.com'),
    ('Сидоров Алексей Андреевич', 'ИВБО-03-25', 'sidorov@example.com')
ON CONFLICT (email) DO NOTHING;
```

### Создание `.env`

1. В корне проекта создайте файл `.env`.
2. Добавьте в него следующие переменные окружения:

```env
DATABASE_URL=postgres://USERNAME:PASSWORD@HOST:PORT/DBNAME?sslmode=disable
ADDR=localhost:8443
CERT_FILE=certs/server.crt
KEY_FILE=certs/server.key
```

- `DATABASE_URL` — параметры подключения к PostgreSQL:
  - `USERNAME` — имя пользователя базы данных.
  - `PASSWORD` — пароль пользователя.
  - `HOST` — адрес сервера базы данных (`localhost` для локальной базы или IP/домен удалённого сервера).
  - `PORT` — порт PostgreSQL (обычно `5432`).
  - `DBNAME` — имя базы данных (например, `pz5_security`).
  - `sslmode` — режим SSL (`disable` для локальной разработки, `require` для продакшена при необходимости).
- `ADDR` — адрес и порт, на котором будет запущен HTTPS-сервер (например, `localhost:8443`).
- `CERT_FILE` — путь к файлу TLS-сертификата (например, `certs/server.crt`).
- `KEY_FILE` — путь к файлу приватного ключа (например, `certs/server.key`).

### Запуск приложения

```sh
go run ./cmd/server
```

## TLS-сертификат: что это и как сгенерировать

TLS-сертификат обеспечивает защищённое соединение между клиентом и сервером по протоколу HTTPS. Для разработки можно использовать самоподписанный сертификат. Для генерации используйте команду:

```powershell
openssl req -x509 -newkey rsa:2048 -nodes `
-keyout certs/server.key `
-out certs/server.crt `
-days 365 `
-subj "/CN=localhost" `
-config NUL
```

В результате будут созданы два файла:
- `certs/server.key` — приватный ключ
- `certs/server.crt` — публичный сертификат

## SQL-инъекции: что это и как защититься

SQL-инъекция — это атака, при которой злоумышленник может внедрить вредоносный SQL-код через пользовательский ввод. Никогда не подставляйте значения напрямую в SQL-запросы!

**Небезопасный пример:**
```go
rows, err := db.Query(fmt.Sprintf("SELECT * FROM students WHERE id = %s", rawID))
```
Так делать нельзя!

**Безопасный пример:**
```go
row := db.QueryRow(
    "SELECT id, full_name, study_group, email FROM students WHERE id = $1",
    id,
)
```
Используйте параметризацию запросов — это защищает от SQL-инъекций.

## Скриншоты

### Генерация самоподписанного сертификата
```powershell
openssl req -x509 -newkey rsa:2048 -nodes `
-keyout certs/server.key `
-out certs/server.crt `
-days 365 `
-subj "/CN=localhost" `
-config NUL
```
![gen_cert](assets/certificate_generation.png)

### Проверка эндпоинта /health по HTTPS
```sh
curl -k https://localhost:8443/health
```
![health](assets/health_endpoint.png)

### Получение информации о студенте по id
```sh
curl -k "https://localhost:8443/students?id=1"
```
![student_by_id](assets/get_student.png)

### Ошибка: студент не найден
```sh
curl -k "https://localhost:8443/students?id=999"
```
![student_not_found](assets/student_not_found.png)
