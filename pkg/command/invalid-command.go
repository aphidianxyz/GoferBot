package command

import (
    "errors"

	telebot "github.com/OvyFlash/telegram-bot-api"
)

type InvalidCommand struct {
    chatID int64
    request string
    sendConfig telebot.MessageConfig
}

// TODO: maybe this would be better described as "UnknownCommand"
func (ic *InvalidCommand) GenerateMessage() error {
    invalidRequest := ic.request + " is not a valid command!"
    ic.sendConfig = telebot.NewMessage(ic.chatID, invalidRequest)
    return nil
}

func (ic *InvalidCommand) SendMessage(api *telebot.BotAPI) error {
    if _, err := api.Send(ic.sendConfig); err != nil {
        return errors.New("Failed to send an InvalidCommand")
    }
    return nil
}
