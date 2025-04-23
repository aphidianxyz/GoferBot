package command

import (
	"strings"

	telebot "github.com/OvyFlash/telegram-bot-api"
)

type Pin struct {
	msg telebot.Message
	sendConfig telebot.Chattable
}

func MakePin(msg telebot.Message) Command {
	if msg.ReplyToMessage == nil {
		pinErrStr := "Please reply to the message you want to pin, with /pin\nReplying to messages before the bot was added without visible chat history will also not work."
		return MakeError(msg, "/pin", pinErrStr)
	}
	return &Pin{msg: msg}
}

func (pc *Pin) GenerateMessage() {
	var replyTargetID int
	replyTargetID = pc.msg.ReplyToMessage.MessageID
	msgTokens := strings.Split(pc.msg.Text, " ")
	notificationFlag := len(msgTokens) > 1 && msgTokens[1] == "-notify"
	pinConfig := telebot.NewPinChatMessage(pc.msg.Chat.ID, replyTargetID, notificationFlag)
	pc.sendConfig = pinConfig
}

func (pc Pin) SendMessage(api *telebot.BotAPI) error {
    if _, err := api.Request(pc.sendConfig); err != nil {
        return err
    }
	return nil
}
