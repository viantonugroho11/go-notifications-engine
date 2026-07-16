# Arsitektur — go-notifications-engine

> Level-0 notification service: semua pengiriman notifikasi (email, push, SMS, dll.) melewati service ini.

---

## Gambaran Sistem

```
Caller (service lain)
        │
        ▼  POST /notifications
┌──────────────────────────────┐
│   HTTP API (Echo)            │  ← cmd/app
│   transport/apis/            │
│   usecase/notifications      │
│   state machine              │
└──────────┬───────────────────┘
           │  GORM + PostgreSQL
           ▼
   ┌───────────────────┐
   │   PostgreSQL       │
   │  notifications     │
   │  notification_logs │
   │  notification_     │
   │    templates       │
   │  notification_     │
   │    inbox           │
   │  users             │
   └──────┬────────────┘
          │  Kafka publish (setiap state change)
          ▼
   ┌─────────────┐
   │   Kafka      │  topic: user-events / sent
   └──────┬──────┘
          │
   ┌──────▼──────────────────────────────────────────────────┐
   │  Consumer (cmd/consumer -consumer=X) — proses terpisah  │
   │                                                          │
   │  "notification" consumer                                 │
   │    └─ update state notification via HTTP API             │
   │                                                          │
   │  "sent" consumer                                         │
   │    └─ resolve sendTo dari person service                 │
   │    └─ render template                                    │
   │    └─ kirim email (SMTP) / push (FCM)                   │
   │    └─ update log state via HTTP API                      │
   └─────────────────────────────────────────────────────────┘
```

---

## Komponen Utama

| Komponen | Lokasi | Fungsi |
|---|---|---|
| HTTP API | `cmd/app`, `transport/apis/` | Entry point CRUD + trigger pengiriman |
| Consumer binary | `cmd/consumer`, `transport/event/kafka/` | Async pipeline: update state, render, kirim |
| State Machine (Notification) | `usecase/notifications/states/` | Transisi state: CREATED → PROCESSING → COMPLETED/FAILED |
| State Machine (Log) | `usecase/notificationlogs/states/` | Transisi state log per user |
| Notification Client | `client/notification/` | HTTP client consumer → API (self-call, lihat keputusan desain) |
| Person Client | `client/person/` | Resolve email/token/phone per user dari person service |
| Email Client | `client/email/` | Adapter SMTP via gomail |
| Firebase Client | `client/firebase/` | Adapter Firebase FCM |

---

## Flow Pengiriman Notifikasi

### Step-by-step

```
1. Caller POST /notifications
   ├─ validasi request
   ├─ buat Notification (state: CREATED)
   ├─ buat NotificationLog per user (state: PENDING)
   ├─ simpan ke PostgreSQL (FullSaveAssociations)
   └─ publish event ke Kafka topic "user-events"

2. Consumer "notification" consume event
   ├─ filter: skip jika action=DELETE atau state=SCHEDULED/PROCESSING
   └─ panggil PUT /notifications/:id → update state → PROCESSING
      └─ state machine publish event baru ke Kafka

3. Consumer "sent" consume event
   ├─ FetchPerson: resolve sendTo dari person service
   │   ├─ email → person.Email
   │   ├─ push  → device token dengan LastActiveAt terbaru
   │   └─ sms/wa/telegram → person.Phone
   ├─ GenerateMessage: fetch template, render subject+body dengan data
   │   └─ update log state → PROCESSING
   └─ Send: kirim ke provider
       ├─ email → SMTP (gomail)
       ├─ push  → Firebase FCM
       └─ update log state → COMPLETED (atau FAILED) + external_ref
```

### State Machine — Notification

```
CREATED
   │
   ├── (schedule_at ada) → SCHEDULED
   │                           │
   │                     (waktu tiba) → PROCESSING
   │
   └── (langsung) → PROCESSING
                         │
                    (kirim berhasil) → COMPLETED
                         │
                    (kirim gagal)  → FAILED
```

### State Machine — Notification Log (per user)

```
PENDING → PROCESSING → COMPLETED
                     → FAILED
```

> **Penting:** State `PROCESSING` di-set saat template selesai di-render (siap kirim). State `COMPLETED`/`FAILED` di-set setelah provider eksternal merespons.

---

## Kafka Event Schema

Semua event menggunakan envelope standar:

```json
{
  "resource_id": "<notification-id>",
  "meta": {
    "event_id": "<uuid>",
    "event_timestamp": "2026-07-17T10:00:00Z",
    "action": "INSERT",
    "resource": "Notification",
    "message_schema_version": 1
  },
  "before": { },
  "after": {
    "notification_id": "...",
    "event_key": "order.created",
    "notification_template_id": "...",
    "channel": "email",
    "category": "transactional",
    "state": "CREATED",
    "data": { "orderId": "ORD-001" },
    "notification_logs": [
      {
        "id": "...",
        "user_id": "...",
        "state": "PENDING"
      }
    ]
  }
}
```

