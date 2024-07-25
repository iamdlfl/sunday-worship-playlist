package main

import (
	"log"
	"net/smtp"

	"gopkg.in/ini.v1"
)

func NewMailer(configFile string) (*mailer, error) {
	m := mailer{
		ConfigFileName: configFile,
		emailTo:        []string{"davidjlynch1017@gmail.com"},
		smtpHost:       "smtp.gmail.com",
		smtpPort:       "587",
	}

	cfg, err := ini.Load(configFile)
	if err != nil {
		log.Println(err)
		return &m, err
	}

	m.emailFrom = cfg.Section("email").Key("from").String()
	m.emailFromPw = cfg.Section("email").Key("pw").String()

	m.auth = smtp.PlainAuth("", m.emailFrom, m.emailFromPw, m.smtpHost)
	return &m, nil
}

type mailer struct {
	ConfigFileName string
	emailFrom      string
	emailFromPw    string
	emailTo        []string
	smtpHost       string
	smtpPort       string
	auth           smtp.Auth
}

func (m mailer) SendMessage(msg string) error {
	log.Println("Sending mail!")
	msgWithHeaders := "To: " + m.emailTo[0] + "\r\n" +
		"From: " + m.emailFrom + "\r\n" +
		"Subject: Sunday Worship Playlist Script\r\n" +
		"\r\n" +
		msg + "\r\n"
	return smtp.SendMail(m.smtpHost+":"+m.smtpPort, m.auth, m.emailFrom, m.emailTo, []byte(msgWithHeaders))
}

func (m mailer) SendMessageTo(msg, to string) error {
	log.Println("Sending mail to!")
	msgWithHeaders := "To: " + to + "\r\n" +
		"From: " + m.emailFrom + "\r\n" +
		"Subject: Sunday Worship Playlist Script\r\n" +
		"\r\n" +
		msg + "\r\n"
	return smtp.SendMail(m.smtpHost+":"+m.smtpPort, m.auth, m.emailFrom, []string{to}, []byte(msgWithHeaders))
}
