package request

type CreateLink struct {
	URL string `json:"url" binding:"required,url"`
}
