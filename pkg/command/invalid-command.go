package command

import (
	"errors"
	"time"

	telebot "github.com/OvyFlash/telegram-bot-api"
)

type InvalidCommand struct {
	msg telebot.Message
    request string
    sendConfig telebot.MessageConfig
}

// TODO: maybe this would be better described as "UnknownCommand"
func (ic *InvalidCommand) GenerateMessage() {
    invalidRequest := ic.request + " is not a valid command!"
    ic.sendConfig = telebot.NewMessage(ic.msg.Chat.ID, invalidRequest)
}

func (ic *InvalidCommand) SendMessage(api *telebot.BotAPI) error {
	var msg telebot.Message
	var err error
    if msg, err = api.Send(ic.sendConfig); err != nil {
        return errors.New("Failed to send an InvalidCommand")
    }
	deleteMessage(api, 10 * time.Second, ic.msg.Chat.ID, ic.msg.MessageID)
	deleteMessage(api, 10 * time.Second, ic.msg.Chat.ID, msg.MessageID)
    return nil
}
