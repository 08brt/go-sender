package main

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/smtp"
	"os"
	"path/filepath"
	"strings"
)

var (
	host       = "smtp.gmail.com"
	username   = ""
	password   = ""
	portNumber = "587"
	filePath   = ""
)

type Message struct {
	To          []string
	Subject     string
	Body        string
	Attachments map[string][]byte
}

func NewMessage(s, b string) *Message {
	return &Message{Subject: s, Body: b, Attachments: make(map[string][]byte)}
}

func (m *Message) AttachFile(src string) error {
	b, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	_, fileName := filepath.Split(src)
	m.Attachments[fileName] = b
	return nil
}

func (m *Message) ToBytes() []byte {
	buf := bytes.NewBuffer(nil)
	withAttachments := len(m.Attachments) > 0
	buf.WriteString(fmt.Sprintf("Subject: %s\n", m.Subject))
	buf.WriteString(fmt.Sprintf("To: %s\n", strings.Join(m.To, ",")))

	buf.WriteString("MIME-Version: 1.0\n")
	writer := multipart.NewWriter(buf)
	boundary := writer.Boundary()
	if withAttachments {
		buf.WriteString(fmt.Sprintf("Content-Type: multipart/mixed; boundary=%s\n", boundary))
		buf.WriteString(fmt.Sprintf("--%s\n", boundary))
	} else {
		buf.WriteString("Content-Type: text/plain; charset=utf-8\n")
	}

	buf.WriteString(m.Body)
	if withAttachments {
		for k, v := range m.Attachments {
			buf.WriteString(fmt.Sprintf("\n\n--%s\n", boundary))
			buf.WriteString(fmt.Sprintf("Content-Type: %s\n", http.DetectContentType(v)))
			buf.WriteString("Content-Transfer-Encoding: base64\n")
			buf.WriteString(fmt.Sprintf("Content-Disposition: attachment; filename=%s\n", k))

			b := make([]byte, base64.StdEncoding.EncodedLen(len(v)))
			base64.StdEncoding.Encode(b, v)
			buf.Write(b)
			buf.WriteString(fmt.Sprintf("\n--%s", boundary))
		}

		buf.WriteString("--")
	}

	return buf.Bytes()
}

func main() {

	reader := bufio.NewReader(os.Stdin)

	for {
		// Get email address
		fmt.Print("Enter email address (or 'exit'): ")
		email, _ := reader.ReadString('\n')
		email = strings.TrimSpace(email)

		// Check if the user wants to exit
		if strings.ToLower(email) == "exit" {
			fmt.Println("Exiting program.")
			break
		}

		// Create a new message
		m := NewMessage("Application for Java Developer Position", "Hello,\n\nI recently came across your LinkedIn post regarding the Java developer position and am very interested in this opportunity. Attached, please find my CV for your review. I believe my skills and experience make me a strong candidate for this role.\n\nThank you for considering my application. I look forward to the possibility of discussing how I can contribute to your team.\n\nBest regards,\nBart")
		m.To = []string{email}

		// Attach the file and handle errors
		err := m.AttachFile(filePath)
		if err != nil {
			fmt.Printf("Failed to attach file: %v\n", err)
			continue // Skip sending the email if file attachment fails
		}

		// Set up the SMTP authentication
		auth := smtp.PlainAuth("", username, password, host)

		// Send the email and handle errors
		err = smtp.SendMail(fmt.Sprintf("%s:%s", host, portNumber), auth, username, m.To, m.ToBytes())
		if err != nil {
			fmt.Printf("Failed to send email: %v\n", err)
		} else {
			fmt.Println("Email sent successfully!")

		}
	}
}
