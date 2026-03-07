# Menambah Handler Kafka Baru

**Pola:** Mirip **apis** (satu tempat routing). Consumer memakai [go-lib/kafka](https://github.com/viantonugroho11/go-lib).

- **Routing handler:** `internal/transport/event/kafka/router.go` — `EventServices`, `Handlers(svc)`, `Keys` / `AvailableKeys()`. Tambah key + satu baris di `Handlers` = consumer baru.
- **Config Kafka (topic, group):** `internal/infrastructure/broker/kafka/registry.go` — `GetConsumerConfigByKey(cfg, key)`. Tambah case untuk key baru.
- **Handler:** satu file per handler di `handler/`, implement `EventHandler[NotificationProducerMessage]`.

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

### 2. Daftarkan routing (event) + config (broker)

- **event/kafka/router.go:** tambah `const KeyX = "x"`, tambah `KeyX` ke slice `Keys`, tambah entry di `Handlers(svc)`: `KeyX: handler.NewXHandler(svc.X)`.
- **broker/kafka/registry.go:** tambah case di `GetConsumerConfigByKey(cfg, key)` untuk key baru (Topic, GroupID, ClientID dari cfg).
- **EventServices:** jika butuh dependency baru, tambah field di struct dan isi dari bootstrap.

Pemanggil: `brokerkafka.NewConsumerRunner(cfg, consumerKey, eventkafka.EventServices{...})`. Flag `-consumer` menentukan consumer mana yang jalan.

Lifecycle: `Start(ctx)` lalu `Close()`. Progress: `ProgressSuccess` (commit), `ProgressError` (retry), `ProgressSkip`/`ProgressDrop` (commit tanpa retry).

---

### 3. Event message tipe lain (bukan NotificationProducerMessage)

Kalau topic pakai **message type beda** (mis. `OrderEvent`, `PaymentEvent`):

1. **Handler + tipe event** di `handler/` (sama seperti di atas), dengan `EventHandler[OrderEvent]` dsb.
2. **Router terpisah** di `event/kafka/`: buat mis. `router_order.go`:
   - `const KeyOrderCreated = "order_created"`
   - `var KeysOrder = []string{KeyOrderCreated}`
   - `AllKeys` di `router.go` digabung: `var AllKeys = append(Keys, KeysOrder...)`
   - `OrderEventServices` struct + `HandlersOrder(svc OrderEventServices) map[string]EventHandler[OrderEvent]`
3. **Registry:** tambah case di `GetConsumerConfigByKey(cfg, key)` untuk `KeyOrderCreated` (topic, group).
4. **Jalankan consumer:** pakai **generic runner** — jangan `NewConsumerRunner` (itu untuk notifikasi), tapi `NewConsumerRunnerFor[OrderEvent](cfg, consumerKey, eventkafka.HandlersOrder(svc))`.

Contoh di bootstrap (untuk consumer order):

```go
orderSvc := usecase.NewOrderService(...)
handlersOrder := eventkafka.HandlersOrder(eventkafka.OrderEventServices{Order: orderSvc})
consumers, err := brokerkafka.NewConsumerRunnerFor[entity.OrderEvent](cfg, "order_created", handlersOrder)
```

Ringkas: **message sama** (NotificationProducerMessage) → pakai `Handlers(svc)` + `NewConsumerRunner`. **Message beda** → buat router + `HandlersX(svc)` + `NewConsumerRunnerFor[TipeEvent]`.
