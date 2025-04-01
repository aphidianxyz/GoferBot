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

const (
	invalidCmdSuffix = " is not a valid command!"
)

// TODO: maybe this would be better described as "UnknownCommand"
func (ic *InvalidCommand) GenerateMessage() {
    invalidRequest := ic.request + invalidCmdSuffix
    ic.sendConfig = telebot.NewMessage(ic.chatID, invalidRequest)
}

func (ic *InvalidCommand) SendMessage(api *telebot.BotAPI) error {
    if _, err := api.Send(ic.sendConfig); err != nil {
        return errors.New("Failed to send an InvalidCommand")
    }
    return nil
}
