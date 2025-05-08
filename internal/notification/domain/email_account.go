package domain

import (
	voInternal "platform/internal/notification/domain/value_object"
	"platform/pkg/domain"
	voExternal "platform/pkg/domain/value_object"
	"time"

	"github.com/google/uuid"
)

const (
	_               = iota
	Login           // Authentication with username and password
	GmailOAuth2     // OAuth2 authentication with Google APIs
	MicrosoftOAuth2 // OAuth2 authentication with Microsoft Authentication
)

type EmailAccount struct {
	domain.AggregateRoot
	projectID              uuid.UUID
	typeID                 int
	email                  voExternal.Email
	displayName            string
	host                   string
	port                   int
	enableSsl              bool
	traditionalCredentials *voInternal.TraditionalCredentials
	oAuth2Credentials      *voInternal.OAuth2Credentials
	tokenInformation       *voInternal.TokenInformation
	createdAt              time.Time
	emailTemplates         []EmailTemplate
	queuedEmails           []QueuedEmail
}

func NewEmailAccount(id, projectID uuid.UUID, typeID int, email voExternal.Email, displayName string, host string, port int, enableSSL bool) *EmailAccount {
	emailAccount := &EmailAccount{
		AggregateRoot:  domain.NewAggregateRoot(id),
		projectID:      projectID,
		typeID:         typeID,
		email:          email,
		displayName:    displayName,
		host:           host,
		port:           port,
		enableSsl:      enableSSL,
		createdAt:      time.Now(),
		emailTemplates: make([]EmailTemplate, 0),
		queuedEmails:   make([]QueuedEmail, 0),
	}

	return emailAccount
}

// GETTERS
func (ea *EmailAccount) GetProjectID() uuid.UUID    { return ea.projectID }
func (ea *EmailAccount) GetEmail() voExternal.Email { return ea.email }
func (ea *EmailAccount) GetDisplayName() string     { return ea.displayName }
func (ea *EmailAccount) GetHost() string            { return ea.host }
func (ea *EmailAccount) GetPort() int               { return ea.port }
func (ea *EmailAccount) GetEnableSSL() bool         { return ea.enableSsl }
func (ea *EmailAccount) GetSmtpType() int           { return ea.typeID }
func (ea *EmailAccount) GetTraditionalCredentials() *voInternal.TraditionalCredentials {
	return ea.traditionalCredentials
}
func (ea *EmailAccount) GetOAuth2Credentials() *voInternal.OAuth2Credentials {
	return ea.oAuth2Credentials
}
func (ea *EmailAccount) GetTokenInformation() *voInternal.TokenInformation {
	return ea.tokenInformation
}
func (ea *EmailAccount) GetCreatedAt() time.Time        { return ea.createdAt }
func (ea *EmailAccount) GetTemplates() []EmailTemplate  { return ea.emailTemplates }
func (ea *EmailAccount) GetQueuedEmails() []QueuedEmail { return ea.queuedEmails }

// SETTERS
func (ea *EmailAccount) SetProjectID(id uuid.UUID)         { ea.projectID = id }
func (ea *EmailAccount) SetEmail(email voExternal.Email)   { ea.email = email }
func (ea *EmailAccount) SetDisplayName(displayName string) { ea.displayName = displayName }
func (ea *EmailAccount) SetHost(host string)               { ea.host = host }
func (ea *EmailAccount) SetPort(port int)                  { ea.port = port }
func (ea *EmailAccount) SetEnableSSL(enableSSL bool)       { ea.enableSsl = enableSSL }
func (ea *EmailAccount) SetSmtpType(smtpType int)          { ea.typeID = smtpType }
func (ea *EmailAccount) SetTraditionalCredentials(credentials *voInternal.TraditionalCredentials) {
	ea.traditionalCredentials = credentials
	ea.oAuth2Credentials = nil
	ea.tokenInformation = nil
}
func (ea *EmailAccount) SetOAuth2Credentials(credentials *voInternal.OAuth2Credentials) {
	ea.traditionalCredentials = nil
	ea.oAuth2Credentials = credentials
}
func (ea *EmailAccount) SetTokenInformation(tokenInformation *voInternal.TokenInformation) {
	ea.traditionalCredentials = nil
	ea.tokenInformation = tokenInformation
}
func (ea *EmailAccount) SetCreatedAt(createdAt time.Time) { ea.createdAt = createdAt }
