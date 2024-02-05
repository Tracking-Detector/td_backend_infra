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
	Done          chan struct{} // Add this channel
}

func NewTestQueueConsumer(queueAdapter queue.IQueueChannelAdapter) *TestQueueConsumer {
	return &TestQueueConsumer{
		queueAdapter:  queueAdapter,
		QueueMessages: make(map[string][]string),
		Done:          make(chan struct{}), // Initialize the channel
	}
}

// Consume captures the consumed RabbitMQ messages and stores them in the QueueMessages map.
func (tqc *TestQueueConsumer) Consume(queueName string, expectedMessageCount int) {
	defer close(tqc.Done) // Close the channel when done

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

		tqc.QueueMessages[queueName] = append(tqc.QueueMessages[queueName], string(msg.Body))

		if len(tqc.QueueMessages[queueName]) >= expectedMessageCount {
			break
		}
	}

	log.Printf("Export ConsumerService finished for queue %s. Received %d messages.", queueName, len(tqc.QueueMessages[queueName]))
}

func (tqc *TestQueueConsumer) ClearMessages() {
	tqc.QueueMessages = make(map[string][]string)
}
