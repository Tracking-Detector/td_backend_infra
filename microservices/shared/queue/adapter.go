package queue

import "github.com/streadway/amqp"

type IQueueChannelAdapter interface {
	QueueDeclare(name string, durable bool, autoDelete bool, exclusive bool, noWait bool, args amqp.Table) (amqp.Queue, error)
	Publish(exchange string, key string, mandatory bool, immediate bool, msg amqp.Publishing) error
	Consume(queue string, consumer string, autoAck bool, exclusive bool, noLocal bool, noWait bool, args amqp.Table) (<-chan amqp.Delivery, error)
	PurgeQueue(name string, noWait bool) (int, error)
}
