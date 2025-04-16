package domain

import (
	"time"

	"github.com/google/uuid"
)

const (
	MAX_SENT_TRIES = 5
)

type QueuedEmail struct {
	emailAccountId uuid.UUID
	to             string
	replyTo        string
	cc             string
	bcc            string
	subject        string
	body           string
	sentAt         time.Time
	sentTries      int
}

func NewQueuedEmail(emailAccountId uuid.UUID, to, replyTo, cc, bcc, subject, body string, sentAt time.Time, sentTries int) *QueuedEmail {
	return &QueuedEmail{
		emailAccountId: emailAccountId,
		to:             to,
		replyTo:        replyTo,
		cc:             cc,
		bcc:            bcc,
		subject:        subject,
		body:           body,
		sentAt:         sentAt,
		sentTries:      sentTries,
	}
}
