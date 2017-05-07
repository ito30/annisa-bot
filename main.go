package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
	"github.com/line/line-bot-sdk-go/linebot"
)

var (
	Bot                *linebot.Client
	channelAccessToken string
	channelSecret      string
)

func Hello(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprint(w, "Welcome!\n")
}

func LineReq(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	events, err := Bot.ParseRequest(r)

	if err != nil {
		fmt.Fprint(w, err)
	}

	var replyToken string

	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			fmt.Fprint(w, event.Type)
			replyToken = event.ReplyToken
		}
	}

	var messages []linebot.Message

	messages = append(messages, linebot.NewTextMessage("Hello~"))

	_, err = Bot.ReplyMessage(replyToken, messages...).Do()

	if err != nil {
		fmt.Fprint(w, err)
	}
}

func readConfig() map[string]string {
	var config map[string]string

	file, e := ioutil.ReadFile("./config.json")
	if e != nil {
		fmt.Printf("File error: %v\n", e)
		os.Exit(1)
	}

	json.Unmarshal(file, &config)

	return config
}

func Init() {
	config := readConfig()

	bot, err := linebot.New(config["channel_secret"], config["channel_access_token"])

	if err != nil {
		log.Println(err)
	}

	Bot = bot
}

func main() {
	Init()

	router := httprouter.New()
	router.GET("/hello", Hello)
	router.POST("/new", LineReq)

	log.Fatal(http.ListenAndServe(":8080", router))
}
