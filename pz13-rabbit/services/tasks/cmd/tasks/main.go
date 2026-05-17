package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type CreateTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type TaskResponse struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

type TaskEvent struct {
	Event  string `json:"event"`
	TaskID string `json:"task_id"`
	TS     string `json:"ts"`
}

func publishTaskCreated(ch *amqp.Channel, queueName, taskID string) error {
	_, err := ch.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	event := TaskEvent{
		Event:  "task.created",
		TaskID: taskID,
		TS:     time.Now().UTC().Format(time.RFC3339),
	}

	body, err := json.Marshal(event)
	if err != nil {
		return err
	}

	return ch.PublishWithContext(
		context.Background(),
		"",
		queueName,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Body:         body,
		},
	)
}

func main() {
	rabbitURL := os.Getenv("RABBIT_URL")
	if rabbitURL == "" {
		rabbitURL = "amqp://guest:guest@localhost:5672/"
	}

	queueName := os.Getenv("QUEUE_NAME")
	if queueName == "" {
		queueName = "task_events"
	}

	conn, err := amqp.Dial(rabbitURL)
	if err != nil {
		log.Fatalf("rabbit connect error: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("channel error: %v", err)
	}
	defer ch.Close()

	http.HandleFunc("/v1/tasks", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req CreateTaskRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "bad json", http.StatusBadRequest)
			return
		}

		taskID := "t_" + time.Now().Format("150405")

		resp := TaskResponse{
			ID:          taskID,
			Title:       req.Title,
			Description: req.Description,
			Status:      "created",
		}

		if err := publishTaskCreated(ch, queueName, taskID); err != nil {
			log.Printf("publish error, task was still created: %v", err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(resp)
	})

	log.Println("tasks service started on :8082")
	log.Fatal(http.ListenAndServe(":8082", nil))
}