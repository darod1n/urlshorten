package models

type BatchRequest struct {
	CorrelationID string `json:"correlation_id"`
	OriginURL     string `json:"original_url"`
}

type BatchResponse struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

type UserURLS struct {
	OriginURL string `json:"original_url"`
	ShortURL  string `json:"short_url"`
}

type ctxKey int8

const (
	_ ctxKey = iota
	CtxKeyUserID
)
