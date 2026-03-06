package kafka

import (
	"go-boilerplate-clean/internal/infrastructure/broker"
)

// Consumers mengelola banyak consumer (go-lib/kafka); Close() menutup semuanya.
type Consumers struct {
	List []broker.Consumer
}

// Close menutup semua consumer.
func (c *Consumers) Close() error {
	for _, cons := range c.List {
		_ = cons.Close()
	}
	return nil
}
