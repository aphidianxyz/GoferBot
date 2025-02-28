package main

import (
	"log"
	"os"

	telebot "github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
    token := os.Getenv("TOKEN")
    bot, err := telebot.NewBotAPI(token)

    if err != nil {
        log.Panic(err)
    }

    bot.Debug = true
    log.Printf("Authorized on account %s", bot.Self.UserName)

    updateConfig := telebot.NewUpdate(0)
    updateConfig.Timeout = 60

    updates, err := bot.GetUpdatesChan(updateConfig)
    if err != nil {
        log.Panic(err)
    }
    for update := range updates {
        if update.Message != nil {
            log.Println("[", update.Message.From.UserName, "]", update.Message.Text)

            msg := telebot.NewMessage(update.Message.Chat.ID, update.Message.Text)
            msg.ReplyToMessageID = update.Message.MessageID

            bot.Send(msg)
        }
    }

}

