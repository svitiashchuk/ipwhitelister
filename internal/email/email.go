// Email operations
package email

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

type Mailer struct {
	APIKey      string
	SenderEmail string
}

func NewMailer(apiKey, senderEmail string) *Mailer {
	return &Mailer{
		APIKey:      apiKey,
		SenderEmail: senderEmail,
	}
}

type EmailContent struct {
	Key     string  `json:"key"` // Mandrill API key
	Message Message `json:"message"`
}

type Message struct {
	Html      string      `json:"html"`
	Text      string      `json:"text"`
	Subject   string      `json:"subject"`
	FromEmail string      `json:"from_email"`
	FromName  string      `json:"from_name"`
	To        []Recipient `json:"to"`
}

type Recipient struct {
	Email string `json:"email"`
	Name  string `json:"name"`
	Type  string `json:"type"`
}

// SendTokenEmail sends an email with the login token
func (m *Mailer) SendTokenEmail(recipientEmail, token string) error {
	emailContent := EmailContent{
		Key: m.APIKey,
		Message: Message{
			Html:      "<p>Here is your token: <b>" + token + "</b></p>",
			Text:      "Here is your token: " + token,
			Subject:   "Token for IP Whitelisting",
			FromEmail: m.SenderEmail,
			FromName:  m.SenderEmail,
			To: []Recipient{
				{
					Email: recipientEmail,
					Name:  recipientEmail,
					Type:  "to",
				},
			},
		},
	}

	jsonData, err := json.Marshal(emailContent)
	if err != nil {
		return err
	}

	// Send the email using Mandrill's API
	resp, err := http.Post("https://mandrillapp.com/api/1.0/messages/send.json", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	respBody := new(bytes.Buffer)
	_, err = respBody.ReadFrom(resp.Body)
	if err != nil {
		return err
	}

	// TODO - check response, log errors if any
	log.Println(resp.Status)
	log.Println(respBody.String())
	defer resp.Body.Close()

	return err
}
