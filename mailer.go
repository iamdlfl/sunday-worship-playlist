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
	return smtp.SendMail(m.smtpHost+":"+m.smtpPort, m.auth, m.emailFrom, m.emailTo, []byte(msg))
}
