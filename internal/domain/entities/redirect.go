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

type Redirect interface {
	GetUUID() string
	GetURL() string
	SetUUID(uuid string)
	SetURL(url string)
	EncryptURL() error
	DecryptURL() error
}

type redirect struct {
	uuid string
	url  string
}

func NewRedirect(uuid string, url string) *redirect {
	return &redirect{
		uuid: uuid,
		url:  url,
	}
}

func (r *redirect) GetUUID() string {
	return r.uuid
}

func (r *redirect) GetURL() string {
	return r.url
}

func (r *redirect) SetUUID(uuid string) {
	r.uuid = uuid
}

func (r *redirect) SetURL(url string) {
	r.url = url
}

func (r *redirect) EncryptURL() error {
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

	ciphertext := aesGCM.Seal(nonce, nonce, []byte(r.url), nil)

	r.url = base64.URLEncoding.EncodeToString(ciphertext)

	return nil
}

func (r *redirect) DecryptURL() error {
	key := os.Getenv("AES_KEY")

	data, err := base64.URLEncoding.DecodeString(r.url)
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

	r.url = string(plaintext)

	return nil
}
