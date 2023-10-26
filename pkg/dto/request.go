package dto

type WebhoookRequest struct {
	JobID     string
	HeaderMap map[string][]string
	Webhook   map[string]interface{}
	RawBody   []byte
}
