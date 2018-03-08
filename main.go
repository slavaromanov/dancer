package main

import (
	"log"
	"net/http"
	"os"
	"time"

	tg "gopkg.in/telegram-bot-api.v4"
)

const one = `ヽ(￣(ｴ)￣)ﾉ`
const two = `(／￣(ｴ)￣)／`

var ids = map[int64]bool{}

func onErr(err error) {
	if err != nil {
		panic(err)
	}
}

func send(bot *tg.BotAPI, c tg.Chattable) int {
	var (
		err error
		msg tg.Message
	)

	if msg, err = bot.Send(c); err != nil {
		return 0
	}
	return msg.MessageID
}

func newM(id int64, t string) tg.MessageConfig {
	m := tg.NewMessage(id, t)
	return m
}

func newMsg(id int64, b *bool) tg.MessageConfig {
	if ids[id] {
		return tg.MessageConfig{}
	}
	ids[id] = true
	log.Println("New chat", id)
	return newM(id, messageFromBool(b))
}

func newEdit(id int64, mid int, b *bool) tg.EditMessageTextConfig {
	m := tg.NewEditMessageText(id, mid, messageFromBool(b))
	return m
}

func messageFromBool(b *bool) string {
	if *b {
		*b = false
		return one
	}
	*b = true
	return two
}

func serve() {
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})))
}

func main() {
	go serve()

	bot, err := tg.NewBotAPI(os.Getenv("BOT_API"))
	onErr(err)

	updateConfig := tg.NewUpdate(0)
	updateConfig.Timeout = 5

	updates, err := bot.GetUpdatesChan(updateConfig)
	onErr(err)

	for update := range updates {
		if update.Message.IsCommand() && update.Message.Command() == "start" {
			p := 0
			id := update.Message.Chat.ID
			b := true

			defer func(prev *int) {
				log.Println("Deleting", p)
				send(bot, tg.DeleteMessageConfig{ChatID: id, MessageID: *prev})
			}(&p)

			go func(prev *int) {
				for mid := send(bot, newMsg(id, &b)); mid != 0; mid = send(bot, newEdit(id, mid, &b)) {
					time.Sleep(time.Millisecond * 500)
					*prev = mid
				}
			}(&p)
		}
	}
}
