package EasyEmail

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"errors"
	"fmt"
	"html/template"
	"net"
	"net/http"
	"net/smtp"
	"os"
)

type Email struct {
	MaxFileSize      int
	AllowedFileTypes []string
}

// #region: Security function calls
type loginAuth struct {
	username, password string
}

func LoginAuth(username, password string) smtp.Auth {
	return &loginAuth{username, password}
}

func (a *loginAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	return "LOGIN", []byte(a.username), nil
}

func (a *loginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		switch string(fromServer) {
		case "Username:":
			return []byte(a.username), nil
		case "Password:":
			return []byte(a.password), nil
		default:
			return nil, errors.New("Unknown from server")
		}
	}
	return nil, nil
}

// #region: Send Basic Email
// This function sends only basic templated emails using the html template provided
func (email *Email) Email_Body_Only(sender string, password string, smtpadd string, smtpHost string, smtpPort string, templname string, recivers []string, subject string, bodytext string) (string, error) {

	// Receiver email address.
	to := recivers

	conn, err := net.Dial("tcp", smtpadd)
	if err != nil {
		println(err)
	}

	c, err := smtp.NewClient(conn, smtpHost)
	if err != nil {
		println(err)
	}

	tlsconfig := &tls.Config{
		ServerName: smtpHost,
	}

	if err = c.StartTLS(tlsconfig); err != nil {
		println(err)
	}

	t, _ := template.ParseFiles(fmt.Sprintf("./%s", templname))

	auth := LoginAuth(sender, password)

	if err = c.Auth(auth); err != nil {
		println(err)
	}

	var body bytes.Buffer

	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("Subject: %s \n%s\n\n", subject, mimeHeaders)))

	t.Execute(&body, struct {
		Message string
	}{
		Message: bodytext,
	})

	// Sending email.
	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, sender, to, body.Bytes())
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	return "Sent Email", nil
}

// #region: Send Email W Attach
// This function sends emails with attchments using the text sent.
func (email *Email) Email_W_Attachments(sender string, password string, smtpadd string, smtpHost string, smtpPort string, recivers []string, subject string, bodytext string, attachmentPath string) (string, error) {

	// Receiver email address.
	to := recivers

	conn, err := net.Dial("tcp", smtpadd)
	if err != nil {
		println(err)
	}

	c, err := smtp.NewClient(conn, smtpHost)
	if err != nil {
		println(err)
	}

	tlsconfig := &tls.Config{
		ServerName: smtpHost,
	}

	if err = c.StartTLS(tlsconfig); err != nil {
		println(err)
	}

	auth := LoginAuth(sender, password)

	if err = c.Auth(auth); err != nil {
		println(err)
	}

	fileBytes, err := os.ReadFile(attachmentPath)
	if err != nil {
		fmt.Println("Error reading attachment:", err)
	}

	// Encode the file to Base64
	encoded := base64.StdEncoding.EncodeToString(fileBytes)
	contentType := http.DetectContentType(fileBytes)

	// Build an email from the info provided
	message := fmt.Sprintf("From: %s\r\n", sender) +
		fmt.Sprintf("To: %s\r\n", to) +
		fmt.Sprintf("Subject: %s\r\n", subject) +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: multipart/mixed; boundary=EndOfSection1234\r\n" +
		"\r\n" +
		"--EndOfSection1234\r\n" +
		fmt.Sprintf("Content-Type: %s; name=%s\r\n", contentType, attachmentPath) +
		"Content-Transfer-Encoding: base64\r\n" +
		fmt.Sprintf("Content-Disposition: attachment; filename=%s\r\n", attachmentPath) +
		"\r\n" +
		encoded +
		"\r\n" +
		"--EndOfSection1234\r\n" +
		"Content-Type: text/plain; charset=utf-8\r\n" +
		"Content-Transfer-Encoding: quoted-printable\r\n" +
		"\r\n" +
		bodytext +
		"\r\n" +
		"--EndOfSection1234--\r\n"

	// Sending email.
	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, sender, to, []byte(message))
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	return "Sent Email", nil
}
