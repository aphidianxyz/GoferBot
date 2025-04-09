package command

import (
	"strings"

	telebot "github.com/OvyFlash/telegram-bot-api"
)

type PinCommand struct {
	api *telebot.BotAPI
	msg telebot.Message
	sendConfig telebot.Chattable
}

func (pc *PinCommand) GenerateMessage() {
	var replyTargetID int
	if pc.msg.ReplyToMessage == nil {
		pc.sendConfig = telebot.NewMessage(pc.msg.Chat.ID, "Please reply to the message you want to pin, with /pin\nReplying to messages before the bot was added without visible chat history will also not work.")
		return
	}
	replyTargetID = pc.msg.ReplyToMessage.MessageID
	msgTokens := strings.Split(pc.msg.Text, " ")
	notificationFlag := len(msgTokens) > 1 && msgTokens[1] == "-notify"
	pinConfig := telebot.NewPinChatMessage(pc.msg.Chat.ID, replyTargetID, notificationFlag)
	pc.sendConfig = pinConfig
}

func (pc PinCommand) SendMessage(api *telebot.BotAPI) error {
    if _, err := api.Request(pc.sendConfig); err != nil {
        return err
    }
	return nil
}
