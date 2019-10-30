package main

import (
	"encoding/base64"
	"errors"
	"net/smtp"
)

type loginAuth struct {
	username, password string
}

func LoginAuth(username, password string) smtp.Auth {
	return &loginAuth{username, password}
}

func (a *loginAuth) Start(server *smtp.ServerInfo) (proto string, toServer []byte, err error) {
	return "LOGIN", []byte(a.username), nil
}

func (a *loginAuth) Next(fromServer []byte, more bool) (toServer []byte, err error) {
	if more {
		switch string(fromServer) {
		case "Username:":
			return []byte(a.username), nil
		case "Password:":
			return []byte(a.password), nil
		default:
			return nil, errors.New("Unkown fromServer")
		}
	}
	return nil, nil
}

func sendMail(to []string, subject string, body string) error {
	auth := LoginAuth(emailAuthUser, emailAuthPass)

	message := "MIME-Version: 1.0" + "\r\n" +
		"Content-Type: text/plain; charset=\"utf-8\"" + "\r\n" +
		"Content-Transfer-Encoding: base64" + "\r\n" +
		"To: " + to[0] + "\r\n" +
		"Subject: " + subject + "\r\n"

	message += "\r\n" + base64.StdEncoding.EncodeToString([]byte(body))

	return smtp.SendMail(emailHostPort, auth, emailFrom, to, []byte(message))
}
