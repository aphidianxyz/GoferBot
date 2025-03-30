package command

import (
    "errors"
    
    telebot "github.com/OvyFlash/telegram-bot-api"
)

type PingCommand struct {
    chatID int64
    sendConfig telebot.Chattable
}

func (pc *PingCommand) GenerateMessage() {
    config := telebot.NewMessage(pc.chatID, "pong")
    pc.sendConfig = config
}

func (pc *PingCommand) SendMessage(api *telebot.BotAPI) error {
    if _, err := api.Send(pc.sendConfig); err != nil {
        return errors.New("Failed to send a PingCommand")
    }
    return nil
}
