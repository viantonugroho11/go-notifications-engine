package broker

import "context"

// Producer mendefinisikan kontrak untuk message broker producer (generic, raw bytes).
// Untuk producer typed per topic (mis. notifikasi), pakai implementasi go-lib di subpackage kafka.
type Producer interface {
	Publish(ctx context.Context, topic string, key, value []byte) (partition int32, offset int64, err error)
	Close() error
}

// ConsumerRunner menjalankan satu consumer sampai context dibatalkan.
// Routing handler di transport/event/kafka (Handlers); config di broker/kafka (registry).
type ConsumerRunner interface {
	Start(ctx context.Context)
	Close() error
}
