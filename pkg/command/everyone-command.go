package command

import (
    telebot "github.com/OvyFlash/telegram-bot-api"
)

type EveryoneCommand struct {
    chatID int64
    allUsers map[int]string
}

func (ec *EveryoneCommand) GenerateMessage() {
    
}

func (ec *EveryoneCommand) SendMessage(api *telebot.BotAPI) error {
    return nil
}
