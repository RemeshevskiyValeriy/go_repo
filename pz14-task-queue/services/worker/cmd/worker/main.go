package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	mainQueue   = "task_jobs"
	dlqQueue    = "task_jobs_dlq"
	maxAttempts = 3
)

type TaskJob struct {
	Job       string `json:"job"`
	TaskID    string `json:"task_id"`
	Attempt   int    `json:"attempt"`
	MessageID string `json:"message_id"`
}

type ProcessedStore struct {
	items map[string]bool
}

func NewProcessedStore() *ProcessedStore {
	return &ProcessedStore{
		items: make(map[string]bool),
	}
}

func (s *ProcessedStore) Exists(id string) bool {
	return s.items[id]
}

func (s *ProcessedStore) MarkDone(id string) {
	s.items[id] = true
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

func processTask(job TaskJob) error {
	log.Printf("processing task_id=%s message_id=%s attempt=%d", job.TaskID, job.MessageID, job.Attempt)

	time.Sleep(2 * time.Second)

	if job.TaskID == "t_fail" {
		return fmt.Errorf("simulated processing error")
	}

	return nil
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

	if err := ch.Qos(1, 0, false); err != nil {
		log.Fatalf("qos error: %v", err)
	}

	msgs, err := ch.Consume(
		mainQueue,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("consume error: %v", err)
	}

	processed := NewProcessedStore()

	log.Println("worker started, waiting for jobs...")

	for d := range msgs {
		var job TaskJob

		if err := json.Unmarshal(d.Body, &job); err != nil {
			log.Printf("bad message: %v", err)
			_ = d.Nack(false, false)
			continue
		}

		if processed.Exists(job.MessageID) {
			log.Printf("duplicate message skipped message_id=%s", job.MessageID)
			_ = d.Ack(false)
			continue
		}

		if err := processTask(job); err != nil {
			log.Printf("processing error task_id=%s attempt=%d error=%v", job.TaskID, job.Attempt, err)

			job.Attempt++

			if job.Attempt <= maxAttempts {
				log.Printf("retry publish task_id=%s next_attempt=%d", job.TaskID, job.Attempt)

				if err := publishJob(ch, mainQueue, job); err != nil {
					log.Printf("retry publish error: %v", err)
				}

				_ = d.Ack(false)
				continue
			}

			log.Printf("max attempts exceeded, sending to DLQ task_id=%s message_id=%s", job.TaskID, job.MessageID)

			if err := publishJob(ch, dlqQueue, job); err != nil {
				log.Printf("dlq publish error: %v", err)
			}

			_ = d.Ack(false)
			continue
		}

		processed.MarkDone(job.MessageID)
		log.Printf("job done and acked task_id=%s message_id=%s", job.TaskID, job.MessageID)

		_ = d.Ack(false)
	}
}