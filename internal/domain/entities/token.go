package entities

import "time"

type Token interface {
	GetUUID() string
	GetApiKey() string
	GetCreatedAt() time.Time
	SetUUID(uuid string)
	SetApiKey(token string)
	SetCreatedAt(createdAt time.Time)
}

type token struct {
	uuid      string
	apiKey    string
	createdAt time.Time
}

func NewToken(uuid string, apiKey string) *token {
	return &token{
		uuid:      uuid,
		apiKey:    apiKey,
		createdAt: time.Now(),
	}
}

func (r *token) GetUUID() string {
	return r.uuid
}

func (r *token) GetApiKey() string {
	return r.apiKey
}

func (r *token) SetUUID(uuid string) {
	r.uuid = uuid
}

func (r *token) SetApiKey(apiKey string) {
	r.apiKey = apiKey
}

func (r *token) GetCreatedAt() time.Time {
	return r.createdAt
}

func (r *token) SetCreatedAt(createdAt time.Time) {
	r.createdAt = createdAt
}
