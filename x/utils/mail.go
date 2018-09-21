package utils

import (
	"bytes"
	"github.com/golang/glog"
	"gopkg.in/gomail.v2"
	"html/template"
)

type Mail struct {
	Subject string
	Body    string
	To      string
}

var mailDialer = gomail.NewDialer("smtp.gmail.com", 465, "trunglenlvn@gmail.com", "gfgsimshbzgwrxwa")

func (mail Mail) Send() {
	m := gomail.NewMessage()
	m.SetHeader("From", "trunglenlvn@gmail.com")
	m.SetHeader("To", mail.To)
	// m.SetAddressHeader("Cc", "dan@example.com", "Dan")
	m.SetHeader("Subject", mail.Subject)
	m.SetBody("text/html", mail.Body)
	var data = map[string]interface{}{
		"name": "trungken",
	}
	var tmpl, err = template.ParseFiles("mail-template/feedback.html")
	buf := new(bytes.Buffer)
	if err = tmpl.Execute(buf, data); err != nil {
		glog.Error(err)
	}
	m.SetBody("text/html", buf.String())
	if err := mailDialer.DialAndSend(m); err != nil {
		panic(err)
	}
}
