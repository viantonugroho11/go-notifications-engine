# Menambah Handler Kafka Baru

**Pola:** Consumer memakai [go-lib/kafka](https://github.com/viantonugroho11/go-lib). Tiap handler satu file, implement `kafka.EventHandler[NotificationProducerMessage]`, dan **daftar sendiri di main** lewat `map[string]kafka.EventHandler[NotificationProducerMessage]`. Router tidak diubah.

## Langkah

### 1. Buat file handler di `handler/`

Handler harus implement `kafka.EventHandler[E]` dari go-lib: **`Name() string`** dan **`Handle(ctx, evt E, headers ...kafka.Header) kafka.Progress`**. Event `E` di-decode JSON otomatis oleh consumer.

Contoh untuk event type baru (misalnya `OrderEvent`):

```go
package handler

import (
	"context"
	"log"

	"github.com/viantonugroho11/go-lib/kafka"
)

// OrderEvent payload dari topic (JSON).
type OrderEvent struct {
	ID     string `json:"id"`
	Amount int    `json:"amount"`
}

// OrderCreator dependency (use case).
type OrderCreator interface {
	CreateFromEvent(ctx context.Context, id string, amount int) error
}

// OrderCreatedHandler menangani event order created.
type OrderCreatedHandler struct {
	creator OrderCreator
}

func NewOrderCreatedHandler(creator OrderCreator) *OrderCreatedHandler {
	return &OrderCreatedHandler{creator: creator}
}

func (h *OrderCreatedHandler) Name() string { return "order-created" }

func (h *OrderCreatedHandler) Handle(ctx context.Context, evt OrderEvent, _ ...kafka.Header) kafka.Progress {
	if err := h.creator.CreateFromEvent(ctx, evt.ID, evt.Amount); err != nil {
		log.Printf("order created handler error: %v", err)
		return kafka.Progress{Status: kafka.ProgressError, Err: err}
	}
	return kafka.Progress{Status: kafka.ProgressSuccess}
}
```

Untuk handler yang memakai **event notifikasi** (topic yang pakai `NotificationProducerMessage`), gunakan `notifEntity.NotificationProducerMessage` dan daftarkan di `RegisterConsumers` yang sama (handler key + consumer config).

### 2. Di main (consumer)

- Construct handler dengan dependency-nya.
- Daftarkan ke map dengan key unik. **Tipe map:** `map[string]kafka.EventHandler[notifEntity.NotificationProducerMessage]` untuk topic notifikasi; untuk event type lain perlu `RegisterConsumers` / constructor consumer yang menerima handler type tersebut (saat ini router hanya mendukung `NotificationProducerMessage`).
- Tambah `ConsumerConfig` untuk topic + group yang pakai handler itu.

```go
kafkaHandlers := map[string]kafka.EventHandler[notifEntity.NotificationProducerMessage]{
	"notification": eventhandler.NewNotificationUpdateHandler(notificationService),
	"sent":         eventhandler.NewNotificationSentHandler(sentSender),
}
consumerConfigs := []eventhandler.ConsumerConfig{
	{Topic: cfg.Kafka.Topic, GroupID: cfg.Kafka.GroupID, HandlerKey: "notification"},
	{Topic: cfg.Kafka.TopicSent, GroupID: cfg.Kafka.GroupID + "-sent", HandlerKey: "sent"},
}
consumers, err := eventhandler.RegisterConsumers(ctx, cfg.KafkaBrokersList(), consumerConfigs, kafkaHandlers)
```

**Tidak perlu** ubah `router.go` — cukup tambah entry di map dan config di main.

Satu `HandlerKey` = satu handler. Satu `ConsumerConfig` = satu consumer (topic + group). Progress: `ProgressSuccess` (commit), `ProgressError` (retry), `ProgressSkip`/`ProgressDrop` (commit tanpa retry).
