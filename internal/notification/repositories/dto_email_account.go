package repositories

import (
	"platform/internal/notification/domain"
	voInternal "platform/internal/notification/domain/value_object"
	voExternal "platform/pkg/domain/value_object"
	"time"

	"github.com/google/uuid"
)

// EmailAccountDTO maps database rows to domain objects and back.
type EmailAccountDTO struct {
	ID           uuid.UUID  `db:"id"`
	ProjectID    uuid.UUID  `db:"project_id"`
	Email        string     `db:"email"`
	DisplayName  string     `db:"display_name"`
	Host         string     `db:"host"`
	Port         int        `db:"port"`
	EnableSsl    bool       `db:"enable_ssl"`
	TypeID       int        `db:"type_id"`
	Username     *string    `db:"username"`
	Password     *string    `db:"password"`
	ClientID     *string    `db:"client_id"`
	TenantID     *string    `db:"tenant_id"`
	ClientSecret *string    `db:"client_secret"`
	AccessToken  *string    `db:"access_token"`
	RefreshToken *string    `db:"refresh_token"`
	ExpireAt     *time.Time `db:"expire_at"`
	CreatedAt    time.Time  `db:"created_at"`
}

// ToDomain converts the DTO into a domain EmailAccount.
func (dto *EmailAccountDTO) ToDomain() *domain.EmailAccount {
	entity := &domain.EmailAccount{}

	// Email is coming from database, we are sure it is valid, so ignore error
	email, _ := voExternal.NewEmail(dto.Email)

	entity.SetID(dto.ID)
	entity.SetProjectID(dto.ProjectID)
	entity.SetEmail(email)
	entity.SetDisplayName(dto.DisplayName)
	entity.SetHost(dto.Host)
	entity.SetPort(dto.Port)
	entity.SetEnableSSL(dto.EnableSsl)
	entity.SetSmtpType(dto.TypeID)
	entity.SetCreatedAt(dto.CreatedAt)

	if dto.TypeID == domain.Login {
		// Username and password can be null in database, so we should check it and if they are null set to empty string
		username := ptrToString(dto.Username)
		password := ptrToString(dto.Password)
		credentials := voInternal.NewTraditionalCredentials(username, password)
		entity.SetTraditionalCredentials(credentials)

		return entity
	}

	// ClientID, tenantID and ClientSecret can be null in database, so we should check it and if they are null set to empty string
	clientID := ptrToString(dto.ClientID)
	tenantID := ptrToString(dto.TenantID)
	clientSecret := ptrToString(dto.ClientSecret)
	credentials := voInternal.NewOAuth2Credentials(clientID, tenantID, clientSecret)
	entity.SetOAuth2Credentials(credentials)

	// AccessToken, RefreshToken and ExpireAt can be null in database, so we should check it and if they are null set to empty string
	accessToken := ptrToString(dto.AccessToken)
	refreshToken := ptrToString(dto.RefreshToken)
	expireAt := ptrToTime(dto.ExpireAt)
	tokenInfo := voInternal.NewTokenInformation(accessToken, refreshToken, expireAt)
	entity.SetTokenInformation(tokenInfo)

	return entity
}

// Convert from entity to database row
func (dto *EmailAccountDTO) ToDTO(ea *domain.EmailAccount) *EmailAccountDTO {
	dto.ID = ea.GetID()
	dto.ProjectID = ea.GetProjectID()
	dto.Email = ea.GetEmail().Value()
	dto.DisplayName = ea.GetDisplayName()
	dto.Host = ea.GetHost()
	dto.Port = ea.GetPort()
	dto.EnableSsl = ea.GetEnableSSL()
	dto.TypeID = ea.GetSmtpType()
	dto.CreatedAt = ea.GetCreatedAt()

	traditionalCredentials := ea.GetTraditionalCredentials()
	if traditionalCredentials != nil {
		username, password := traditionalCredentials.Credentials()
		dto.Username = ptrToStringValue(username)
		dto.Password = ptrToStringValue(password)
		return dto
	}

	oauth2Credentials := ea.GetOAuth2Credentials()
	if oauth2Credentials != nil {
		clientID, tenantID, clientSecret := oauth2Credentials.Credentials()
		dto.ClientID = ptrToStringValue(clientID)
		dto.TenantID = ptrToStringValue(tenantID)
		dto.ClientSecret = ptrToStringValue(clientSecret)
	}

	tokenInfo := ea.GetTokenInformation()
	if tokenInfo != nil {
		accessToken, refreshToken, expireAt := tokenInfo.TokenInformation()
		dto.AccessToken = ptrToStringValue(accessToken)
		dto.RefreshToken = ptrToStringValue(refreshToken)
		dto.ExpireAt = ptrToTimeValue(expireAt)
	}

	return dto
}

// ToValues returns a flat slice of fields in order for inserts/updates.
func (dto *EmailAccountDTO) GetValues() []any {
	return []any{
		dto.ID,
		dto.ProjectID,
		dto.Email,
		dto.DisplayName,
		dto.Host,
		dto.Port,
		dto.EnableSsl,
		dto.TypeID,
		dto.Username,
		dto.Password,
		dto.ClientID,
		dto.TenantID,
		dto.ClientSecret,
		dto.AccessToken,
		dto.RefreshToken,
		dto.ExpireAt,
		dto.CreatedAt,
	}
}
