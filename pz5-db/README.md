<h1>
Практическое задание №5<br><br>
Ремешевский В.А.<br>
ПИМО-01-25<br></h1>

<h3>Тема: подключение к PostgreSQL через database/sql, выполнение простых запросов (INSERT, SELECT)</h3>

---

# PZ5-DB

## Описание проекта

**PZ5-DB** — простейший таск-менеджер, демонстрирующий работу с базой данных PostgreSQL через Go, включая базовые операции CRUD, пул соединений и работу с транзакциями.

### Окружение

- Go 1.25.1
- PostgreSQL 16.10 (Ubuntu 16.10-0ubuntu0.24.04.1)

---

## Структура проекта

```
pz5-db/
├── assets/
├── db.go         
├── main.go   
├── repository.go
├── go.mod
├── .gitignore
└── README.md
```
---

## Как начать работу

### Инициализация и установка зависимостей

```sh
cd pz5-db
go mod init example.com/pz5-db
go get github.com/jackc/pgx/v5/stdlib
go get github.com/joho/godotenv
```

### Запуск приложения

```powershell
go run .
```

## Сборка приложения (.exe для Windows)

```sh
go build -o bin/server.exe ./...
```

### Запуск приложения

```sh
bin\server.exe
```

## Скриншоты

### Создание БД/таблицы в psql
![psql_db_creation](assets/psql_db_creation.png)
---

### Успешный вывод ```go run .```
```sh
go run .
```
![go_run_output](assets/go_run_output.png)
---

### ```SELECT * FROM tasks``` в psql
```sh
SELECT * FROM tasks
```
![psql_selection](assets/psql_selection.png)
---

### ```SELECT * FROM tasks``` в pgAdmin 4
```sh
SELECT * FROM tasks
```
![pgadmin_selection](assets/pgadmin_selection.png)
---

# Код-листинг

```main.go```:
```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL not set — please create a .env file or export the variable")
	}

	db, err := openDB(dsn)
	if err != nil {
		log.Fatalf("openDB error: %v", err)
	}
	defer db.Close()

	repo := NewRepo(db)

	// 1) Вставим пару задач
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	titles := []string{"Сделать ПЗ №5", "Купить кофе", "Проверить отчёты"}
	for _, title := range titles {
		id, err := repo.CreateTask(ctx, title)
		if err != nil {
			log.Fatalf("CreateTask error: %v", err)
		}
		log.Printf("Inserted task id=%d (%s)", id, title)
	}

	// Массовая вставка задач
	ctxMany, cancelMany := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelMany()

	mass_titles := []string{"Пресс качат", "Т) Беигт", "Турник", "Анжуманя"}
	if err := repo.CreateMany(ctxMany, mass_titles); err != nil {
		log.Fatalf("CreateMany error: %v", err)
	}

	fmt.Println("Mass insertion completed")

	// 2) Прочитаем список задач
	ctxList, cancelList := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancelList()

	tasks, err := repo.ListTasks(ctxList)
	if err != nil {
		log.Fatalf("ListTasks error: %v", err)
	}

	// 3) Напечатаем
	fmt.Println("=== Tasks ===")
	for _, t := range tasks {
		fmt.Printf("#%d | %-24s | done=%-5v | %s\n",
			t.ID, t.Title, t.Done, t.CreatedAt.Format(time.RFC3339))
	}

	// Вывод списка выполненных задач (done=true)
	ctxDone, cancelDone := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancelDone()

	doneTasks, err := repo.ListDone(ctxDone, true)
	if err != nil {
		log.Fatalf("ListDone error: %v", err)
	}

	fmt.Println("=== Done tasks ===")
	for _, t := range doneTasks {
		fmt.Printf("#%d | %-24s | done=%-5v | %s\n",
			t.ID, t.Title, t.Done, t.CreatedAt.Format(time.RFC3339))
	}

	// Поиск задачи по ее ID
	ctxFind, cancelFind := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancelFind()

	taskID := 2
	task, err := repo.FindByID(ctxFind, taskID)
	if err != nil {
		log.Fatalf("FindByID error: %v", err)
	}

	if task == nil {
		fmt.Printf("Task with ID=%d not found\n", taskID)
	} else {
		fmt.Printf("=== Task #%d details ===\n", task.ID)
		fmt.Printf("Title: %s\n", task.Title)
		fmt.Printf("Done: %v\n", task.Done)
		fmt.Printf("Created at: %s\n", task.CreatedAt.Format(time.RFC3339))
	}
}
```

