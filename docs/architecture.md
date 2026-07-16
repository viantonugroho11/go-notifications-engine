# Architecture — go-notifications-engine

> Level-0 notification service: semua pengiriman notifikasi (email, push, SMS, dll.) melewati service ini.

---

## Gambaran sistem

```
Caller (service lain)
    │
    ▼  POST /notifications
┌──────────────────────────┐
│   HTTP API (Echo)        │  ← cmd/app
│   transport/apis/        │
│   usecase/notifications  │
│   state machine          │
└──────────┬───────────────┘
           │  GORM
           ▼
      PostgreSQL
  (notifications, logs,
   templates, inbox)
           │
           │  Kafka publish (setiap state change)
           ▼
    ┌─────────────┐
    │  Kafka      │  topic: notification / sent
    └──────┬──────┘
           │
    ┌──────▼──────────────────────────────────────────────────┐
    │  Consumer proses terpisah (cmd/consumer -consumer=X)    │
    │                                                          │
    │  notification consumer  →  update state via API          │
    │  sent consumer          →  kirim email/push via API      │
    └─────────────────────────────────────────────────────────┘
```

**Komponen utama:**

| Komponen | Lokasi | Fungsi |
|---|---|---|
| HTTP API | `cmd/app`, `transport/apis` | Entry point semua operasi CRUD + trigger pengiriman |
| Consumer | `cmd/consumer`, `transport/event/kafka` | Proses async: update state, render template, kirim notif |
| State Machine | `usecase/notifications/states` | Mengatur transisi state notifikasi |
| Notification Client | `client/notification` | HTTP client consumer → API (lihat keputusan desain di bawah) |
| Person Client | `client/person` | Resolve email/token/phone per user |
| External Clients | `client/email`, `client/firebase` | Adapter ke provider eksternal |

---

## Keputusan desain

### Consumer memanggil API via HTTP, bukan langsung ke repository

**Keputusan:** Consumer (`cmd/consumer`) tidak boleh mengimpor repository atau usecase secara langsung. Semua mutasi state dilakukan lewat REST API (`client/notification`).

**Alasan:**

1. **Single source of truth untuk business logic.** State machine, validasi, dan publish Kafka event semuanya ada di `usecase/notifications`. Kalau consumer bypass lewat repo langsung, logika ini ter-skip — state bisa invalid, event tidak ter-publish, dan tidak ada audit trail.

2. **Consumer adalah proses terpisah.** `cmd/consumer` dan `cmd/app` deploy sebagai binary berbeda. Consumer tidak punya akses ke semua repo (inbox, log, template) — hanya `notificationRepo` dan `templateRepo` yang di-init di `bootstrap/consumer.go`. Dependency graph sengaja dibuat minimal.

3. **Konsistensi API kontrak.** Setiap update melewati handler yang sama dengan yang dipanggil caller eksternal. Bug yang ketahuan di HTTP path otomatis fix di consumer path juga.

4. **Isolation failure.** Kalau API service down, consumer kembali error dan Kafka retry — tidak ada partial state yang masuk langsung ke DB tanpa validasi.

**Trade-off yang diterima:**

- Tambah satu network hop (~1-5ms) per step consumer. Acceptable karena pipeline ini async dan bukan real-time.
- Consumer bergantung pada API service up. Mitigasi: consumer return `ProgressError` → Kafka retry otomatis.
- Loop dependency semu (consumer → API → DB) tapi tidak ada circular import di kode.

**Apa yang TIDAK boleh diubah:** Jangan inject `NotificationRepository` atau usecase langsung ke consumer event handler. Kalau butuh akses data yang tidak tersedia via API, tambah endpoint baru di API layer.

---

## Flow pengiriman notifikasi

```
1. Caller POST /notifications
   └─ create notification + logs di DB
   └─ publish event ke Kafka topic "notification"

2. Consumer "notification" consume event
   └─ jika state CREATED: panggil PUT /notifications/:id (update state → PROCESSING)
   └─ state machine publish event baru ke Kafka

3. Consumer "sent" consume event
   └─ resolve sendTo dari person service
   └─ render subject+body dari template
   └─ kirim via email/Firebase
   └─ update log state → COMPLETED / FAILED

4. Jika category push/in-app: buat inbox entry via POST /notification-inbox
```

### State machine notifikasi

```
CREATED → SCHEDULED → PROCESSING → SENT → COMPLETED
                                        → FAILED
```

### State log notifikasi

```
PENDING → PROCESSING → SENT → COMPLETED
                            → FAILED
```

> **Catatan:** State `SENT` pada log harus di-set hanya setelah pengiriman berhasil ke provider eksternal (email/Firebase), bukan saat template selesai di-render.

---

## Analisis skala — level-0 service

Service ini adalah **level-0**: semua notifikasi dari seluruh sistem melewati sini. Berikut analisis tiap komponen:

### Bottleneck dan kapasitas

| Layer | Estimasi kapasitas default | Risiko utama |
|---|---|---|
| Echo HTTP server | ~10.000 req/s (single instance) | CPU bound saat JSON marshal besar |
| PostgreSQL (GORM) | ~1.000–5.000 write/s (tergantung HW) | Connection pool habis saat lonjakan |
| Kafka producer | Sangat tinggi (batch async) | Tidak ada; Kafka bukan bottleneck |
| Kafka consumer | 1 consumer per partition | Throughput dibatasi jumlah partition |
| Email client | Tergantung provider (mis. SES: 14/s default, bisa ratusan/s setelah limit naik) | Rate limit provider |
| Firebase FCM | 500–1.000 req/s per connection | Perlu batching untuk throughput tinggi |

