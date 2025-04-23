package command

import (
	"fmt"
	"strings"
	"time"

	telebot "github.com/OvyFlash/telegram-bot-api"
)

type Help struct {
	msg telebot.Message
	cmdInfo string
	sendConfig telebot.Chattable
}

func MakeHelp(msg telebot.Message, request string, commandJSON CommandJSON) Command {
	var cmdInfo string
	trimmed := strings.TrimLeft(request, "/")
	if request == "" {
		cmdInfo = commandJSON.formatAllCommandInfo()
	} else {
		cmdInfo = commandJSON.formatCommandInfo(trimmed)
	}
	if cmdInfo == "" {
		cmdDoesntExist := fmt.Sprintf("command: \"%v\" does not exist", request)
		return MakeError(msg, "/help", cmdDoesntExist)
	}
	return &Help{msg: msg, cmdInfo: cmdInfo}
}

func (hc *Help) GenerateMessage() {
	msgConfig := telebot.NewMessage(hc.msg.Chat.ID, hc.cmdInfo) 
	msgConfig.ParseMode = "MarkDown"
	hc.sendConfig = msgConfig
}

func (hc Help) SendMessage(api *telebot.BotAPI) error {
	var msg telebot.Message
	var err error
    if msg, err = api.Send(hc.sendConfig); err != nil {
        return err
    }
	deleteMessage(api, 2 * time.Minute, msg.Chat.ID, msg.MessageID)
	deleteMessage(api, 2 * time.Minute, hc.msg.Chat.ID, hc.msg.MessageID)
    return nil
}

func getCommandInfo(request string, commandJSON CommandJSON) string {
	commandInfo := commandJSON.formatCommandInfo(request)
	return commandInfo
}
