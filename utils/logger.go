package utils

import (
	"log"
	"os"
)

// AuditLogger writes audit events asynchronously.
// Оптимизировано: убрано Lshortfile для повышения производительности
type AuditLogger struct {
	logger *log.Logger
}

func NewAuditLogger() AuditLogger {
	return AuditLogger{
		// Убираем Lshortfile для повышения производительности
		logger: log.New(os.Stdout, "[audit] ", log.LstdFlags),
	}
}

func (l AuditLogger) Log(action string, userID int) {
	l.logger.Printf("action=%s user_id=%d", action, userID)
}

// Notifier simulates an async notifier.
// Оптимизировано: убрано Lshortfile для повышения производительности
type Notifier struct {
	logger *log.Logger
}

func NewNotifier() Notifier {
	return Notifier{
		// Убираем Lshortfile для повышения производительности
		logger: log.New(os.Stdout, "[notify] ", log.LstdFlags),
	}
}

func (n Notifier) Send(userID int, event string) {
	n.logger.Printf("user_id=%d event=%s", userID, event)
}
