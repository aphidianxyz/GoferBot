package command

import (
	"strings"
	"time"

	telebot "github.com/OvyFlash/telegram-bot-api"
)

type PinCommand struct {
	api *telebot.BotAPI
	msg telebot.Message
	sendConfig telebot.Chattable
	failed bool
}

func (pc *PinCommand) GenerateMessage() {
	var replyTargetID int
	pc.failed = false 
	if pc.msg.ReplyToMessage == nil {
		pc.sendConfig = telebot.NewMessage(pc.msg.Chat.ID, "Please reply to the message you want to pin, with /pin\nReplying to messages before the bot was added without visible chat history will also not work.")
		pc.failed = true
		return
	}
	replyTargetID = pc.msg.ReplyToMessage.MessageID
	msgTokens := strings.Split(pc.msg.Text, " ")
	notificationFlag := len(msgTokens) > 1 && msgTokens[1] == "-notify"
	pinConfig := telebot.NewPinChatMessage(pc.msg.Chat.ID, replyTargetID, notificationFlag)
	pc.sendConfig = pinConfig
}

func (pc PinCommand) SendMessage(api *telebot.BotAPI) error {
	if pc.failed { // HACK: don't like how we're making sendmessage handle this specific case
		var msg telebot.Message
		var err error
		if msg, err = api.Send(pc.sendConfig); err != nil {
			return err
		}
		deleteMessage(api, 10 * time.Second, msg.Chat.ID, pc.msg.MessageID)
		deleteMessage(api, 10 * time.Second, msg.Chat.ID, msg.MessageID)
		return nil
	}
    if _, err := api.Request(pc.sendConfig); err != nil {
        return err
    }
	return nil
}
