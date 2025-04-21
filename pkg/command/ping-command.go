package command

import (
	"errors"
	"time"

	telebot "github.com/OvyFlash/telegram-bot-api"
)

type PingCommand struct {
	msg telebot.Message
    sendConfig telebot.Chattable
}

func MakePingCommand(msg telebot.Message) Command {
	return &PingCommand{msg: msg}
}

func (pc *PingCommand) GenerateMessage() {
    config := telebot.NewMessage(pc.msg.Chat.ID, "pong!")
	config.ReplyParameters.MessageID = pc.msg.MessageID
    pc.sendConfig = config
}

func (pc *PingCommand) SendMessage(api *telebot.BotAPI) error {
	var msg telebot.Message
	var err error
    if msg, err = api.Send(pc.sendConfig); err != nil {
        return errors.New("Failed to send a PingCommand")
    }
	deleteMessage(api, 5 * time.Second, msg.Chat.ID, msg.MessageID)
	deleteMessage(api, 5 * time.Second, msg.Chat.ID, pc.msg.MessageID)
    return nil
}
