package utils

import (
	"log"
	"os"
)

// AuditLogger writes audit events asynchronously.
type AuditLogger struct {
	logger *log.Logger
}

func NewAuditLogger() AuditLogger {
	return AuditLogger{
		logger: log.New(os.Stdout, "[audit] ", log.LstdFlags|log.Lshortfile),
	}
}

func (l AuditLogger) Log(action string, userID int) {
	l.logger.Printf("action=%s user_id=%d", action, userID)
}

// Notifier simulates an async notifier.
type Notifier struct {
	logger *log.Logger
}

func NewNotifier() Notifier {
	return Notifier{
		logger: log.New(os.Stdout, "[notify] ", log.LstdFlags|log.Lshortfile),
	}
}

func (n Notifier) Send(userID int, event string) {
	n.logger.Printf("user_id=%d event=%s", userID, event)
}