## Ответы на вопросы

**Что такое пул соединений ```sql.DB``` и зачем его настраивать?**  

```*sql.DB``` — это абстракция пула соединений к СУБД. Она управляет набором реальных TCP/сокет соединений и повторно использует их для запросов. Настройка (```MaxOpenConns```, ```MaxIdleConns```, ```ConnMaxLifetime```) нужна, чтобы:
- ограничить потребление ресурсов БД (не открывать слишком много соединений),
- уменьшить накладные расходы на установку/закрытие соединений,
- обеспечить стабильную производительность под нагрузкой.

---

**Почему используем плейсхолдеры ```$1```, ```$2```?**

Плейсхолдеры (```$1```, ```$2``` ...) — это параметризация SQL-запросов в PostgreSQL/pgx. Они защищают от SQL‑инъекций и позволяют драйверу корректно сериализовать значения и кэшировать план выполнения.

---

**Чем отличаются ```Query```, ```QueryRow``` и ```Exec```?**
- ```Query``` — выполняет запрос, возвращает множество строк (rows), нужно итерироваться и закрывать rows.  
- ```QueryRow``` — для запроса, возвращающего одну строку; упрощённый вариант (возвращает Row и Scan).  
- ```Exec``` — выполняет команду без возврата строк (INSERT/UPDATE/DELETE без RETURNING), возвращает результат (RowsAffected) и ошибку.

---

## Обоснование транзакций и настроек пула

### **Транзакции (почему и где применены)**  
В проекте транзакция используется для массовой вставки (```CreateMany```). Это оправдано: при батчевой вставке каждая вставка должна выполняться в рамках одной логической операции — при ошибке нужно откатить уже вставленные строки, чтобы не получить частичный результат. Для одиночных INSERT/UPDATE/DELETE транзакция часто не требуется, так как каждый SQL‑запрос атомарен сам по себе.


### **Обоснование настроек пула (параметры для удалённого сервера)**  
Хост: удалённый сервер Ubuntu 24.04 с характеристиками: 1 CPU, 2 GB RAM, 40 GB диск. Исходя из этого:
- Ограничить число одновременно открытых соединений, чтобы не исчерпать ресурсы процессора и памяти. Рекомендуемая конфигурация для такого сервера: ```MaxOpenConns``` = 3–6, ```MaxIdleConns``` = 1–2. В проекте установлены ```db.SetMaxOpenConns(5)``` и ```db.SetMaxIdleConns(2)``` — это компромисс между производительностью и безопасностью использования ресурсов.  
- ```ConnMaxLifetime``` = 30m — позволяет периодически обновлять соединения и избегать долгоживущих проблем с сетью/аутентификацией.  
Таким образом выбранные параметры позволяют обслуживать небольшой объём параллельных запросов без перегрузки БД и ОС.

---

## Использование файла `.env`

В проекте файл `.env` используется для хранения конфиденциальных данных и параметров подключения к базе данных.  

**Зачем нужен `.env`:**

- Сохраняет настройки подключения (DSN) отдельно от кода, чтобы не хранить пароли в репозитории.
- Позволяет легко менять параметры для локальной разработки, тестового или удалённого сервера.
- Обеспечивает безопасность и удобство при работе в команде.

---

### Создание `.env`

1. В корне проекта создайте файл `.env`.
2. Добавьте в него переменную `DATABASE_URL` с параметрами подключения к PostgreSQL:

```env
DATABASE_URL=postgres://USERNAME:PASSWORD@HOST:PORT/DBNAME?sslmode=disable
```

- ```USERNAME``` — имя пользователя базы данных.
- ```PASSWORD``` — пароль пользователя.
- ```HOST``` — адрес сервера базы данных (```localhost``` для локальной базы или IP/домен удалённого сервера).
- ```PORT``` — порт PostgreSQL (обычно ```5432```).
- ```DBNAME``` — имя базы данных.
- ```sslmode``` — режим SSL (```disable``` для локальной разработки, ```require``` для продакшена при необходимости).