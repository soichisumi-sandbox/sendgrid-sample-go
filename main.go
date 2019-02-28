package main

import (
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/sendgrid/sendgrid-go"
	"log"
	"os"
)

type Mail struct {
	Subject          string             `json:"subject"`
	Personalizations []Personalizations `json:"personalizations"`
	From             MailUser           `json:"from"`
	Content          []Contents         `json:"content"`
}

type Personalizations struct {
	To []MailUser `json:"to"`
}

type MailUser struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

type Contents struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Println(err)
	}

	subject := "テスト"
	contents := "テストですよ"

	SendMail(subject, contents)
}

// メールの中身を作成して、メールを送信する
func SendMail(subject, contents string) {
	apiKey := os.Getenv("API_KEY")
	host := "https://api.sendgrid.com"
	endpoint := "/v3/mail/send"

	request := sendgrid.GetRequest(apiKey, endpoint, host)
	request.Method = "POST"

	// メールの内容をJSONで作成する
	mail := Mail{
		Subject: subject,
		Personalizations: []Personalizations{
			{To: []MailUser{{
				Email: os.Getenv("RECEIVER_ADDRESS1"),
				Name:  os.Getenv("RECEIVER_NAME1"),
			},
				{
					Email: os.Getenv("RECEIVER_ADDRESS2"),
					Name:  os.Getenv("RECEIVER_NAME2"),
				},
			}}},
		From: MailUser{
			Email: os.Getenv("SENDER_ADDRESS"),
			Name:  os.Getenv("SENDER_NAME"),
		},
		Content: []Contents{{
			Type:  "text/plain",
			Value: contents,
		}},
	}

	data, err := json.Marshal(mail)

	log.Println(string(data))

	if err != nil {
		log.Println(err)
	}

	request.Body = data

	response, err := sendgrid.API(request)
	if err != nil {
		log.Println(err)
	} else {
		fmt.Println(response.StatusCode)
		fmt.Println(response.Body)
		fmt.Println(response.Headers)
	}
}