**Action values:** `INSERT`, `UPDATE`, `DELETE`

Consumer handler baca dari `evt.Meta.Action` (bukan top-level `action`).

---

## Keputusan Desain

### Consumer memanggil API via HTTP (self-call), bukan langsung ke repository

**Keputusan:** Consumer tidak boleh impor repository atau usecase secara langsung. Semua mutasi state lewat REST API (`client/notification`).

**Alasan:**

1. **Single source of truth untuk business logic.** State machine, validasi, dan Kafka publish semuanya ada di `usecase/notifications`. Bypass langsung ke repo → logic ter-skip, state bisa invalid.

2. **Consumer adalah binary terpisah.** `cmd/consumer` dan `cmd/app` deploy berbeda. Consumer hanya init dependency minimal (notificationRepo, templateRepo) — tidak ada akses ke semua usecase.

3. **Konsistensi API kontrak.** Setiap update lewat handler yang sama dengan caller eksternal. Fix di HTTP path otomatis fix di consumer path.

4. **Failure isolation.** API down → consumer error → Kafka retry otomatis. Tidak ada partial state masuk DB tanpa validasi.

**Trade-off yang diterima:**

- Tambah satu network hop (~1–5ms) per step consumer. Acceptable karena pipeline async.
- Consumer bergantung pada API service up. Mitigasi: consumer return `ProgressError` → Kafka retry.

**Aturan:** Jangan inject `NotificationRepository` atau usecase langsung ke consumer event handler. Kalau butuh data baru via consumer, tambah endpoint di API layer.

---

### Template rendering di entity layer

Go `text/template` dengan `Option("missingkey=zero")` — variable yang tidak ada di `data` di-render sebagai string kosong, tidak error.

```go
// Contoh template
"Halo {{.customerName}}, pesanan {{.orderId}} sudah diterima."

// Contoh data
{"customerName": "Budi", "orderId": "ORD-001"}
```

**Trade-off:** Tidak support logic kompleks (loop, conditional). Cocok untuk notifikasi sederhana. Kalau butuh logic kompleks, pertimbangkan Sprig atau library template lain.

---

## Data Model

### notifications

| Kolom | Tipe | Keterangan |
|---|---|---|
| `id` | UUID | Primary key |
| `event_key` | varchar | Identifier event (e.g. `order.created`) |
| `notification_template_id` | UUID | FK ke `notification_templates` |
| `data` | jsonb | Variable untuk template rendering |
| `channel` | varchar | `email`, `push`, `sms`, dll. |
| `category` | varchar | `transactional`, `promo`, `system`, `other` |
| `state` | varchar | State machine notification |
| `schedule_at` | timestamp | Waktu pengiriman terjadwal (nullable) |
| `created_by` | varchar | Identifier caller/sistem |
| `created_at` | timestamp | Waktu dibuat |
| `updated_at` | timestamp | Waktu update terakhir |

### notification_logs

| Kolom | Tipe | Keterangan |
|---|---|---|
| `id` | UUID | Primary key |
| `notification_id` | UUID | FK ke `notifications` |
| `user_id` | UUID | Target user |
| `send_to` | text | Alamat pengiriman (email/token/phone) — diisi saat FetchPerson |
| `rendered_subject` | text | Subject setelah template rendering |
| `rendered_message` | text | Body setelah template rendering |
| `state` | varchar | State machine log |
| `retry_count` | int | Jumlah retry |
| `error_message` | text | Pesan error dari provider |
| `external_ref` | text | Message ID dari provider (FCM message ID, SES message ID) |
| `sent_at` | timestamp | Waktu berhasil terkirim |
| `created_at` | timestamp | Waktu dibuat |

### notification_templates

| Kolom | Tipe | Keterangan |
|---|---|---|
| `id` | UUID | Primary key |
| `name` | varchar | Nama template |
| `channel` | varchar | Channel target |
| `template_type` | varchar | `transactional`, `promo`, `system`, `other` |
| `subject` | text | Template subject (Go template) |
| `body` | text | Template body (Go template) |
| `payload_schema` | jsonb | JSON schema validasi payload `data` |

### notification_inbox

