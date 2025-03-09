package notification

type NotificationRequest struct {
	Email   string `json:"email"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

type NotificationResponse struct {
	Message string `json:"message"`
}
