package response

type CreateApiKeyResponse struct {
	UUID   string `json:"uuid"`
	Secret string `json:"secret"`
	Slug   string `json:"slug"`
}
