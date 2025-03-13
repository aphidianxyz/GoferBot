package command

import (
    "errors"

	telebot "github.com/go-telegram-bot-api/telegram-bot-api"
)

type HelloCommand struct {
    chatID int64
    firstName, lastName, userName string
    sendConfig telebot.MessageConfig
}

func (hc *HelloCommand) GenerateMessage() error {
    helloString := "Hello, " + hc.firstName + " " + hc.lastName + "!\nAKA: " + hc.userName 
    config := telebot.NewMessage(hc.chatID, helloString)
    hc.sendConfig = config 
    return nil
} 

func (hc *HelloCommand) SendMessage(api *telebot.BotAPI) error {
    if _, err := api.Send(hc.sendConfig); err != nil {
        return errors.New("Failed to send a HelloCommand")
    }
    return nil
}
