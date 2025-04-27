package repositories

import (
	"platform/internal/notification/domain"
	internalValueObjects "platform/internal/notification/domain/value_object"
	"platform/pkg/domain/valueobject"
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
	email, _ := valueobject.NewEmail(dto.Email)
	emailAccount := domain.NewEmailAccount(dto.ID, dto.ProjectID, dto.TypeID, email, dto.DisplayName, dto.Host, dto.Port, dto.EnableSsl)
	emailAccount.SetCreatedAt(dto.CreatedAt)

	dto.applyTraditionalCredentials(emailAccount)
	dto.applyOAuth2Credentials(emailAccount)
	dto.applyTokenInformation(emailAccount)

	return emailAccount
}

// Convert from entity to database row
func ToDTO(emailAccount *domain.EmailAccount) *EmailAccountDTO {
	dto := &EmailAccountDTO{
		ID:          emailAccount.GetID(),
		ProjectID:   emailAccount.GetProjectID(),
		Email:       emailAccount.Email.GetValue(),
		DisplayName: emailAccount.GetDisplayName(),
		Host:        emailAccount.GetHost(),
		Port:        emailAccount.GetPort(),
		TypeID:      emailAccount.GetSmtpType(),
		EnableSsl:   emailAccount.IsSslEnabled(),
		CreatedAt:   emailAccount.GetCreatedDate(),
	}

	if emailAccount.TraditionalCredentials != nil {
		username, password := emailAccount.TraditionalCredentials.GetCredentials()
		dto.Username = ptrToStringValue(username)
		dto.Password = ptrToStringValue(password)
	}

	if emailAccount.OAuth2Credentials != nil {
		clientID, tenantID, clientSecret := emailAccount.OAuth2Credentials.GetCredentials()
		dto.ClientID = ptrToStringValue(clientID)
		dto.TenantID = ptrToStringValue(tenantID)
		dto.ClientSecret = ptrToStringValue(clientSecret)
	}

	if emailAccount.TokenInformation != nil {
		accessToken, refreshToken, expireAt := emailAccount.TokenInformation.GetTokenInformation()
		dto.AccessToken = ptrToStringValue(accessToken)
		dto.RefreshToken = ptrToStringValue(refreshToken)
		dto.ExpireAt = ptrToTimeValue(expireAt)
	}

	return dto
}

// ToValues returns a flat slice of fields in order for inserts/updates.
func (dto *EmailAccountDTO) ToValues() []any {
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

func (dto *EmailAccountDTO) applyTraditionalCredentials(ea *domain.EmailAccount) {
	if ea.GetSmtpType() != domain.Login {
		return
	}

	username := ptrToString(dto.Username)
	password := ptrToString(dto.Password)
	ea.TraditionalCredentials.SetTraditionalCredentials(username, password)
}

func (dto *EmailAccountDTO) applyOAuth2Credentials(ea *domain.EmailAccount) {
	if ea.GetSmtpType() == domain.Login {
		return
	}

	clientID := ptrToString(dto.ClientID)
	tenantID := ptrToString(dto.TenantID)
	clientSecret := ptrToString(dto.ClientSecret)
	ea.OAuth2Credentials = internalValueObjects.NewOAuth2Credentials(clientID, tenantID, clientSecret)
}

func (dto *EmailAccountDTO) applyTokenInformation(ea *domain.EmailAccount) {
	if ea.GetSmtpType() == domain.Login {
		return
	}

	accessToken := ptrToString(dto.AccessToken)
	refreshToken := ptrToString(dto.RefreshToken)
	expireAt := ptrToTime(dto.ExpireAt)
	ea.TokenInformation = internalValueObjects.NewTokenInformation(accessToken, refreshToken, expireAt)
}

// ptrToString creates a pointer to s if non-nil; otherwise returns nil.
func ptrToString(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}

// ptrToStringValue returns a *string if the value is non-empty.
func ptrToStringValue(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// ptrToTime returns time value or zero.
func ptrToTime(t *time.Time) time.Time {
	if t != nil {
		return *t
	}
	return time.Time{}
}

// ptrToTimeValue returns *time.Time if non-zero.
func ptrToTimeValue(t time.Time) *time.Time {
	if !t.IsZero() {
		return &t
	}
	return nil
}
