# go-notifications-engine

Notification engine untuk mengelola pengiriman notifikasi multi-channel (email, push/FCM, SMS, WhatsApp, dll.) dengan template, per-user delivery log, dan inbox.

Stack: **Go · Echo · GORM + PostgreSQL · Kafka · Redis · Firebase FCM · SMTP**

---

## Fitur

| Fitur | Keterangan |
|---|---|
| Multi-channel | Email, Push (FCM), SMS, WhatsApp, Telegram, Line (Email & Push aktif) |
| Template engine | Go `text/template` — subject & body dengan variable substitution dari `data` payload |
| Per-user log | Setiap notifikasi punya satu log per user (`notification_logs`) dengan state machine |
| Inbox | Record in-app per user untuk notifikasi push/in-app |
| Kafka pipeline | Producer publish event setiap state change; consumer proses async delivery |
| State machine | Explicit state transition untuk notification dan notification log |
| Scheduled notification | Field `schedule_at` untuk notifikasi terjadwal |
| Self-contained consumer | Consumer binary terpisah, kirim email/push tanpa sharing DB connection dengan API |

---

## Persyaratan

- **Go** 1.24+
- **PostgreSQL** 14+
- **Kafka** (ZooKeeper atau KRaft)
- **Redis** 6+
- **Firebase project** (untuk push notification)
- **SMTP server** (untuk email)

---

## Quick Start

### 1. Clone dan download dependency

```bash
git clone <repository-url>
cd go-notifications-engine
go mod download
```

### 2. Jalankan infrastruktur lokal

```bash
docker compose up -d --build
```

Port default: App **8080** · PostgreSQL **5432** · Kafka **9092** · Redis **6379**

### 3. Konfigurasi

Salin dan edit konfigurasi:

```bash
cp configs/config.yaml config/config.yaml
```

Edit `config/config.yaml` minimal:

```yaml
port: "8080"
notification_base_url: "http://localhost:8080"

dbhost: "localhost"
dbport: "5432"
dbuser: "postgres"
dbpassword: "postgres"
dbname: "appdb"

kafkabrokers: "localhost:9092"

redisaddr: "localhost:6379"

# Email (opsional)
emailhost: "smtp.gmail.com"
emailport: 587
emailuser: "noreply@example.com"
emailpassword: "secret"

# Firebase (opsional, untuk push)
fcmprojectid: "my-firebase-project"
```

### 4. Jalankan HTTP server

```bash
go run ./cmd/app
```

Server otomatis jalankan `AutoMigrate` saat startup.

### 5. Jalankan Kafka consumer (proses terpisah)

```bash
# Consumer update state notification
go run ./cmd/consumer -consumer notification

# Consumer kirim email/push
go run ./cmd/consumer -consumer sent
```

---

## HTTP API

Base URL: `http://localhost:8080`

### Health Check

```
GET /healthz
→ 200 "ok"
```

### Notifications

| Method | Path | Deskripsi |
|---|---|---|
| `POST` | `/notifications` | Buat notifikasi baru (bisa ke banyak user sekaligus) |
| `GET` | `/notifications` | List notifikasi (dengan filter & pagination) |
| `GET` | `/notifications/:id` | Detail notifikasi + logs |
| `PUT` | `/notifications/:id` | Update notifikasi |
| `DELETE` | `/notifications/:id` | Hapus notifikasi |

**Buat notifikasi:**

```bash
curl -X POST http://localhost:8080/notifications \
  -H "Content-Type: application/json" \
  -d '{
    "event_key": "order.created",
    "notification_template_id": "<template-uuid>",
    "channel": "email",
    "category": "transactional",
    "data": {
      "orderId": "ORD-001",
      "customerName": "Budi"
    },
    "user_ids": ["<user-uuid-1>", "<user-uuid-2>"]
  }'
```

**Query params untuk GET list:**

