package entities

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
	"os"
	"time"

	"github.com/BrunoPolaski/go-crud/src/configuration/rest_err"
)

type ApiKey interface {
	GetUUID() string
	GetSecret() string
	SetUUID(uuid string)
	SetSecret(secret string)
	EncryptSecret() *rest_err.RestErr
	DecryptSecret() *rest_err.RestErr
	GetSlug() string
	SetSlug(slug string)
	GetCreatedAt() time.Time
	SetCreatedAt(createdAt time.Time)
}

type apiKey struct {
	uuid      string
	secret    string
	slug      string
	createdAt time.Time
}

func NewApiKey(uuid string, secret string, slug string) *apiKey {
	return &apiKey{
		uuid:   uuid,
		secret: secret,
		slug:   slug,
	}
}

func (r *apiKey) GetUUID() string {
	return r.uuid
}

func (r *apiKey) GetSecret() string {
	return r.secret
}

func (r *apiKey) SetUUID(uuid string) {
	r.uuid = uuid
}

func (r *apiKey) SetSecret(secret string) {
	r.secret = secret
}

func (r *apiKey) GetSlug() string {
	return r.slug
}

func (r *apiKey) SetSlug(slug string) {
	r.slug = slug
}

func (r *apiKey) EncryptSecret() *rest_err.RestErr {
	key := os.Getenv("AES_KEY")

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return rest_err.NewInternalServerError(err.Error())
	}

	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return rest_err.NewInternalServerError(err.Error())
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return rest_err.NewInternalServerError(err.Error())
	}

	ciphertext := aesGCM.Seal(nonce, nonce, []byte(r.secret), nil)

	r.secret = base64.URLEncoding.EncodeToString(ciphertext)

	return nil
}

func (r *apiKey) DecryptSecret() *rest_err.RestErr {
	key := os.Getenv("AES_KEY")

	data, err := base64.URLEncoding.DecodeString(r.secret)
	if err != nil {
		return rest_err.NewInternalServerError(err.Error())
	}

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return rest_err.NewInternalServerError(err.Error())
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return rest_err.NewInternalServerError(err.Error())
	}

	nonceSize := aesGCM.NonceSize()
	if len(data) < nonceSize {
		return rest_err.NewInternalServerError("ciphertext too short")
	}

	ciphertext := data[nonceSize:]
	nonce := data[:nonceSize]

	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return rest_err.NewInternalServerError(err.Error())
	}

	r.secret = string(plaintext)

	return nil
}

func (r *apiKey) GetCreatedAt() time.Time {
	return r.createdAt
}

func (r *apiKey) SetCreatedAt(createdAt time.Time) {
	r.createdAt = createdAt
}
