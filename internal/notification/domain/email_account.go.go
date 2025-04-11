package domain

import (
	"errors"
	"fmt"
	"platform/pkg/domain"
	"platform/pkg/domain/valueobject"
	"time"

	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const (
	Login           = iota // Authentication with username and password
	GmailOAuth2            // OAuth2 authentication with Google APIs
	MicrosoftOAuth2        // OAuth2 authentication with Microsoft Authentication
)

type EmailAccount struct {
	domain.AggregateRoot
	accessToken  string
	clientId     string
	clientSecret string
	createdAt    time.Time
	displayName  string
	email        valueobject.Email
	enableSsl    bool
	expireAt     time.Time
	host         string
	password     string
	port         int
	projectId    uuid.UUID
	refreshToken string
	tenantId     string
	typeId       int
	username     string
	// Relationships
	emailTemplates []EmailTemplate
	queuedEmails   []QueuedEmail
}

// AccessToken		 ---> can be null
// ClientId          ---> (WHEN -> Gmail | Microsoft) required
// ClientSecret		 ---> (WHEN -> Gmail | Microsoft) required
// CreatedAt         ---> time.Now()
// DisplayName       ---> required
// Email             ---> required
// Host				 ---> can be empty but not mail sending
// ID 				 ---> required
// Password			 ---> can be null
// Port				 ---> can be empty but not mail sending
// ProjectId   		 ---> required
// RefreshToken		 ---> can be null
// Microsoft		 ---> TenantId required
// TypeId			 ---> Login | Gmail | Microsoft
// Username			 ---> (WHEN -> Login) required

func NewLoginEmailAccount(projectId uuid.UUID, email valueobject.Email, displayName, host, username, password string, port int, enableSsl bool) *EmailAccount {
	return &EmailAccount{
		AggregateRoot: domain.NewAggregateRoot(uuid.New()),
		accessToken:   "",
		clientId:      "",
		clientSecret:  "",
		createdAt:     time.Now(),
		displayName:   displayName,
		email:         email,
		enableSsl:     enableSsl,
		expireAt:      time.Time{},
		host:          host,
		password:      password,
		port:          port,
		projectId:     projectId,
		refreshToken:  "",
		tenantId:      "",
		typeId:        Login,
		username:      username,
	}
}

func NewMicrosoftEmailAccount(projectId uuid.UUID, email valueobject.Email, displayName, host, clientId, clientSecret, tenantId string, port int, enableSsl bool) *EmailAccount {
	return &EmailAccount{
		AggregateRoot: domain.NewAggregateRoot(uuid.New()),
		accessToken:   "",
		clientId:      clientId,
		clientSecret:  clientSecret,
		createdAt:     time.Now(),
		displayName:   displayName,
		email:         email,
		enableSsl:     enableSsl,
		expireAt:      time.Time{},
		host:          host,
		password:      "",
		port:          port,
		projectId:     projectId,
		refreshToken:  "",
		tenantId:      tenantId,
		typeId:        MicrosoftOAuth2,
		username:      "",
	}
}

func NewGmailEmailAccount(projectId uuid.UUID, email valueobject.Email, displayName, host, clientId, clientSecret string, port int, enableSsl bool) *EmailAccount {
	return &EmailAccount{
		AggregateRoot: domain.NewAggregateRoot(uuid.New()),
		accessToken:   "",
		clientId:      clientId,
		clientSecret:  clientSecret,
		createdAt:     time.Now(),
		displayName:   displayName,
		email:         email,
		enableSsl:     enableSsl,
		expireAt:      time.Time{},
		host:          host,
		password:      "",
		port:          port,
		projectId:     projectId,
		refreshToken:  "",
		tenantId:      "",
		typeId:        GmailOAuth2,
		username:      "",
	}
}

func (ea *EmailAccount) GetEmail() string {
	return ea.email.GetValue()
}

func (ea *EmailAccount) GetHost() string {
	return ea.host
}

func (ea *EmailAccount) GetPassword() string {
	return ea.password
}

func (ea *EmailAccount) GetAddr() string {
	return fmt.Sprintf("%s:%d", ea.host, ea.port)
}

func (ea *EmailAccount) GetTypeId() int {
	return ea.typeId
}

func (ea *EmailAccount) GetUsername() string {
	return ea.username
}

func (ea *EmailAccount) GetGmailInfo() (clientId, clientSecret string) {
	return ea.clientId, ea.clientSecret
}

func (ea *EmailAccount) GetExchangeInfo() (clientID, clientSecret, tenantID string) {
	return ea.clientId, ea.clientSecret, ea.tenantId
}

func (ea *EmailAccount) GetTokenInfo() (accessToken, refreshToken string, expireAt time.Time, err error) {
	if ea.typeId == Login {
		return "", "", time.Time{}, errors.New("access and refresh token are not usable for SMTP username-password type")
	}
	return ea.accessToken, ea.refreshToken, ea.expireAt, nil
}

func (ea *EmailAccount) SetTokenInfo(accessToken, refreshToken string, expireAt time.Time) error {
	if ea.typeId == Login {
		return errors.New("token info is not usable for SMTP username-password type")
	}

	ea.accessToken = accessToken
	ea.refreshToken = refreshToken
	ea.expireAt = expireAt
	return nil
}

func (ea *EmailAccount) GetGmailOAuthTokenUrl() (string, error) {
	if ea.typeId != GmailOAuth2 {
		return "", errors.New("only gmail email account can have oauth token url")
	}

	if ea.clientId == "" || ea.clientSecret == "" {
		return "", errors.New("clientId and clientSecret are required")
	}

	conf := &oauth2.Config{
		ClientID:     ea.clientId,
		ClientSecret: ea.clientSecret,
		Endpoint:     google.Endpoint,
		Scopes:       []string{"https://mail.google.com/"},
		RedirectURL:  "https://api.platform.com/auth-return",
	}

	// To get access we should go to link below
	// Add custom query param: emailAccountId
	authURL := conf.AuthCodeURL("state-token",
		oauth2.AccessTypeOffline,
		oauth2.SetAuthURLParam("emailAccountId", ea.ID.String()),
	)

	return authURL, nil
}

/* return-auth endpoint'i nasıl olacak:
func HandleGmailOAuthCallback(c *fiber.Ctx) error {
	code := c.Query("code")
	emailAccountId := c.Query("emailAccountId")

	if code == "" || emailAccountId == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Missing code or emailAccountId")
	}

	// EmailAccount'ı ID üzerinden bul
	emailAccount, err := emailAccountRepo.GetByID(context.Background(), emailAccountId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("EmailAccount not found")
	}

	conf := &oauth2.Config{
		ClientID:     emailAccount.ClientID,
		ClientSecret: emailAccount.ClientSecret,
		Endpoint:     google.Endpoint,
		Scopes:       []string{"https://mail.google.com/"},
		RedirectURL:  "https://api.platform.com/auth-return",
	}

	// code'u AccessToken ve RefreshToken ile takas et
	token, err := conf.Exchange(context.Background(), code)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to exchange token: " + err.Error())
	}

	// AccessToken, RefreshToken ve expiry değerlerini DB'de güncelle
	emailAccount.AccessToken = token.AccessToken
	emailAccount.RefreshToken = token.RefreshToken
	emailAccount.TokenExpiry = token.Expiry // optional: saklamıyorsan null geç

	err = emailAccountRepo.Update(context.Background(), emailAccount)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to update email account: " + err.Error())
	}

	return c.SendString("Authorization successful! You can now send email via Gmail 🎉")
}
*/