| Param | Tipe | Contoh |
|---|---|---|
| `event_key` | string | `order.created` |
| `channel` | string | `email` |
| `category` | string (comma) | `transactional,promo` |
| `state` | string (comma) | `CREATED,PROCESSING` |
| `ids` | string (comma) | `uuid1,uuid2` |
| `page` | int | `1` |
| `limit` | int | `20` |
| `offset` | int | `0` |

### Notification Templates

| Method | Path | Deskripsi |
|---|---|---|
| `POST` | `/notification-templates` | Buat template |
| `GET` | `/notification-templates` | List template |
| `GET` | `/notification-templates/:id` | Detail template |
| `PUT` | `/notification-templates/:id` | Update template |
| `DELETE` | `/notification-templates/:id` | Hapus template |

**Buat template:**

```bash
curl -X POST http://localhost:8080/notification-templates \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Order Created Email",
    "channel": "email",
    "template_type": "transactional",
    "subject": "Pesanan {{.orderId}} berhasil dibuat",
    "body": "Halo {{.customerName}}, pesanan Anda dengan ID {{.orderId}} sudah diterima."
  }'
```

Template menggunakan Go `text/template` syntax. Variable diisi dari field `data` pada saat notifikasi dibuat.

### Notification Logs

| Method | Path | Deskripsi |
|---|---|---|
| `GET` | `/notification-logs` | List logs (filter by user, notification, state) |
| `GET` | `/notification-logs/:id` | Detail log |
| `PUT` | `/notification-logs/:id` | Update log (state, error, external ref) |

### Notification Inbox

| Method | Path | Deskripsi |
|---|---|---|
| `GET` | `/notification-inbox` | List inbox per user |
| `GET` | `/notification-inbox/:id` | Detail inbox item |
| `PUT` | `/notification-inbox/:id` | Update (mark as read) |

### Users

| Method | Path | Deskripsi |
|---|---|---|
| `POST` | `/users` | Buat user |
| `GET` | `/users` | List user |
| `GET` | `/users/:id` | Detail user |
| `PUT` | `/users/:id` | Update user |
| `DELETE` | `/users/:id` | Hapus user |

---

## Kafka Event Schema

Semua event mengikuti schema standar:

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
    "notification_logs": [...]
  }
}
```

**Topics:**

| Topic | Consumer key | Fungsi |
|---|---|---|
| `user-events` | `notification` | Update state notification dari event CDC/aplikasi lain |
| *(sent topic)* | `sent` | Trigger delivery email/push setelah state PROCESSING |

---

## Channel yang Didukung

| Channel | Nilai | Status |
|---|---|---|
| Email | `email` | Aktif (SMTP via gomail) |
| Push Notification | `push` | Aktif (Firebase FCM) |
| SMS | `sms` | Didefinisikan, belum diimplementasi |
| WhatsApp | `whatsapp` | Didefinisikan, belum diimplementasi |
| Telegram | `telegram` | Didefinisikan, belum diimplementasi |
| Line | `line` | Didefinisikan, belum diimplementasi |

---

## State Machine

### Notification State

```
CREATED → SCHEDULED (jika ada schedule_at)
        → PROCESSING
          → SENT
            → COMPLETED
            → FAILED
```

### Notification Log State (per user)

```
PENDING → PROCESSING → COMPLETED
                     → FAILED
```

---

## Struktur Proyek

```
cmd/
  app/                      HTTP server entrypoint
  consumer/                 Kafka consumer entrypoint (-consumer <key>)
internal/
  bootstrap/                Wiring semua dependency
  config/                   Config struct (App, DB, Kafka, Redis, Email, FCM)
  entity/                   Domain struct — bebas dari framework & DB tag
    notifications/
    notificationlogs/
    notificationtemplates/
    notificationinbox/
    users/
  repository/               Interface + GORM model + postgres impl per aggregate
  usecase/                  Business logic & state machine
    notifications/states/   State machine explicit (created, processing, sent, …)
    notificationlogs/states/
    event/                  Usecase untuk event pipeline (FetchPerson, GenerateMessage, Send)
  transport/
    apis/                   Echo handler, router, DTO
    event/kafka/            Kafka consumer handler & routing
  infrastructure/
    broker/kafka/           Producer, consumer runner, registry
    cache/redis/            Redis client
    database/postgres/      GORM connect + AutoMigrate
  client/
    notification/           HTTP client ke API sendiri (self-call dari consumer)
    person/                 HTTP client ke person service (resolve email/token/phone)
    email/                  SMTP adapter (gomail)
    firebase/               FCM adapter
