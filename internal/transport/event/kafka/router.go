package kafka

import (
	"sort"

	notifEntity "go-boilerplate-clean/internal/entity/notifications"
	"go-boilerplate-clean/internal/transport/event/kafka/handler"
	usecasenotif "go-boilerplate-clean/internal/usecase/notifications"

	libkafka "github.com/viantonugroho11/go-lib/kafka"
)

// Key consumer (satu key = satu route). Tambah constant + entry di Handlers + config di broker.
const (
	KeyNotification = "notification"
	KeySent         = "sent"
)

// Keys daftar key untuk event notifikasi. Kalau ada tipe event lain, buat KeysOrder dsb. dan gabung di init().
var Keys = []string{KeyNotification, KeySent}

// AllKeys gabungan key semua tipe event (validasi/help). Di init() digabung dengan KeysOrder dll.
var AllKeys []string

func init() {
	AllKeys = append([]string{}, Keys...)
	// nanti: AllKeys = append(AllKeys, KeysOrder...)
}

// AvailableKeys untuk flag/help (mis. cmd/consumer).
func AvailableKeys() []string {
	out := make([]string, len(AllKeys))
	copy(out, AllKeys)
	sort.Strings(out)
	return out
}

// EventServices dependency untuk semua event handler (mirip apis.Services).
type EventServices struct {
	Notification usecasenotif.NotificationService
	SentSender   handler.NotificationSender
}

// Handlers mengembalikan map key -> handler. Consumer decode message sebagai NotificationsEventMessage (Action, After, Before).
func Handlers(svc EventServices) map[string]libkafka.EventHandler[notifEntity.NotificationsEventMessage] {
	return map[string]libkafka.EventHandler[notifEntity.NotificationsEventMessage]{
		KeyNotification: handler.NewNotificationUpdateHandler(svc.Notification),
		KeySent:         handler.NewNotificationSentHandler(svc.SentSender),
	}
}
