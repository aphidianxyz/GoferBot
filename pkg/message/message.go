package message

import (
	//telebot "github.com/OvyFlash/telegram-bot-api"
)

type Message interface {
	GetID()
	GetText()
	GetPhoto()
	GetReply()
}
