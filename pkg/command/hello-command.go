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

func MakeHelloCommand(msg telebot.Message) Command {
	return &HelloCommand{chatID: msg.Chat.ID, firstName: msg.From.FirstName, lastName: msg.From.LastName, userName: msg.From.UserName}
}

func (hc *HelloCommand) GenerateMessage() {
    helloString := "Hello, " + hc.firstName + " " + hc.lastName + "!\nAKA: " + hc.userName 
    hc.sendConfig = telebot.NewMessage(hc.chatID, helloString)
} 

func (hc *HelloCommand) SendMessage(api *telebot.BotAPI) error {
    if _, err := api.Send(hc.sendConfig); err != nil {
        return errors.New("Failed to send a HelloCommand")
    }
    return nil
}
