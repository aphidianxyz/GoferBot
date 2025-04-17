package command

import (
	"log"
	"time"

	telebot "github.com/OvyFlash/telegram-bot-api"
)

type Command interface {
	// SetInputs() Command -- verify inputs, create a command 
    GenerateMessage()
    SendMessage(api *telebot.BotAPI) error
}

func deleteMessage(api *telebot.BotAPI, delay time.Duration, chatID int64, msgID int) {
	deleteRequest := telebot.NewDeleteMessage(chatID, msgID)
	go func() {
		time.Sleep(delay)
		if _, err := api.Request(deleteRequest); err != nil {
			log.Printf("Failed to delete message: %v", err)
		}
	}()
}
