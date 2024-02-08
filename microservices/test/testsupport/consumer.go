package testsupport

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/Tracking-Detector/td_backend_infra/microservices/shared/queue"
)

type TestQueueConsumer struct {
	queueAdapter  queue.IQueueChannelAdapter
	QueueMessages map[string][]string
	WaitGroup     sync.WaitGroup
	cancelContext context.Context
	cancelFunc    context.CancelFunc
}

func NewTestQueueConsumer(queueAdapter queue.IQueueChannelAdapter) *TestQueueConsumer {

	return &TestQueueConsumer{
		queueAdapter:  queueAdapter,
		QueueMessages: make(map[string][]string),
		WaitGroup:     sync.WaitGroup{},
	}
}

// Consume captures the consumed RabbitMQ messages and stores them in the QueueMessages map.
func (tqc *TestQueueConsumer) Consume(queueName string, expectedMessageCount int) error {
	ctx, cancel := context.WithCancel(context.Background())
	tqc.cancelContext = ctx
	tqc.cancelFunc = cancel
	tqc.WaitGroup.Add(expectedMessageCount)

	msgs, err := tqc.queueAdapter.Consume(
		queueName,
		"TestConsumer_"+queueName, // Consumer tag
		true,                      // Auto-ack
		false,                     // Exclusive
		false,                     // No-local
		false,                     // No-wait
		nil,                       // Args
	)
	if err != nil {
		return fmt.Errorf("failed to register a consumer: %v", err)
	}

	log.Printf("ConsumerService started for queue %s. Waiting for %d messages...", queueName, expectedMessageCount)

	for {
		select {
		case msg, ok := <-msgs:
			if !ok {
				log.Printf("Queue closed. Shutting down ConsumerService for queue %s.", queueName)
				return nil
			}

			messageBody := string(msg.Body)
			log.Printf("Received message: %s", messageBody)

			if _, ok := tqc.QueueMessages[queueName]; !ok {
				tqc.QueueMessages[queueName] = []string{}
			}
			tqc.QueueMessages[queueName] = append(tqc.QueueMessages[queueName], messageBody)
			tqc.WaitGroup.Done()
		case <-tqc.cancelContext.Done():
			log.Printf("Shutting down ConsumerService for queue %s.", queueName)
			return nil
		}
	}
}

// Stop stops the TestQueueConsumer.
func (tqc *TestQueueConsumer) Stop() {
	tqc.cancelFunc()
	// Add additional cleanup logic if needed
}

// WaitForMessages waits for the specified number of messages to be received.
func (tqc *TestQueueConsumer) WaitForMessages(queueName string, expectedMessageCount int) {

	if _, ok := tqc.QueueMessages[queueName]; !ok {
		tqc.QueueMessages[queueName] = []string{}
	}

	currentMessageCount := len(tqc.QueueMessages[queueName])
	if currentMessageCount < expectedMessageCount {
		log.Printf("Waiting for %d more messages in queue %s...", expectedMessageCount-currentMessageCount, queueName)
		tqc.WaitGroup.Wait()
	}

	log.Printf("ConsumerService finished for queue %s. Received %d messages.", queueName, len(tqc.QueueMessages[queueName]))
}

// ClearMessages resets the QueueMessages map.
func (tqc *TestQueueConsumer) ClearMessages() {

	tqc.QueueMessages = make(map[string][]string)
}
