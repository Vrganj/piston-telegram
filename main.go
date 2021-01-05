package main

import (
	"time"

	telebot "gopkg.in/tucnak/telebot.v2"
)

func main() {
	config, err := loadConfig()

	if err != nil {
		panic(err)
	}

	bot, err := telebot.NewBot(telebot.Settings{
		Token:  config.Token,
		Poller: &telebot.LongPoller{Timeout: 5 * time.Second},
	})

	if err != nil {
		panic(err)
	}

	bot.Handle("/run", func(message *telebot.Message) {
		runCommand(message, bot)
	})

	bot.Start()
}
