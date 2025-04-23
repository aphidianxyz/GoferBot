package command

import (
	"errors"
	"time"

	telebot "github.com/OvyFlash/telegram-bot-api"
)

type Error struct {
	msg telebot.Message
	originCmdName string
	errMsg string
	sendConfig telebot.Chattable
}

func MakeError(msg telebot.Message, originCmdName, errMsg string) Command {
	return &Error{msg: msg, originCmdName: originCmdName, errMsg: errMsg}
}

func (e *Error) GenerateMessage() {
	errorMessage := e.originCmdName + " - " + e.errMsg
	e.sendConfig = telebot.NewMessage(e.msg.Chat.ID, errorMessage)
}

func (e *Error) SendMessage(api *telebot.BotAPI) error {
	var msg telebot.Message
	var err error
    if msg, err = api.Send(e.sendConfig); err != nil {
        return errors.New("Failed to send an ErrorCommand")
    }
	deleteMessage(api, 10 * time.Second, e.msg.Chat.ID, e.msg.MessageID)
	deleteMessage(api, 10 * time.Second, e.msg.Chat.ID, msg.MessageID)
    return nil
}
