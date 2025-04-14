package http

import "net/http"

type responseRecorder struct {
	Headers    map[string]string
	StatusCode int
	Body       string
}

func NewResponseRecorder() *responseRecorder {
	return &responseRecorder{
		Headers: make(map[string]string),
	}
}

func (rr *responseRecorder) Header() http.Header {
	h := http.Header{}
	for k, v := range rr.Headers {
		h.Add(k, v)
	}
	return h
}

func (rr *responseRecorder) WriteHeader(statusCode int) {
	rr.StatusCode = statusCode
}

func (rr *responseRecorder) Write(body []byte) (int, error) {
	rr.Body = string(body)
	return len(body), nil
}
