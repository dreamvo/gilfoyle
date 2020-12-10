package mocks

import (
	"github.com/streadway/amqp"
	"github.com/stretchr/testify/mock"
)

type MockedChannel struct {
	mock.Mock
}

func (m *MockedChannel) Publish(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error {
	args := m.Called(exchange, key, mandatory, immediate, msg)
	return args.Error(0)

}

func (m *MockedChannel) QueueDeclare(name string, durable, autoDelete, exclusive, noWait bool, args amqp.Table) (amqp.Queue, error) {
	arguments := m.Called(name, durable, autoDelete, exclusive, noWait)
	return amqp.Queue{}, arguments.Error(0)
}

func (m *MockedChannel) Consume(queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args amqp.Table) (<-chan amqp.Delivery, error) {
	arguments := m.Called(queue, consumer, autoAck, exclusive, noLocal, noWait)
	return make(chan amqp.Delivery), arguments.Error(0)
}
