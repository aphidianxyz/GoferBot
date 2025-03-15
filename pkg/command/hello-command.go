package command

import (
    "errors"
    
    telebot "github.com/OvyFlash/telegram-bot-api"
)

type HelloCommand struct {
    chatID int64
    firstName, lastName, userName string
    sendConfig telebot.Chattable
}

func (hc *HelloCommand) GenerateMessage() {
    helloString := "Hello, " + hc.firstName + " " + hc.lastName + "!\nAKA: " + hc.userName 
    config := telebot.NewMessage(hc.chatID, helloString)
    hc.sendConfig = config 
} 

func (hc *HelloCommand) SendMessage(api *telebot.BotAPI) error {
    if _, err := api.Send(hc.sendConfig); err != nil {
        return errors.New("Failed to send a HelloCommand")
    }
    return nil
}
