package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	mainQueue = "task_jobs"
	dlqQueue  = "task_jobs_dlq"
)

type CreateJobRequest struct {
	TaskID string `json:"task_id"`
}

type TaskJob struct {
	Job       string `json:"job"`
	TaskID    string `json:"task_id"`
	Attempt   int    `json:"attempt"`
	MessageID string `json:"message_id"`
}

type CreateJobResponse struct {
	Status    string `json:"status"`
	TaskID    string `json:"task_id"`
	MessageID string `json:"message_id"`
}

func generateMessageID() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "fallback-message-id"
	}
	return hex.EncodeToString(b)
}

func declareQueues(ch *amqp.Channel) error {
	_, err := ch.QueueDeclare(
		dlqQueue,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	_, err = ch.QueueDeclare(
		mainQueue,
		true,
		false,
		false,
		false,
		nil,
	)
	return err
}

func publishJob(ch *amqp.Channel, queue string, job TaskJob) error {
	body, err := json.Marshal(job)
	if err != nil {
		return err
	}

	return ch.PublishWithContext(
		context.Background(),
		"",
		queue,
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

	if err := declareQueues(ch); err != nil {
		log.Fatalf("queue declare error: %v", err)
	}

	http.HandleFunc("/v1/jobs/process-task", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req CreateJobRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "bad json", http.StatusBadRequest)
			return
		}

		if req.TaskID == "" {
			http.Error(w, "task_id is required", http.StatusBadRequest)
			return
		}

		job := TaskJob{
			Job:       "process_task",
			TaskID:    req.TaskID,
			Attempt:   1,
			MessageID: generateMessageID(),
		}

		if err := publishJob(ch, mainQueue, job); err != nil {
			log.Printf("publish job error: %v", err)
			http.Error(w, "publish error", http.StatusInternalServerError)
			return
		}

		log.Printf("job accepted task_id=%s message_id=%s attempt=%d", job.TaskID, job.MessageID, job.Attempt)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
		_ = json.NewEncoder(w).Encode(CreateJobResponse{
			Status:    "accepted",
			TaskID:    job.TaskID,
			MessageID: job.MessageID,
		})
	})

	log.Println("tasks service started on :8082")
	log.Fatal(http.ListenAndServe(":8082", nil))
}