| Kolom | Tipe | Keterangan |
|---|---|---|
| `id` | UUID | Primary key |
| `user_id` | UUID | Pemilik inbox |
| `notification_log_id` | UUID | FK ke `notification_logs` |
| `subject` | text | Subject notifikasi (denormalized) |
| `message` | text | Body notifikasi (denormalized) |
| `is_read` | boolean | Status baca |
| `read_at` | timestamp | Waktu dibaca |
| `created_at` | timestamp | Waktu dibuat |

---

## Aturan Dependency

```
transport/apis    → usecase → repository (interface) → entity
transport/event   → usecase → client/ → entity
bootstrap/        → semua (satu-satunya tempat wiring konkret)
entity/           → tidak boleh impor layer lain
client/           → entity (tidak boleh impor usecase/repository)
```

---

## Analisis Skala

### Estimasi kapasitas default (single instance)

| Layer | Kapasitas | Bottleneck |
|---|---|---|
| Echo HTTP server | ~10.000 req/s | CPU saat marshal JSON besar |
| PostgreSQL (GORM) | ~1.000–5.000 write/s | Connection pool habis saat lonjakan |
| Kafka producer | Sangat tinggi (batch async) | Bukan bottleneck |
| Kafka consumer | 1 consumer per partition | Throughput dibatasi jumlah partition |
| Email (SMTP) | Tergantung provider (SES: ~14/s default) | Rate limit provider |
| Firebase FCM | ~500–1.000 req/s per connection | Butuh batching untuk throughput tinggi |

### Titik lemah

**1. Consumer single-threaded per key.**
Scale: naikkan partition Kafka, deploy multiple consumer instance (consumer group otomatis balance).

**2. External call blocking.**
Consumer tunggu response email/FCM sebelum proses message berikutnya. Provider lambat → throughput turun.
Mitigasi: timeout agresif (email: 5s, FCM: 3s), goroutine pool untuk parallel log.

**3. Template tidak di-cache.**
Fetch template via HTTP setiap render. Volume tinggi → N+1 HTTP call.
Mitigasi: cache di Redis, key `template:{id}`, TTL 5–10 menit.

**4. Tidak ada rate limiter.**
Caller bisa flood `POST /notifications`.
Mitigasi: Echo rate limiter middleware per IP atau API key.

**5. Tidak ada Dead Letter Queue.**
Message gagal terus di-retry. Bug permanen → consumer stuck.
Mitigasi: max retry di go-lib/kafka, route ke DLQ topic setelah max retry.

### Rekomendasi untuk traffic tinggi (>10.000 notif/menit)

```
1. Partition Kafka ≥ 4 per topic + ≥ 2 consumer instance
2. PostgreSQL connection pool: MaxOpenConns=50, MaxIdleConns=10, ConnMaxLifetime=5m
3. Redis cache untuk template (TTL 5 menit)
4. Timeout eksplisit: notification client=10s, person=5s, email=5s, FCM=3s
5. FCM Multicast API untuk push ke banyak device sekaligus
6. Dead Letter Queue untuk message yang melebihi max retry
```

---

## Dependency Eksternal

| Dependency | Tujuan | Fallback |
|---|---|---|
| PostgreSQL | Primary store | Tidak ada — service tidak bisa berjalan tanpa DB |
| Kafka | Event bus | Kalau Kafka down, create tetap jalan; publish error di-log, tidak gagalkan request |
| Redis | Cache (belum aktif di hot path) | Optional |
| Person service | Resolve email/token/phone user | Consumer error → Kafka retry |
| Email provider (SMTP) | Kirim email | Retry via Kafka |
| Firebase FCM | Kirim push notification | Retry via Kafka |

---

## Channel yang Direncanakan

| Channel | Status | Provider |
|---|---|---|
| `email` | Aktif | SMTP (gomail) |
| `push` | Aktif | Firebase FCM |
| `sms` | Planned | Twilio / vendor lokal |
| `whatsapp` | Planned | WhatsApp Business API |
| `telegram` | Planned | Telegram Bot API |
| `line` | Planned | LINE Messaging API |

Untuk menambah channel baru: implementasi `ChannelSender` interface di `client/`, register di `notification_send_notification_usecase.go` switch case, tambah config dan client init di `bootstrap/consumer.go`.

---

## Backlog & Known Issues

| Issue | Severity | Status |
|---|---|---|
| Tidak ada autentikasi/otorisasi di HTTP routes | P1 | Open |
| Template tidak di-cache Redis | P2 | Open |
| Tidak ada DLQ untuk message gagal berulang | P2 | Open |
| Timeout eksplisit di semua HTTP client | P2 | Open |
| Rate limiter API belum ada | P2 | Open |
| SMS/WA/Telegram channel belum diimplementasi | P2 | Planned |
| Scheduled notification belum ada worker | P2 | Planned |
| Tidak ada unit/integration test | P2 | Open |
