package testsupport

import (
	"log"
	"sync"
	"tds/shared/queue"
)

type TestQueueConsumer struct {
	sync.Mutex
	queueAdapter  queue.IQueueChannelAdapter
	QueueMessages map[string][]string
	WaitGroup     sync.WaitGroup
}

func NewTestQueueConsumer(queueAdapter queue.IQueueChannelAdapter) *TestQueueConsumer {
	return &TestQueueConsumer{
		queueAdapter:  queueAdapter,
		QueueMessages: make(map[string][]string),
		WaitGroup:     sync.WaitGroup{},
	}
}

// Consume captures the consumed RabbitMQ messages and stores them in the QueueMessages map.
func (tqc *TestQueueConsumer) Consume(queueName string, expectedMessageCount int) {

	msgs, err := tqc.queueAdapter.Consume(
		queueName, // queue name
		"",
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	log.Printf("Export ConsumerService started for queue %s. Waiting for %d messages...", queueName, expectedMessageCount)

	for msg := range msgs {
		tqc.Lock()
		tqc.QueueMessages[queueName] = append(tqc.QueueMessages[queueName], string(msg.Body))
		tqc.Unlock()
		if len(tqc.QueueMessages[queueName]) >= expectedMessageCount {
			break
		}
	}

	log.Printf("Export ConsumerService finished for queue %s. Received %d messages.", queueName, len(tqc.QueueMessages[queueName]))
}

// WaitForMessages waits until the specified number of messages have been consumed for a specific queue.
func (tqc *TestQueueConsumer) WaitForMessages(queueName string, expectedMessageCount int) {
	tqc.WaitGroup.Add(1)
	go func() {
		defer tqc.WaitGroup.Done()
		tqc.Consume(queueName, expectedMessageCount)
	}()
	// Wait for the specified number of messages to be consumed
	tqc.WaitGroup.Wait()
}

// Wait waits until all consumers finish consuming messages.
func (tqc *TestQueueConsumer) Wait() {
	tqc.WaitGroup.Wait()
}

func (tqc *TestQueueConsumer) ClearMessages() {
	tqc.QueueMessages = make(map[string][]string)
}
