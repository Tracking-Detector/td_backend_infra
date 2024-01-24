package queue

import "github.com/streadway/amqp"

type RabbitMQChannelAdapter struct {
	ch *amqp.Channel
}

func NewRabbitMQChannelAdapter(ch *amqp.Channel) *RabbitMQChannelAdapter {
	return &RabbitMQChannelAdapter{
		ch: ch,
	}
}

func (a *RabbitMQChannelAdapter) QueueDeclare(name string, durable bool, autoDelete bool, exclusive bool, noWait bool, args amqp.Table) (amqp.Queue, error) {
	return a.ch.QueueDeclare(name, durable, autoDelete, exclusive, noWait, args)
}

func (a *RabbitMQChannelAdapter) Publish(exchange string, key string, mandatory bool, immediate bool, msg amqp.Publishing) error {
	return a.ch.Publish(exchange, key, mandatory, immediate, msg)
}
