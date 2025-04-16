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
	projectID             uuid.UUID
	typeID                int
	email                 valueobject.Email
	displayName           string
	host                  string
	port                  int
	enableSsl             bool
	createdAt             time.Time
	TraditionalCredential *internalValueObject.TraditionalCredential
	OAuthCredentials      *internalValueObject.OAuthCredential
	TokenInformation      *internalValueObject.TokenInformation
	emailTemplates        []EmailTemplate
	queuedEmails          []QueuedEmail
}

func NewEmailAccount(projectID uuid.UUID, typeID int, email valueobject.Email, displayName, host string, port int, enableSSL bool) *EmailAccount {
	return &EmailAccount{
		AggregateRoot: domain.NewAggregateRoot(uuid.New()),
		projectID:     projectID,
		typeID:        typeID,
		email:         email,
		displayName:   displayName,
		host:          host,
		port:          port,
		enableSsl:     enableSSL,
		createdAt:     time.Now(),
	}
}

func EmailAccountFromDB(id uuid.UUID, projectID uuid.UUID, typeID int, email valueobject.Email, displayName, host string, port int, enableSSL bool, createdAt time.Time) *EmailAccount {
	return &EmailAccount{
		AggregateRoot: domain.NewAggregateRoot(id),
		projectID:     projectID,
		typeID:        typeID,
		email:         email,
		displayName:   displayName,
		host:          host,
		port:          port,
		enableSsl:     enableSSL,
		createdAt:     createdAt,
	}
}

func (ea *EmailAccount) GetProjectID() uuid.UUID { return ea.projectID }
func (ea *EmailAccount) GetTypeID() int          { return ea.typeID }
func (ea *EmailAccount) GetEmail() string        { return ea.email.GetValue() }
func (ea *EmailAccount) GetDisplayName() string  { return ea.displayName }
func (ea *EmailAccount) GetHost() string         { return ea.host }
func (ea *EmailAccount) GetPort() int            { return ea.port }
func (ea *EmailAccount) GetAddr() string         { return fmt.Sprintf("%s:%d", ea.host, ea.port) }
func (ea *EmailAccount) GetCreatedAt() time.Time { return ea.createdAt }
func (ea *EmailAccount) IsSslEnabled() bool      { return ea.enableSsl }
