package domain

import (
	"github.com/google/uuid"
)

type EmailTemplateName string

const (
	USER_EMAIL_VALIDATION EmailTemplateName = "USER_EMAIL_VALIDATION"
)

type EmailTemplate struct {
	emailAccountId    uuid.UUID
	name              EmailTemplateName
	language          string
	subject           string
	body              string
	bccEmailAddresses string
	allowDirectReply  bool
}

func NewEmailTemplate(emailAccountId uuid.UUID, name EmailTemplateName, language, subject, body, bccEmailAddress string, allowDirectReply bool) *EmailTemplate {
	return &EmailTemplate{
		emailAccountId:    emailAccountId,
		name:              name,
		language:          language,
		subject:           subject,
		body:              body,
		bccEmailAddresses: bccEmailAddress,
		allowDirectReply:  allowDirectReply,
	}
}