### Titik lemah yang perlu diperhatikan

**1. Consumer single-threaded per key.**
Setiap consumer key (`notification`, `sent`) berjalan di satu proses. Throughput dibatasi jumlah partition Kafka dan kecepatan pemrosesan per message. Untuk scale:
- Naikkan partition count di topic
- Deploy multiple instance consumer (Kafka consumer group otomatis balance)

**2. External call synchronous dan blocking.**
Consumer menunggu response email/Firebase sebelum lanjut ke message berikutnya. Kalau provider lambat (timeout 30s), throughput turun drastis.
- Tambah timeout agresif di HTTP client (rekomendasi: 5s untuk email, 3s untuk FCM)
- Pertimbangkan goroutine pool untuk process multiple log per notifikasi secara paralel

**3. Tidak ada Redis caching di hot path.**
Template di-fetch via HTTP setiap kali pesan di-render. Kalau volume tinggi, ini bisa jadi N+1 call ke API.
- Cache template di Redis dengan TTL 5–10 menit
- Key: `template:{id}`, invalidate saat template di-update

**4. Tidak ada rate limiter di HTTP API.**
Caller bisa flood POST /notifications tanpa batas.
- Tambah rate limiter middleware di Echo (per IP atau per API key)

**5. Tidak ada DLQ (Dead Letter Queue).**
Message yang gagal terus di-retry tanpa batas. Bila ada bug permanen (mis. template tidak ada), consumer akan stuck.
- Konfigurasi max retry di go-lib/kafka
- Route message yang melebihi max retry ke topic DLQ terpisah
- Monitor DLQ via alert

### Rekomendasi untuk traffic tinggi (>10.000 notif/menit)

```
1. Scale horizontal consumer:
   - Naikkan partition Kafka ke ≥4 per topic
   - Deploy ≥2 instance consumer (consumer group handle balance otomatis)

2. Connection pool PostgreSQL:
   - Set MaxOpenConns: 25-50
   - Set MaxIdleConns: 10
   - Set ConnMaxLifetime: 5 menit

3. Cache template di Redis:
   - TTL 5 menit cukup untuk template yang jarang berubah

4. Timeout eksplisit di semua HTTP client:
   - notification client: 10s
   - person client: 5s
   - email client: 5s
   - firebase client: 3s

5. Batching Firebase:
   - Gunakan FCM Multicast API untuk push ke banyak device sekaligus
```

---

## Dependency eksternal

| Dependency | Tujuan | Fallback |
|---|---|---|
| PostgreSQL | Primary store semua data | Tidak ada; service tidak bisa berfungsi tanpa DB |
| Kafka | Event bus antar proses | Kalau Kafka down, create tetap jalan (publish error di-log, tidak gagalkan request) |
| Redis | Cache (belum aktif di hot path) | Tidak ada; optional |
| Person service | Resolve email/token/phone user | Consumer return error → Kafka retry |
| Email provider | Kirim email | Retry via Kafka |
| Firebase FCM | Kirim push notification | Retry via Kafka |

---

## Struktur direktori

```
cmd/
  app/          → HTTP server entrypoint
  consumer/     → Kafka consumer entrypoint (-consumer flag)
internal/
  entity/       → Domain struct (tidak boleh ada framework/DB tag)
  repository/   → Interface + model/ + postgres/ per aggregate
  usecase/      → Business logic + state machine + event publisher
  transport/
    apis/       → Echo handler, router, DTO
    event/kafka → Kafka consumer handler + router
  infrastructure/
    broker/kafka → Producer + consumer runner
    cache/redis  → Redis client
    database/    → PostgreSQL connection + migrate
  client/       → HTTP client ke service lain (person, notification, email, firebase)
  bootstrap/    → Wiring semua dependency (services.go, consumer.go, db.go)
```

**Aturan dependency:**
- `transport` boleh impor `usecase` dan `entity`
- `usecase` boleh impor `repository` (interface) dan `entity`
- `entity` tidak boleh impor layer lain
- `client/` adalah adapter eksternal, boleh diimpor oleh `usecase` dan `transport/event`
- `bootstrap/` boleh impor semua — ini satu-satunya tempat wiring konkret

---

## Known issues & backlog

| Issue | Severity | Status |
|---|---|---|
| Bug: push device token tidak di-assign di loop | P0 | Open |
| Bug: push send selalu FAILED karena messageID masuk `remark` | P0 | Open |
| `SentSender: nil` di `bootstrap/consumer.go` | P0 | Open |
| State log di-set SENT saat render, bukan saat kirim | P1 | Open |
| `idx_notifications_send_time` → kolom tidak ada (harusnya `schedule_at`) | P1 | Open |
| `idx_notification_logs_channel` → kolom di-comment | P1 | Open |
| Tidak ada autentikasi/otorisasi di HTTP routes | P1 | Open |
| Tidak ada DLQ untuk message yang gagal berulang | P2 | Open |
| Template tidak di-cache (fetch setiap message) | P2 | Open |
| Tidak ada timeout eksplisit di HTTP client | P2 | Open |
| Rate limiter belum ada di API | P2 | Open |
