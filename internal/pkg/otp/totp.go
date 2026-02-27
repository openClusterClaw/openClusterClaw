package otp

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"github.com/skip2/go-qrcode"
)

const (
	// SecretSize is the size of the generated secret in bytes
	SecretSize = 32
	// CodeLength is the length of the TOTP code
	CodeLength = 6
	// BackupCodeCount is the number of backup codes generated
	BackupCodeCount = 10
	// TempTokenExpiration is the expiration time for temporary login tokens
	TempTokenExpiration = 5 * time.Minute
)

// Service handles TOTP operations
type Service struct {
	encryptionKey []byte
	issuer       string
}

// NewService creates a new TOTP service
func NewService(encryptionKey string, issuer string) *Service {
	return &Service{
		encryptionKey: []byte(encryptionKey),
		issuer:       issuer,
	}
}

// GenerateSecret generates a new TOTP secret
func (s *Service) GenerateSecret(username string) (*otp.Key, error) {
	return totp.Generate(totp.GenerateOpts{
		Issuer:      s.issuer,
		AccountName:  username,
		SecretSize:   SecretSize,
		Digits:      CodeLength,
	})
}

// GenerateBackupCodes generates backup codes for account recovery
func (s *Service) GenerateBackupCodes(count int) ([]string, error) {
	codes := make([]string, count)
	for i := 0; i < count; i++ {
		b := make([]byte, 4) // 4 bytes = 8 hex chars
		if _, err := rand.Read(b); err != nil {
			return nil, fmt.Errorf("failed to generate backup code: %w", err)
		}
		// Convert to uppercase hex for better readability
		codes[i] = strings.ToUpper(fmt.Sprintf("%08X", b))
	}
	return codes, nil
}

// ValidateCode validates a TOTP code against a secret
func (s *Service) ValidateCode(code, secret string) (bool, error) {
	// Validate the code against the secret
	valid := totp.Validate(code, secret)
	if !valid {
		return false, nil
	}
	return true, nil
}

// ValidateCodeWithWindow validates a TOTP code with a time window
// Allow codes from previous and next time window (30 seconds)
func (s *Service) ValidateCodeWithWindow(code, secret string) (bool, error) {
	valid, err := totp.ValidateCustom(
		code,
		secret,
		time.Now(),
		totp.ValidateOpts{
			Digits:    CodeLength,
			Algorithm: otp.AlgorithmSHA512,
			Skew:      1, // Allow 1 time window before and after
		},
	)
	if err != nil {
		return false, fmt.Errorf("failed to validate code: %w", err)
	}
	return valid, nil
}

// EncryptSecret encrypts a secret using AES-256-GCM
func (s *Service) EncryptSecret(secret string) (string, error) {
	if len(s.encryptionKey) != 32 {
		return "", fmt.Errorf("encryption key must be 32 bytes")
	}

	block, err := aes.NewCipher(s.encryptionKey)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return "", fmt.Errorf("failed to generate nonce: %w", err)
	}

	// Encrypt the secret
	ciphertext := gcm.Seal(nonce, nonce, []byte(secret), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// DecryptSecret decrypts an encrypted secret
func (s *Service) DecryptSecret(encrypted string) (string, error) {
	if len(s.encryptionKey) != 32 {
		return "", fmt.Errorf("encryption key must be 32 bytes")
	}

	ciphertext, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return "", fmt.Errorf("failed to decode encrypted secret: %w", err)
	}

	block, err := aes.NewCipher(s.encryptionKey)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt secret: %w", err)
	}

	return string(plaintext), nil
}

// GenerateQRCode generates a QR code image as base64 string
func (s *Service) GenerateQRCode(key *otp.Key) (string, error) {
	// Generate QR code with 256x256 pixels
	qrCode, err := qrcode.Encode(key.String(), qrcode.Medium, 256)
	if err != nil {
		return "", fmt.Errorf("failed to generate QR code: %w", err)
	}

	// Convert to base64
	base64Img := base64.StdEncoding.EncodeToString(qrCode)
	return fmt.Sprintf("data:image/png;base64,%s", base64Img), nil
}

// GenerateQRCodeFromSecret generates a QR code directly from a secret string
func (s *Service) GenerateQRCodeFromSecret(secret, username string) (string, error) {
	// Create OTP URI manually
	otpURL := fmt.Sprintf("otpauth://totp/%s:%s?secret=%s&issuer=%s&digits=%d",
		s.issuer, username, secret, s.issuer, CodeLength)

	// Generate QR code from URL
	qrCode, err := qrcode.Encode(otpURL, qrcode.Medium, 256)
	if err != nil {
		return "", fmt.Errorf("failed to generate QR code: %w", err)
	}

	base64Img := base64.StdEncoding.EncodeToString(qrCode)
	return fmt.Sprintf("data:image/png;base64,%s", base64Img), nil
}
