package queue

import "github.com/streadway/amqp"

type IQueueChannelAdapter interface {
	QueueDeclare(name string, durable bool, autoDelete bool, exclusive bool, noWait bool, args amqp.Table) (amqp.Queue, error)
	Publish(exchange string, key string, mandatory bool, immediate bool, msg amqp.Publishing) error
}
