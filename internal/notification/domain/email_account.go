package domain

import (
	"fmt"
	internalValueObject "platform/internal/notification/domain/value_object"
	"platform/pkg/domain"
	"platform/pkg/domain/valueobject"
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
	email                  valueobject.Email
	displayName            string
	host                   string
	port                   int
	enableSsl              bool
	createdAt              time.Time
	TraditionalCredentials *internalValueObject.TraditionalCredential
	OAuth2Credentials      *internalValueObject.OAuth2Credential
	TokenInformation       *internalValueObject.TokenInformation
	emailTemplates         []EmailTemplate
	queuedEmails           []QueuedEmail
}

func NewEmailAccount(id uuid.UUID, projectID uuid.UUID, typeID int, email valueobject.Email, displayName, host string, port int, enableSSL bool) *EmailAccount {
	emailAccount := &EmailAccount{
		AggregateRoot: domain.NewAggregateRoot(id),
		projectID:     projectID,
		typeID:        typeID,
		email:         email,
		displayName:   displayName,
		host:          host,
		port:          port,
		enableSsl:     enableSSL,
		createdAt:     time.Now(),
	}

	return emailAccount
}

// GETTERS
func (ea *EmailAccount) GetProjectID() uuid.UUID        { return ea.projectID }
func (ea *EmailAccount) GetEmail() valueobject.Email    { return ea.email }
func (ea *EmailAccount) GetSmtpType() int               { return ea.typeID }
func (ea *EmailAccount) GetDisplayName() string         { return ea.displayName }
func (ea *EmailAccount) GetHost() string                { return ea.host }
func (ea *EmailAccount) GetPort() int                   { return ea.port }
func (ea *EmailAccount) GetAddress() string             { return fmt.Sprintf("%s:%d", ea.host, ea.port) }
func (ea *EmailAccount) GetEnableSSL() bool             { return ea.enableSsl }
func (ea *EmailAccount) IsSslEnabled() bool             { return ea.enableSsl }
func (ea *EmailAccount) GetCreatedDate() time.Time      { return ea.createdAt }
func (ea *EmailAccount) GetTemplates() []EmailTemplate  { return ea.emailTemplates }
func (ea *EmailAccount) GetQueuedEmails() []QueuedEmail { return ea.queuedEmails }

// SETTERS
func (ea *EmailAccount) SetProjectID(id uuid.UUID)         { ea.projectID = id }
func (ea *EmailAccount) SetEmail(email valueobject.Email)  { ea.email = email }
func (ea *EmailAccount) SetDisplayName(displayName string) { ea.displayName = displayName }
func (ea *EmailAccount) SetHost(host string)               { ea.host = host }
func (ea *EmailAccount) SetPort(port int)                  { ea.port = port }
func (ea *EmailAccount) SetEnableSSL(enableSSL bool)       { ea.enableSsl = enableSSL }
func (ea *EmailAccount) SetSMTPType(smtpType int)          { ea.typeID = smtpType }
func (ea *EmailAccount) SetCreatedAt(createdAt time.Time)  { ea.createdAt = createdAt }
func (ea *EmailAccount) SetTraditionalCredentials(credentials *internalValueObject.TraditionalCredential) {
	ea.TraditionalCredentials = credentials
	ea.OAuth2Credentials = nil
	ea.TokenInformation = nil
}
func (ea *EmailAccount) SetOAuth2Credentials(credentials *internalValueObject.OAuth2Credential) {
	ea.TraditionalCredentials = nil
	ea.OAuth2Credentials = credentials
}
func (ea *EmailAccount) SetTokenInformation(tokenInformation *internalValueObject.TokenInformation) {
	ea.TraditionalCredentials = nil
	ea.TokenInformation = tokenInformation
}
