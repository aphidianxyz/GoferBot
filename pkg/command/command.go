package command

import (
	telebot "github.com/OvyFlash/telegram-bot-api"
)

type Command interface {
    GenerateMessage()
    SendMessage(api *telebot.BotAPI) error
}
