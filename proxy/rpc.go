package rpc

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/streadway/amqp"
)

const (
	retrySleep = 2 * time.Second
)

type RpcClient struct {
	url, exchange string
	cfg           *amqp.Config

	conn       *amqp.Connection
	channel    *amqp.Channel
	replyQueue amqp.Queue

	confirmations <-chan amqp.Confirmation
	deliveries    <-chan amqp.Delivery

	m  *chanMap
	mu sync.Mutex
}

func NewRPCClient(url, exchange string, cfg *amqp.Config) *RpcClient {
	return &RpcClient{
		m:        NewChanMap(),
		url:      url,
		exchange: exchange,
		cfg:      cfg,
	}
}

func (c *RpcClient) Init() error {
	var err error

	if c.conn, err = amqp.DialConfig(c.url, *c.cfg); err != nil {
		return fmt.Errorf("Failed connecting to AMQP: %v", err)
	}

	if c.channel, err = c.conn.Channel(); err != nil {
		return fmt.Errorf("Failed opening a channel: %v", err)
	}

	if err = c.channel.Confirm(false); err != nil {
		return fmt.Errorf("Failed putting channel into confirm mode: %v", err)
	}
	c.confirmations = c.channel.NotifyPublish(make(chan amqp.Confirmation, 1))

	if err := c.channel.ExchangeDeclare(
		c.exchange, // name
		"direct",   // kind
		true,       // durable
		false,      // autoDelete
		false,      // internal
		false,      // noWait
		nil,        // args
	); err != nil {
		return fmt.Errorf("ExchangeDeclare: %v", err)
	}

	if c.replyQueue, err = c.channel.QueueDeclare(
		"",    // name
		false, // durable
		true,  // delete when usused
		true,  // exclusive
		false, // noWait
		nil,   // arguments
	); err != nil {
		return fmt.Errorf("Failed to declare a queue: %v", err)
	}

	if c.deliveries, err = c.channel.Consume(
		c.replyQueue.Name, // queue
		"",                // consumer
		true,              // auto-ack
		false,             // exclusive
		false,             // no-local
		false,             // no-wait
		nil,               // args
	); err != nil {
		return fmt.Errorf("Failed to register a consumer: %v", err)
	}

	go c.loop()

	return nil
}

// Call sends an RPC message to the provited routing_key (key) with empty body. It expectes to
// receive a response on the response channel.
func (c *RpcClient) Call(key string) (<-chan []byte, error) {
	id := RandomString(32)

	ch := make(chan []byte)
	c.m.Add(id, ch)

	c.mu.Lock()
	defer c.mu.Unlock()

	if err := c.channel.Publish(
		c.exchange,
		key,
		true,  // mandatory
		false, // immidiate
		amqp.Publishing{
			CorrelationId: id,
			ReplyTo:       c.replyQueue.Name,
		}); err != nil {
		c.m.Delete(id)
		return nil, err
	}

	if confirmed := <-c.confirmations; !confirmed.Ack {
		c.m.Delete(id)
		return nil, errors.New("Nack received")
	}

	return ch, nil
}

func (c *RpcClient) loop() {
	returns := c.channel.NotifyReturn(make(chan amqp.Return))
	errors := c.channel.NotifyClose(make(chan *amqp.Error))

	for {
		select {
		case ret, ok := <-returns:
			if !ok {
				log.Printf("Returns channel closed")
				return
			}
			log.Printf("Published message returned: %v", ret)
			if ch, ok := c.m.Get(ret.CorrelationId); ok {
				ch <- nil
			}
			c.m.Delete(ret.CorrelationId)
		case ret, ok := <-c.deliveries:
			if !ok {
				log.Printf("Messages channel closed")
				return
			}
			if ch, ok := c.m.Get(ret.CorrelationId); ok {
				ch <- ret.Body
			}
			c.m.Delete(ret.CorrelationId)
		case err, ok := <-errors:
			if !ok {
				log.Printf("Errors channel closed")
				return
			}
			log.Printf("Shutdown: %v", err)
			c.shutdown()
			for {
				log.Printf("Reconnecting...")
				if err := c.Init(); err != nil {
					log.Printf("amqpClient init: %v", err)
					time.Sleep(retrySleep)
					continue
				}
				log.Printf("Connected")
				return
			}
		}
	}
}

func (c *RpcClient) shutdown() {
	if err := c.channel.Close(); err != nil {
		log.Printf("Failed closing channel: %v", err)
	}
	if err := c.conn.Close(); err != nil {
		log.Printf("Failed closing connection: %v", err)
	}
}
