package entities

type Redirect struct {
	UUID string `json:"uuid" validate:"uuid"`
	URL  string `json:"url" validate:"url"`
}

func NewRedirect(uuid string, url string) *Redirect {
	return &Redirect{
		UUID: uuid,
		URL:  url,
	}
}
