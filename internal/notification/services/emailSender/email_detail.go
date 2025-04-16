package email_sender

import (
	"errors"
	"platform/pkg/domain/valueobject"
)

var (
	ErrSubjectRequired = errors.New("subject is required")
	ErrBodyRequired    = errors.New("body is required")
)

type EmailDetail struct {
	subject            string
	body               string
	from               valueobject.Email
	to                 valueobject.Email
	replyTo            *valueobject.Email
	cc                 []string
	bcc                []string
	attachmentFilePath *string
	attachmentFileName *string
	attachedDownloadId *int
	headers            map[string]string
}

func BaseEmailDetail(subject, body string, from, to valueobject.Email) (*EmailDetail, error) {
	if subject == "" {
		return nil, ErrSubjectRequired
	}
	if body == "" {
		return nil, ErrBodyRequired
	}

	return &EmailDetail{
		subject: subject,
		body:    body,
		from:    from,
		to:      to,
		cc:      []string{},
		bcc:     []string{},
		headers: make(map[string]string),
	}, nil
}

func (ed *EmailDetail) WithReplyTo(replyTo *valueobject.Email) *EmailDetail {
	ed.replyTo = replyTo
	return ed
}

func (ed *EmailDetail) WithCc(cc []string) *EmailDetail {
	ed.cc = cc
	return ed
}

func (ed *EmailDetail) WithBcc(bcc []string) *EmailDetail {
	ed.bcc = bcc
	return ed
}

func (ed *EmailDetail) WithAttachmentFilePath(path *string) *EmailDetail {
	ed.attachmentFilePath = path
	return ed
}

func (ed *EmailDetail) WithAttachmentFileName(name *string) *EmailDetail {
	ed.attachmentFileName = name
	return ed
}

func (ed *EmailDetail) WithAttachmentDownloadID(id *int) *EmailDetail {
	ed.attachedDownloadId = id
	return ed
}

func (ed *EmailDetail) WithHeaders(headers map[string]string) *EmailDetail {
	ed.headers = headers
	return ed
}
