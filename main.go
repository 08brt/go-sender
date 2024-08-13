package main

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"fmt"
	"log"
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

func NewMessage(subject, body string) *Message {
	return &Message{Subject: subject, Body: body, Attachments: make(map[string][]byte)}
}

func (m *Message) AttachFile(src string) error {
	b, err := os.ReadFile(src)
	if err != nil {
		return fmt.Errorf("error reading file %s: %w", src, err)
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
		writeBoundary(writer, buf)
	} else {
		buf.WriteString("Content-Type: text/plain; charset=utf-8\n")
	}

	buf.WriteString(m.Body)
	if withAttachments {
		for fileName, content := range m.Attachments {
			writeBoundary(writer, buf)
			buf.WriteString(fmt.Sprintf("Content-Type: %s\n", http.DetectContentType(content)))
			buf.WriteString("Content-Transfer-Encoding: base64\n")
			buf.WriteString(fmt.Sprintf("Content-Disposition: attachment; filename=%s\n", fileName))

			b64Content := make([]byte, base64.StdEncoding.EncodedLen(len(content)))
			base64.StdEncoding.Encode(b64Content, content)
			buf.Write(b64Content)
			buf.WriteString(fmt.Sprintf("\n--%s", boundary))
		}

		buf.WriteString("--")
	}

	return buf.Bytes()
}

func writeBoundary(writer *multipart.Writer, buf *bytes.Buffer) {
	buf.WriteString(fmt.Sprintf("--%s\n", writer.Boundary()))
}

func getUserInput(reader *bufio.Reader) (string, error) {
	fmt.Print("Enter email address (or 'exit'): ")
	email, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("error reading input: %w", err)
	}
	return strings.TrimSpace(email), nil
}

func sendEmail(m *Message) error {
	auth := smtp.PlainAuth("", username, password, host)
	return smtp.SendMail(fmt.Sprintf("%s:%s", host, portNumber), auth, username, m.To, m.ToBytes())
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	for {
		email, err := getUserInput(reader)
		if err != nil {
			log.Printf("Error: %v\n", err)
			continue
		}

		if strings.ToLower(email) == "exit" {
			fmt.Println("Exiting program.")
			break
		}

		m := NewMessage("Application for Java Developer Position", "Hello,\n\nI recently came across your LinkedIn post regarding the Java developer position and am very interested in this opportunity. Attached, please find my CV for your review. I believe my skills and experience make me a strong candidate for this role.\n\nThank you for considering my application. I look forward to the possibility of discussing how I can contribute to your team.\n\nBest regards,\nBart")
		m.To = []string{email}

		if err := m.AttachFile(filePath); err != nil {
			log.Printf("Failed to attach file: %v\n", err)
			continue
		}

		if err := sendEmail(m); err != nil {
			log.Printf("Failed to send email: %v\n", err)
		} else {
			fmt.Println("Email sent successfully!")
		}
	}
}