configs/
  config.yaml               Baseline konfigurasi lokal
db/
  notifications.sql         DDL tabel notifications
  notification-logs.sql     DDL tabel notification_logs
  openapi.yaml              OpenAPI 3.0 spec lengkap
docs/
  architecture.md           Arsitektur detail, keputusan desain, analisis skala
```

---

## Konfigurasi Lengkap

| Key | Env Var | Default | Keterangan |
|---|---|---|---|
| `port` | `PORT` | `8080` | Port HTTP server |
| `notification_base_url` | `NOTIFICATION_BASE_URL` | `http://localhost:8080` | Base URL API (untuk self-call dari consumer) |
| `dbhost` | `DB_HOST` | `postgres` | PostgreSQL host |
| `dbport` | `DB_PORT` | `5432` | PostgreSQL port |
| `dbuser` | `DB_USER` | `postgres` | PostgreSQL user |
| `dbpassword` | `DB_PASSWORD` | `postgres` | PostgreSQL password |
| `dbname` | `DB_NAME` | `appdb` | PostgreSQL database name |
| `databaseurl` | `DATABASE_URL` | *(kosong)* | DSN lengkap — override semua field DB di atas |
| `kafkabrokers` | `KAFKA_BROKERS` | `kafka:9092` | Kafka broker address |
| `kafkatopic` | `KAFKA_TOPIC` | `user-events` | Topic utama |
| `redisaddr` | `REDIS_ADDR` | `redis:6379` | Redis address |
| `redispassword` | `REDIS_PASSWORD` | *(kosong)* | Redis password |
| `emailhost` | `EMAIL_HOST` | *(kosong)* | SMTP host |
| `emailport` | `EMAIL_PORT` | `587` | SMTP port |
| `emailuser` | `EMAIL_USER` | *(kosong)* | SMTP username |
| `emailpassword` | `EMAIL_PASSWORD` | *(kosong)* | SMTP password |
| `fcmprojectid` | `FCM_PROJECT_ID` | *(kosong)* | Firebase project ID |

---

## Development

```bash
# Format
go fmt ./...

# Lint
go vet ./...

# Build semua binary
go build ./...

# Jalankan tests
go test ./...

# Tidy dependency
go mod tidy
```

---

## Troubleshooting

| Gejala | Solusi |
|---|---|
| `config load` error | Pastikan file `config/config.yaml` ada atau env var di-set |
| PostgreSQL connection failed | Cek `dbhost`, `dbport`, credentials; pastikan Compose running |
| Kafka producer error | Cek `kafkabrokers`, topic sudah dibuat di broker |
| Consumer `-consumer sent` gagal startup | Pastikan `kafkatopicsent` di-config atau topic tersedia |
| Email tidak terkirim | Cek `emailhost`, `emailuser`, `emailpassword`; consumer log akan print error |
| Push tidak terkirim | Cek `fcmprojectid`; pastikan service account credential tersedia |
| Consumer stuck / tidak proses | Cek `notification_base_url` — consumer self-call ke API, pastikan API up |

---

## Dokumentasi Lanjutan

- [Arsitektur & Keputusan Desain](docs/architecture.md) — flow lengkap, keputusan self-call HTTP, analisis skala
- [OpenAPI Spec](db/openapi.yaml) — kontrak API lengkap dengan request/response schema
- [Contributing](CONTRIBUTING.md) — branch workflow, konvensi kode

---

## License

MIT
