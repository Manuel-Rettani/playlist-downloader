package service

import (
	"fmt"
	gomail "gopkg.in/mail.v2"
	"log"
)

type IMailService interface {
	SendMail(fileLink string, playListName string) error
}

type MailService struct {
	senderMail    string
	dialer        *gomail.Dialer
	requesterMail string
}

func NewMailService(smtpServer string, smtpPort int, senderMail string, password string, requesterMail string) *MailService {
	dialer := gomail.NewDialer(smtpServer, smtpPort, senderMail, password)
	return &MailService{
		senderMail:    senderMail,
		dialer:        dialer,
		requesterMail: requesterMail,
	}
}

func (m MailService) SendMail(fileLink string, playListName string) error {
	log.Println("Sending mail to: ", m.requesterMail)
	message := gomail.NewMessage()

	message.SetHeader("From", m.senderMail)
	message.SetHeader("To", m.requesterMail)
	message.SetHeader("Subject", fmt.Sprintf("Download link for playlist %s", playListName))
	message.SetBody("text/plain", fmt.Sprintf("This is the download link for playlist %s\n\n%s", playListName, fileLink))

	if err := m.dialer.DialAndSend(message); err != nil {
		return err
	} else {
		log.Println("Mail sent successfully")
	}

	return nil
}
