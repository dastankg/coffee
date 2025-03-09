package notification

import (
	"coffee/configs"
	"coffee/pkg/req"
	"coffee/pkg/res"
	"fmt"
	"net/http"
	"net/smtp"
)

type NotificationHandler struct {
	*configs.Config
}

type NotificationHandlerDeps struct {
	*configs.Config
}

func NewNotificationHandler(router *http.ServeMux, deps NotificationHandlerDeps) {
	handler := &NotificationHandler{
		Config: deps.Config,
	}
	router.HandleFunc("POST /notification/send", handler.SendEmail())
}

// @Summary Отправка Email
// @Description Отправка Email
// @Tags Notifaction
// @Accept json
// @Produce json
// @Param request body NotificationRequest true "Данные для отправки"
// @Success 200 {object} NotificationResponse "Успешная отправка"
// @Router /notification/send [post]
func (handler *NotificationHandler) SendEmail() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[NotificationRequest](&w, r)
		if err != nil {
			return
		}

		smtpHost := handler.Config.Smtp.SmtpHost
		smtpPort := handler.Config.Smtp.SmtpPort

		from := handler.Config.Smtp.From
		password := handler.Config.Smtp.Password
		fmt.Println(smtpHost, smtpPort)
		fmt.Println(from, password)
		to := []string{
			body.Email,
		}

		subject := body.Subject
		bodyMessage := body.Body
		message := []byte("To: " + to[0] + "\r\n" +
			"Subject: " + subject + "\r\n" +
			"\r\n" +
			bodyMessage + "\r\n")

		auth := smtp.PlainAuth("", from, password, smtpHost)

		err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
		if err != nil {
			fmt.Println(err)
			return
		}
		data := NotificationResponse{
			Message: "Succes",
		}
		res.Json(w, data, http.StatusOK)
	}
}
