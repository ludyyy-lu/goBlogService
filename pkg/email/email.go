package email

import (
	"crypto/tls"

	"gopkg.in/gomail.v2"
)

type Email struct {
	*SMTPInfo
}
type SMTPInfo struct {
	Host     string
	Port     int
	IsSSL    bool
	UserName string
	Password string
	From     string
}

func NewEmail(info *SMTPInfo) *Email {
	return &Email{SMTPInfo: info}
}

func (e *Email) SendMail(to []string, subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", e.From)     //发件人
	m.SetHeader("To", to...)        //收件人
	m.SetHeader("Subject", subject) //邮件主题
	m.SetBody("text/html", body)    //邮件正文

	dialer := gomail.NewDialer(e.Host, e.Port, e.UserName, e.Password)
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: e.IsSSL}
	return dialer.DialAndSend(m)
}
