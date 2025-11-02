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
