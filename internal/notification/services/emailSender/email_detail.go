package email_sender

import (
	"errors"
	vo "platform/pkg/domain/value_object"
)

var (
	ErrSubjectRequired = errors.New("subject is required")
	ErrBodyRequired    = errors.New("body is required")
)

type EmailDetail struct {
	subject            string
	body               string
	from               vo.Email
	to                 vo.Email
	replyTo            *vo.Email
	cc                 []string
	bcc                []string
	attachmentFilePath *string
	attachmentFileName *string
	attachedDownloadId *int
	headers            map[string]string
}

func BaseEmailDetail(subject, body string, from, to vo.Email) (*EmailDetail, error) {
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

func (ed *EmailDetail) WithReplyTo(replyTo *vo.Email) *EmailDetail {
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
