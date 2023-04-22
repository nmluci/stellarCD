package dto

type WebhoookRequest struct {
	JobID   string
	Webhook map[string]interface{}
}
