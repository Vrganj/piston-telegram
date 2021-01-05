package main

import (
	"log"
	"time"

	telebot "gopkg.in/tucnak/telebot.v2"
)

func main() {
	config, err := loadConfig()

	if err != nil {
		log.Println(err)
		return
	}

	if config.Token == "" {
		log.Println("Provide a token in config.json")
		return
	}

	bot, err := telebot.NewBot(telebot.Settings{
		Token:  config.Token,
		Poller: &telebot.LongPoller{Timeout: 5 * time.Second},
	})

	if err != nil {
		log.Println(err)
		return
	}

	bot.Handle("/run", func(message *telebot.Message) {
		runCommand(message, bot)
	})

	bot.Start()
}
