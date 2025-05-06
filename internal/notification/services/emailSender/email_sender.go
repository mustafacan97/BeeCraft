// Package email_sender provides functions to send MIME-formatted emails via SMTP
// using various authentication methods including classic login and OAuth2 (Gmail and Microsoft).
package email_sender

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/base64"
	"errors"
	"fmt"
	"mime"
	"mime/multipart"
	"net"
	"net/smtp"
	"net/textproto"
	"os"
	"path/filepath"
	"platform/internal/notification/domain"
	voInternal "platform/internal/notification/domain/value_object"
	"platform/internal/notification/services/encryption"
	"strings"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func SendEmail(encryption encryption.EncryptionService, emailAccount *domain.EmailAccount, request *EmailDetail) error {
	// Build the MIME message into a buffer.
	messageBuffer, err := buildMIMEMessage(request)
	if err != nil {
		return err
	}

	// Build an SMTP client with proper authentication.
	smtpClient, err := buildSmtpClient(encryption, emailAccount)
	if err != nil {
		return err
	}
	defer smtpClient.Quit()

	// Set the sender, recipients (To, Cc, Bcc) and send the email.
	if err = smtpClient.Mail(request.from.Value()); err != nil {
		return err
	}

	if err = smtpClient.Rcpt(request.to.Value()); err != nil {
		return err
	}

	for _, addr := range request.cc {
		if err = smtpClient.Rcpt(strings.TrimSpace(addr)); err != nil {
			return err
		}
	}
	for _, addr := range request.bcc {
		if err = smtpClient.Rcpt(strings.TrimSpace(addr)); err != nil {
			return err
		}
	}

	// Write the data and close writer to complete sending.
	dataWriter, err := smtpClient.Data()
	if err != nil {
		return err
	}

	_, err = dataWriter.Write(messageBuffer.Bytes())
	if err != nil {
		return err
	}

	return dataWriter.Close()
}

