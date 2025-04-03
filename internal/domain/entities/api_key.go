package entities

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"os"
)

type ApiKey interface {
	GetUUID() string
	GetSecret() string
	SetUUID(uuid string)
	SetSecret(secret string)
	EncryptSecret() error
	DecryptSecret() error
	GetSlug() string
	SetSlug(slug string)
}

type apiKey struct {
	uuid   string
	secret string
	slug   string
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

func (r *apiKey) EncryptSecret() error {
	key := os.Getenv("AES_KEY")

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return err
	}

	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	ciphertext := aesGCM.Seal(nonce, nonce, []byte(r.secret), nil)

	r.secret = base64.URLEncoding.EncodeToString(ciphertext)

	return nil
}

func (r *apiKey) DecryptSecret() error {
	key := os.Getenv("AES_KEY")

	data, err := base64.URLEncoding.DecodeString(r.secret)
	if err != nil {
		return err
	}

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	nonceSize := aesGCM.NonceSize()
	if len(data) < nonceSize {
		return errors.New("ciphertext is too short")
	}

	ciphertext := data[nonceSize:]
	nonce := data[:nonceSize]

	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return err
	}

	r.secret = string(plaintext)

	return nil
}
