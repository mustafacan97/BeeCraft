package repositories

import (
	"platform/internal/notification/domain"
	internalValueObject "platform/internal/notification/domain/value_object"
	"platform/pkg/domain/valueobject"
	"time"

	"github.com/google/uuid"
)

type EmailAccountDTO struct {
	id           uuid.UUID  `db:"id"`
	projectID    uuid.UUID  `db:"project_id"`
	email        string     `db:"email"`
	displayName  string     `db:"display_name"`
	host         string     `db:"host"`
	port         int        `db:"port"`
	typeID       int        `db:"type_id"`
	enableSsl    bool       `db:"enable_ssl"`
	createdAt    time.Time  `db:"created_at"`
	username     *string    `db:"username"`
	password     *string    `db:"password"`
	clientID     *string    `db:"client_id"`
	clientSecret *string    `db:"client_secret"`
	tenantID     *string    `db:"tenant_id"`
	accessToken  *string    `db:"access_token"`
	refreshToken *string    `db:"refresh_token"`
	expireAt     *time.Time `db:"token_expire_at"`
}

func (dto *EmailAccountDTO) ToDomain() (*domain.EmailAccount, error) {
	email, err := valueobject.NewEmail(dto.email)
	if err != nil {
		return nil, err
	}

	emailAccount := domain.EmailAccountFromDB(
		dto.id,
		dto.projectID,
		dto.typeID,
		email,
		dto.displayName,
		dto.host,
		dto.port,
		dto.enableSsl,
		dto.createdAt)

	switch dto.typeID {
	case domain.Login:
		dto.applyLogin(emailAccount)
	case domain.GmailOAuth2:
	case domain.MicrosoftOAuth2:
		dto.applyOAuth2(emailAccount)
		dto.applyTokenInformation(emailAccount)
	}

	return emailAccount, nil
}

func ToDTO(emailAccount *domain.EmailAccount) *EmailAccountDTO {
	emailAccountDTO := &EmailAccountDTO{
		id:          emailAccount.GetID(),
		projectID:   emailAccount.GetProjectID(),
		email:       emailAccount.GetEmail(),
		displayName: emailAccount.GetDisplayName(),
		host:        emailAccount.GetHost(),
		port:        emailAccount.GetPort(),
		typeID:      emailAccount.GetTypeID(),
		enableSsl:   emailAccount.IsSslEnabled(),
		createdAt:   emailAccount.GetCreatedAt(),
	}

	username, password := emailAccount.TraditionalCredential.GetCredentials([]byte(emailAccount.GetEmail()))
	if username != "" {
		emailAccountDTO.username = &username
	}
	if password != "" {
		emailAccountDTO.password = &password
	}

	accessToken, refreshToken, expireAt := emailAccount.TokenInformation.GetTokenInformation()
	if accessToken != "" {
		emailAccountDTO.accessToken = &accessToken
	}
	if refreshToken != "" {
		emailAccountDTO.refreshToken = &refreshToken
	}
	if !expireAt.IsZero() {
		emailAccountDTO.expireAt = &expireAt
	}

	clientID, clientSecret, tenantID := emailAccount.OAuthCredentials.GetCredentials()
	if clientID != "" {
		emailAccountDTO.clientID = &clientID
	}
	if clientSecret != "" {
		emailAccountDTO.clientSecret = &clientSecret
	}
	if tenantID != "" {
		emailAccountDTO.tenantID = &tenantID
	}

	return emailAccountDTO
}

func (dto *EmailAccountDTO) ToValues() []any {
	return []any{
		dto.id,
		dto.projectID,
		dto.email,
		dto.displayName,
		dto.host,
		dto.port,
		dto.typeID,
		dto.enableSsl,
		dto.createdAt,
		dto.username,
		dto.password,
		dto.clientID,
		dto.clientSecret,
		dto.tenantID,
		dto.accessToken,
		dto.refreshToken,
		dto.expireAt,
	}
}

func (dto *EmailAccountDTO) applyLogin(emailAccount *domain.EmailAccount) {
	username := ""
	password := ""
	if dto.username != nil {
		username = *dto.username
	}
	if dto.password != nil {
		password = *dto.password
	}
	emailAccount.TraditionalCredential = internalValueObject.NewTraditionalCredential(username, password, []byte(dto.email))
}

func (dto *EmailAccountDTO) applyOAuth2(emailAccount *domain.EmailAccount) {
	clientID := ""
	clientSecret := ""
	tenantID := ""
	if dto.clientID != nil {
		clientID = *dto.clientID
	}
	if dto.clientSecret != nil {
		clientSecret = *dto.clientSecret
	}
	if dto.tenantID != nil {
		tenantID = *dto.tenantID
	}
	emailAccount.OAuthCredentials = internalValueObject.NewOAuthCredential(clientID, clientSecret, tenantID)
}

func (dto *EmailAccountDTO) applyTokenInformation(emailAccount *domain.EmailAccount) {
	accessToken := ""
	refreshToken := ""
	expireAt := time.Time{}
	if dto.accessToken != nil {
		accessToken = *dto.accessToken
	}
	if dto.refreshToken != nil {
		refreshToken = *dto.refreshToken
	}
	if dto.expireAt != nil {
		expireAt = *dto.expireAt
	}
	emailAccount.TokenInformation = internalValueObject.NewTokenInformation(accessToken, refreshToken, expireAt)
}
