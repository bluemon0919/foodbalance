package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/leboncoin/dialogflow-go-webhook"
	df "github.com/leboncoin/dialogflow-go-webhook"
	"github.com/line/line-bot-sdk-go/linebot"
)

// Astois is equivalent to Atoi for slice.
func Astois(ss []string) ([]int, error) {
	var is []int
	for _, s := range ss {
		i, err := strconv.Atoi(s)
		if err != nil {
			return nil, err
		}
		is = append(is, i)
	}
	return is, nil
}

// WebpageHandler handles input form
func WebpageHandler(w http.ResponseWriter, r *http.Request) {
	tpl := template.Must(template.ParseFiles("input.html"))
	tpl.Execute(w, nil)
}

// WebpagePostHandler handles input form posts and redirects
func WebpagePostHandler(w http.ResponseWriter, r *http.Request) {
	// HTTPメソッドをチェック（POSTのみ許可）
	if r.Method != http.MethodPost {
		return
	}
	r.ParseForm()

	name := r.Form["Name"][0]
	if len(name) == 0 {
		return
	}
	timeZone, err := strconv.Atoi(r.Form["TimeZone"][0])
	if err != nil {
		timeZone = 0
	}
	is, err := Astois(r.Form["Group"])
	if err != nil {
		return
	}
	t := time.Now()
	ts := t.Format(dateFormat)
	m := Create(ts, name, timeZone, is[0], is[1], is[2], is[3], is[4])
	Put(m)

	http.Redirect(w, r, "/", 303)
}

// DialogflowParam is the input structure from DialogFlow
type DialogflowParam struct {
	Name string `json:"name"`
}

// DialogflowHandler handles posts from DialogFlow
func DialogflowHandler(w http.ResponseWriter, r *http.Request) {
	var dfRequest *df.Request
	if err := json.NewDecoder(r.Body).Decode(&dfRequest); err != nil {
		code := http.StatusBadRequest
		log.Println("Error:", err)
		http.Error(w, http.StatusText(code), code)
		return
	}

	// https://cloud.google.com/dialogflow/docs/fulfillment-how?hl=ja
	switch dfRequest.QueryResult.Intent.DisplayName {
	case "put":
		DialogflowPost(w, dfRequest)
	}
}

var linebotCannelSecret string
var linebotCannelAccessToken string

func init() {
	linebotCannelSecret = os.Getenv("CannelSecret")
	linebotCannelAccessToken = os.Getenv("CannelAccessToken")
}

/*
func LinebotMessage() {
	client := &http.Client{}
	bot, err := linebot.New(linebotCannelSecret, linebotCannelAccessToken, linebot.WithHTTPClient(client))
}
*/

// Webhook receives an http request and sends a message to LineBot
func Webhook(w http.ResponseWriter, r *http.Request) {
	/*
		events, err := linebot.ParseRequest(linebotCannelSecret, r)
		if err != nil {
			if err == linebot.ErrInvalidSignature {
				w.WriteHeader(400)
			} else {
				w.WriteHeader(500)
			}
			return
		}
		for _, event := range events {
			log.Printf("Got event %v", event)
			switch event.Type {
			case linebot.EventTypeMessage:
				switch message := event.Message.(type) {
				case *linebot.TextMessage:

					fmt.Println(message)
				}
			}
		}
	*/

	// ここから先は動く
	msg := linebot.NewTextMessage("届きました")
	dff := &dialogflow.Fulfillment{
		FulfillmentMessages: dialogflow.Messages{
			dialogflow.Message{
				Platform: dialogflow.Line,
				RichMessage: dialogflow.PayloadWrapper{Payload: map[string]interface{}{
					"line": msg,
				}},
			},
		},
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(dff); err != nil {
		log.Println("Error:", err)
	}
}

func DialogflowPost(w http.ResponseWriter, dfRequest *df.Request) {
	// Dialogflowへの応答メッセージを返す.
	var dfParam DialogflowParam
	if err := json.Unmarshal([]byte(dfRequest.QueryResult.Parameters), &dfParam); err != nil {
		code := http.StatusBadRequest
		log.Println("Error:", err)
		http.Error(w, http.StatusText(code), code)
		return
	}

	t := time.Now()
	ts := t.Format(dateFormat)
	m := Create(ts, dfParam.Name, 1, 1, 1, 1, 1, 1) // TODO : name以外を入力する
	Put(m)
}
