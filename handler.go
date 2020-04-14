package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/line/line-bot-sdk-go/linebot"
)

// Handler handles message processing from linebot, sending reply message, and Get / Put to datastore.
type Handler struct {
	client *linebot.Client
	regist *Registration
}

// NewHandler creates handler
func NewHandler(client *linebot.Client, regist *Registration) *Handler {
	return &Handler{
		client: client,
		regist: regist,
	}
}

// WebhookHandler hooks messages from linebot
func (h *Handler) WebhookHandler(w http.ResponseWriter, r *http.Request) {
	events, err := h.client.ParseRequest(r)
	if err != nil {
		log.Print(err)
		return
	}
	for _, event := range events {
		if event.Type != linebot.EventTypeMessage {
			return
		}

		switch message := event.Message.(type) {
		case *linebot.TextMessage:
			h.replyMessageExec(event, message)

		case *linebot.StickerMessage:
			replyMessage := fmt.Sprintf(
				"sticker id is %s, stickerResourceType is ...", message.StickerID)
			if _, err := h.client.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyMessage)).Do(); err != nil {
				log.Print(err)
			}
		}
	}
}

// replyMessageExec sends a reply message to linebot
func (h *Handler) replyMessageExec(event *linebot.Event, message *linebot.TextMessage) {
	switch message.Text {
	case "入力":
		url := "https://foodbalance001.appspot.com/input" + "?userid=" + event.Source.UserID
		linebotMessage := "登録には以下のリンクをクリックしてください." + url
		resp := linebot.NewTextMessage(linebotMessage)
		_, err := h.client.ReplyMessage(event.ReplyToken, resp).Do()
		if err != nil {
			log.Print(err)
		}
	case "表示":
		// 当日の集計を表示
		date := time.Now().Format(dateFormat)
		query := datastore.NewQuery("RegistrationData").Filter("UserID = ", event.Source.UserID).Filter("Date = ", date)
		var regists []RegistrationData
		if err := h.regist.GetAll(context.Background(), query, &regists); err != nil {
			log.Print("Get失敗", err)
			return
		}

		group := SumGroup(regists)
		message := "今日食べたものは...\n"
		for _, regist := range regists {
			message = message + fmt.Sprintf("・%s\n", regist.Name)
		}
		message = message + "\n"
		message = message + fmt.Sprintf("主食 : %d\n副菜 : %d\n主菜 : %d\n乳製品 : %d\n果物 : %d\n",
			group.GrainDishes, group.VegetableDishes, group.FishAndMealDishes, group.Milk, group.Fruit)
		resp := linebot.NewTextMessage(message)
		_, err := h.client.ReplyMessage(event.ReplyToken, resp).Do()
		if err != nil {
			log.Print(err)
		}
	}
}

// FormField is a field to replace with the input form
type formField struct {
	Userid string
}

type errorField struct {
	Detail string
}

// InputformHandler handles the input form
func (h *Handler) InputformHandler(w http.ResponseWriter, r *http.Request) {
	v := r.URL.Query()
	if v == nil {
		tpl := template.Must(template.ParseFiles("error.html"))
		tpl.Execute(w, errorField{Detail: "Not find URL value"})
		return
	}

	var fd formField
	for key, vs := range v {
		log.Printf("%s = %s\n", key, vs[0])
		fd.Userid = vs[0] // UserIDを入力フォームのsubmitの際にクエリとして与える
	}
	tpl := template.Must(template.ParseFiles("input.html"))
	tpl.Execute(w, fd)
}

// PostHandler handles post requests from forms
func (h *Handler) PostHandler(w http.ResponseWriter, r *http.Request) {
	v := r.URL.Query()
	if v == nil {
		tpl := template.Must(template.ParseFiles("error.html"))
		tpl.Execute(w, errorField{Detail: "Not find URL value"})
		return
	}

	var userid string
	for key, vs := range v {
		log.Printf("%s = %s\n", key, vs[0])
		userid = vs[0] // UserIDをdatastoreに登録するデータのメンバに含める
	}
	if r.Method != http.MethodPost {
		return
	}

	registData, err := h.convertRegistration(r, userid)
	if err != nil {
		tpl := template.Must(template.ParseFiles("error.html"))
		tpl.Execute(w, errorField{Detail: "Error convert parameter"})
		return
	}

	err = h.regist.Put(context.Background(), datastore.NameKey("RegistrationData", "", nil), registData)
	if err != nil {
		tpl := template.Must(template.ParseFiles("error.html"))
		tpl.Execute(w, errorField{Detail: "Error writing to datastore"})
		return
	}

	tpl := template.Must(template.ParseFiles("success.html"))
	tpl.Execute(w, nil)
}

// convertRegistration parses data from the input form and converts it to type Registration
func (h *Handler) convertRegistration(r *http.Request, userid string) (*RegistrationData, error) {
	if err := r.ParseForm(); err != nil {
		return nil, err
	}
	name := r.Form["Name"][0]
	if len(name) == 0 {
		return nil, fmt.Errorf("Error: %s", "Name is empty")
	}
	timeZone, err := strconv.Atoi(r.Form["TimeZone"][0])
	if err != nil {
		timeZone = 0
	}
	is, err := Astois(r.Form["Group"])
	if err != nil {
		return nil, err
	}
	ts := time.Now().Format(dateFormat)
	group := Group{is[0], is[1], is[2], is[3], is[4]}
	return NewRegistationData(userid, ts, name, timeZone, group), nil
}

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

// SumGroup returns sum of groups
func SumGroup(regists []RegistrationData) Group {
	var group Group
	for _, r := range regists {
		group.Sum(r.BalanceGroup)
	}
	return group
}

// Sum adds the entered Group to your own Group
func (g *Group) Sum(in Group) {
	g.GrainDishes += in.GrainDishes
	g.VegetableDishes += in.VegetableDishes
	g.FishAndMealDishes += in.FishAndMealDishes
	g.Milk += in.Milk
	g.Fruit += in.Fruit
}
