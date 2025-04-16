package request

type CreateApiKeyRequest struct {
	Slug string `json:"slug" binding:"required"`
}
