package command

import (
	"errors"
	"time"

	telebot "github.com/OvyFlash/telegram-bot-api"
)

type ErrorCommand struct {
	msg telebot.Message
	originCommand string
	error string
	sendConfig telebot.Chattable
}

func (ec *ErrorCommand) GenerateMessage() {
	errorMessage := ec.originCommand + " failed: " + ec.error
	ec.sendConfig = telebot.NewMessage(ec.msg.Chat.ID, errorMessage)
}

func (ec *ErrorCommand) SendMessage(api *telebot.BotAPI) error {
	var msg telebot.Message
	var err error
    if msg, err = api.Send(ec.sendConfig); err != nil {
        return errors.New("Failed to send an InvalidCommand")
    }
	deleteMessage(api, 10 * time.Second, ec.msg.Chat.ID, ec.msg.MessageID)
	deleteMessage(api, 10 * time.Second, ec.msg.Chat.ID, msg.MessageID)
    return nil
}