func buildMIMEMessage(request *EmailDetail) (*bytes.Buffer, error) {
	var buf bytes.Buffer

	// Create a multipart writer for a mixed MIME message.
	mixedWriter := multipart.NewWriter(&buf)

	// Write essential headers.
	buf.WriteString("MIME-Version: 1.0\r\n")
	buf.WriteString(fmt.Sprintf("Content-Type: multipart/mixed; boundary=%s\r\n", mixedWriter.Boundary()))
	buf.WriteString("Content-Transfer-Encoding: base64\r\n")
	buf.WriteString(fmt.Sprintf("Date: %s\r\n", time.Now().Format(time.RFC1123Z)))
	buf.WriteString(fmt.Sprintf("Subject: %s\r\n", request.subject))
	buf.WriteString(fmt.Sprintf("From: %s\r\n", request.from.Value()))
	buf.WriteString(fmt.Sprintf("To: %s\r\n", request.to.Value()))
	if len(request.cc) > 0 {
		buf.WriteString(fmt.Sprintf("Cc: %s\r\n", strings.Join(request.cc, ", ")))
	}
	if request.replyTo != nil {
		buf.WriteString(fmt.Sprintf("Reply-To: %s\r\n", request.replyTo.Value()))
	}
	for k, v := range request.headers {
		buf.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	buf.WriteString("\r\n")

	// Add HTML text part.
	if err := addHTMLPart(mixedWriter, request.body); err != nil {
		return nil, err
	}

	// Add attachment if provided.
	if request.attachmentFilePath != nil && fileExists(*request.attachmentFilePath) {
		if err := addAttachment(mixedWriter, *request.attachmentFilePath, *request.attachmentFileName); err != nil {
			return nil, err
		}
	}

	// Close the multipart writer to flush the boundary.
	if err := mixedWriter.Close(); err != nil {
		return nil, err
	}

	return &buf, nil
}

func addHTMLPart(mixedWriter *multipart.Writer, htmlContent string) error {
	header := make(textproto.MIMEHeader)
	header.Set("Content-Type", "text/html; charset=UTF-8")
	part, err := mixedWriter.CreatePart(header)
	if err != nil {
		return err
	}
	_, err = part.Write([]byte(htmlContent))
	return err
}

func fileExists(filePath string) bool {
	_, err := os.ReadFile(filePath)
	return err == nil
}

func addAttachment(mixedWriter *multipart.Writer, filePath, fileName string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	// Determine MIME type using file extension.
	mimeType := mime.TypeByExtension(filepath.Ext(filePath))
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	h := make(textproto.MIMEHeader)
	h.Set("Content-Type", fmt.Sprintf("%s; name=%q", mimeType, fileName))
	h.Set("Content-Transfer-Encoding", "base64")
	h.Set("Content-Disposition", fmt.Sprintf(`attachment; filename=%q`, fileName))

	part, err := mixedWriter.CreatePart(h)
	if err != nil {
		return err
	}

	// Base64 encoded and written to part (RFC compliant with a line break every 76 characters)
	encoder := base64.NewEncoder(base64.StdEncoding, NewBase64LineWriter(part))
	defer encoder.Close()

	_, err = encoder.Write(data)
	return err
}

func buildSmtpClient(encryption encryption.EncryptionService, ea *domain.EmailAccount) (*smtp.Client, error) {
	tlsConfig := &tls.Config{ServerName: ea.GetHost()}
	addr := fmt.Sprintf("%s:%d", ea.GetHost(), ea.GetPort())
	var client *smtp.Client
	var err error

	// If SSL is enabled (implicit TLS), use tls.Dial. Otherwise use net.Dial.
	if ea.GetEnableSSL() {
		conn, err := tls.Dial("tcp", addr, tlsConfig)
		if err != nil {
			return nil, err
		}
		client, err = smtp.NewClient(conn, ea.GetHost())
		if err != nil {
			return nil, err
		}
	} else {
		conn, err := net.Dial("tcp", addr)
		if err != nil {
			return nil, err
		}
		client, err = smtp.NewClient(conn, ea.GetHost())
		if err != nil {
			return nil, err
		}
	}

	// If the server supports STARTTLS, upgrade to a secure connection.
	if ok, _ := client.Extension("STARTTLS"); ok {
		if err = client.StartTLS(tlsConfig); err != nil {
			return nil, err
		}
	}

	// Authenticate based on the email account type.
	switch ea.GetSmtpType() {
	case domain.Login:
		username, password := ea.GetTraditionalCredentials().Credentials()
		rawPassword, err := encryption.Decrypt(password)
		if err != nil {
			return nil, err
		}
		auth := smtp.PlainAuth("", username, rawPassword, ea.GetHost())
		if err := client.Auth(auth); err != nil {
			return nil, err
		}
	case domain.GmailOAuth2:
	case domain.MicrosoftOAuth2:
		token, err := getOAuth2Credentials(ea)
		if err != nil {
			return nil, err
		}
		if err := client.Auth(NewOAuth2Auth(ea.GetEmail().Value(), token.AccessToken)); err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("unsupported auth method")
	}

	return client, nil
}

func getOAuth2Credentials(emailAccount *domain.EmailAccount) (*oauth2.Token, error) {
	accountType := emailAccount.GetSmtpType()
	clientID, clientSecret, tenantID := emailAccount.GetOAuth2Credentials().Credentials()
	if clientID == "" || clientSecret == "" {
		return nil, errors.New("ClientId and ClientSecret are required")
	}
	if accountType == domain.MicrosoftOAuth2 && tenantID == "" {
		return nil, errors.New("ClientID, ClientSecret and TenantID are required")
	}

	accessToken, refreshToken, expireAt := emailAccount.GetTokenInformation().TokenInformation()
	token := &oauth2.Token{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Expiry:       expireAt,
		TokenType:    "Bearer",
	}

	// If token is expired, refresh it
	if !token.Valid() {
		var conf oauth2.Config

		if accountType == domain.GmailOAuth2 {
			conf = *gmailOAuth2Configs(clientID, clientSecret)
		} else if accountType == domain.MicrosoftOAuth2 {
			conf = *exchangeOAuth2Configs(clientID, clientSecret, tenantID)
		}

		ts := conf.TokenSource(context.Background(), token)
		newToken, err := ts.Token()
		if err != nil {
			return nil, err
		}

		emailAccount.SetTokenInformation(voInternal.NewTokenInformation(newToken.AccessToken, newToken.RefreshToken, newToken.Expiry))
		return newToken, nil
	}

	return token, nil
}

func gmailOAuth2Configs(clientID, clientSecret string) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint:     google.Endpoint,
		Scopes:       []string{"https://mail.google.com/"},
		RedirectURL:  "http://localhost:8080/callback",
	}
}

func exchangeOAuth2Configs(clientID, clientSecret, tenantID string) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  "http://localhost:8080/callbackk",
		Scopes: []string{
			"https://outlook.office365.com/SMTP.Send",
			"offline_access",
		},
		Endpoint: oauth2.Endpoint{
			AuthURL:  fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/authorize", tenantID),
			TokenURL: fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/token", tenantID),
		},
	}
}
