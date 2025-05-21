package messaging

import (
	"context"
	"log"
	"sync"
)

// Consumer defines the interface that all consumers must implement
type Consumer interface {
	Start(ctx context.Context) error
	Name() string
}

// ConsumerManager handles multiple consumers
type ConsumerManager struct {
	consumers []Consumer
	wg        sync.WaitGroup
}

// NewConsumerManager creates a new consumer manager
func NewConsumerManager(consumers ...Consumer) *ConsumerManager {
	return &ConsumerManager{
		consumers: consumers,
	}
}

// StartAll starts all registered consumers
func (m *ConsumerManager) StartAll(ctx context.Context) error {
	for _, consumer := range m.consumers {
		m.wg.Add(1)
		go func(c Consumer) {
			defer m.wg.Done()

			// Start the consumer with error recovery
			for {
				select {
				case <-ctx.Done():
					log.Printf("Consumer %s shutting down...", c.Name())
					return
				default:
					if err := c.Start(ctx); err != nil {
						log.Printf("Error in consumer %s: %v. Restarting...", c.Name(), err)
						continue
					}
					return
				}
			}
		}(consumer)
	}

	return nil
}

// WaitForShutdown waits for all consumers to finish
func (m *ConsumerManager) WaitForShutdown() {
	m.wg.Wait()
}
