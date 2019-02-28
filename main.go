package main

import (
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
)

// requert.Bodyに格納するJSONの元となるメールを表す構造体
type Mail struct {
	Subject          string             `json:"subject"`
	Personalizations []Personalizations `json:"personalizations"`
	From             MailUser           `json:"from"`
	Content          []Contents         `json:"content"`
}

// 封筒のようなもの
// メールのメタデータを表す構造体
type Personalizations struct {
	To []MailUser `json:"to"`
}

// メールのユーザーを表す構造体
type MailUser struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

// メールの中身を表す構造体
type Contents struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

// .envファイルを読み込んで、ロードする
func Env_load() {
	// .envファイルを読み込んで、ロード
	err := godotenv.Load()
	if err != nil {
		log.Println(err)
	}
}

func main() {
	subject := "テスト"
	contents := "テストですよ"

	SendMail(subject, contents)
}

// メールの中身を作成して、メールを送信する
func SendMail(subject, contents string) {

	Env_load()

	// .envファイルに格納したAPI KEYを取得
	apiKey := os.Getenv("SENDGRID_API_KEY")
	// ホスト
	host := "https://api.sendgrid.com"
	// エンドポイント
	endpoint := "/v3/mail/send"

	// API KEYとエンドポイント、ホストからrestパッケージのRequestを生成
	request := sendgrid.GetRequest(apiKey, endpoint, host)
	// requestのMethodをPostに
	request.Method = "POST"

	// メールの内容をJSONで作成する
	mail := Mail{
		Subject: subject,
		Personalizations: []Personalizations{
			{To: []MailUser{{
				Email: os.Getenv("RECEIVER_USER_ADDRESS1"),
				Name:  os.Getenv("RECEIVER_USER_NAME1"),
			},
				{
					Email: os.Getenv("RECEIVER_USER_ADDRESS2"),
					Name:  os.Getenv("RECEIVER_USER_NAME2"),
				},
			}}},
		From: MailUser{
			Email: os.Getenv("SENDER_USER_ADDRESS"),
			Name:  os.Getenv("SENDER_USER_NAME"),
		},
		Content: []Contents{{
			Type:  "text/plain",
			Value: contents,
		}},
	}

	// GoのコードをJSON化
	data, err := json.Marshal(mail)

	log.Println(string(data))

	if err != nil {
		log.Println(err)
	}

	// JSON化したmailの内容をrequest.Bodyに代入
	request.Body = data

	// sendGridのAPIにリクエストをセット
	// 戻り値でresponseが返ってくる
	response, err := sendgrid.API(request)
	if err != nil {
		log.Println(err)
	} else {
		fmt.Println(response.StatusCode)
		fmt.Println(response.Body)
		fmt.Println(response.Headers)
	}
}
