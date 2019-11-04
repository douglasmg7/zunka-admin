package main

import (
	"encoding/base64"
	"errors"
	"log"
	"net/smtp"
)

type loginAuth struct {
	username, password string
}

func LoginAuth(username, password string) smtp.Auth {
	return &loginAuth{username, password}
}

func (a *loginAuth) Start(server *smtp.ServerInfo) (proto string, toServer []byte, err error) {
	// log.Println("Start smtp.ServerInfo:", server, "Response:", a.username)
	log.Println("Start smtp.ServerInfo:", server)
	// return "LOGIN", []byte(a.username), nil
	return "LOGIN", []byte{}, nil
}

func (a *loginAuth) Next(fromServer []byte, more bool) (toServer []byte, err error) {
	if more {
		switch string(fromServer) {
		case "Username:":
			log.Println("Next Username:", a.username)
			return []byte(a.username), nil
		case "Password:":
			log.Println("Next Password:", a.password)
			return []byte(a.password), nil
		default:
			log.Println("Next default:", string(fromServer))
			return nil, errors.New("Unkown fromServer")
		}
	}
	return nil, nil
}

func sendMail(to []string, subject string, body string) error {
	auth := LoginAuth(emailAuthUser, emailAuthPass)
	// auth := LoginAuth("zunka", emailAuthPass)

	message := "MIME-Version: 1.0" + "\r\n" +
		"Content-Type: text/plain; charset=\"utf-8\"" + "\r\n" +
		"Content-Transfer-Encoding: base64" + "\r\n" +
		"To: " + to[0] + "\r\n" +
		"Subject: " + subject + "\r\n"

	log.Printf("auth: %s\n", auth)
	log.Printf("SendMail to: %s, from: %s, host: %s, message: %sbody: %s\n", to, emailFrom, emailHostPort, message, body)

	message += "\r\n" + base64.StdEncoding.EncodeToString([]byte(body))

	return smtp.SendMail(emailHostPort, auth, emailFrom, to, []byte(message))
}
