package command

import (
	"fmt"
	"strings"
	"time"

	telebot "github.com/OvyFlash/telegram-bot-api"
)

type HelpCommand struct {
	msg telebot.Message
    request string
	commandJSON CommandJSON
    sendConfig telebot.Chattable
	failed bool
}

func (hc *HelpCommand) GenerateMessage() {
	hc.failed = false
	var cmdInfo string
	trimmed := strings.TrimLeft(hc.request, "/")
	if hc.request == ""{
		cmdInfo = hc.commandJSON.formatAllCommandInfo()
	} else {
		cmdInfo = hc.commandJSON.formatCommandInfo(trimmed)
	}
	if cmdInfo == "" {
		hc.failed = true
		ErrCmdDoesntExist := fmt.Sprintf("Error: command: \"%v\" does not exist", hc.request)
		hc.sendConfig = telebot.NewMessage(hc.msg.Chat.ID, ErrCmdDoesntExist)
		return
	}
	msgConfig := telebot.NewMessage(hc.msg.Chat.ID, cmdInfo) 
	msgConfig.ParseMode = "MarkDown"
	hc.sendConfig = msgConfig
}

func (hc HelpCommand) SendMessage(api *telebot.BotAPI) error {
	var msg telebot.Message
	var err error
    if msg, err = api.Send(hc.sendConfig); err != nil {
        return err
    }
	var delay time.Duration = 60 * time.Second
	if hc.failed {
		delay = 10 * time.Second
	}
	deleteMessage(api, delay, msg.Chat.ID, msg.MessageID)
	deleteMessage(api, delay, hc.msg.Chat.ID, hc.msg.MessageID)
    return nil
}

func getCommandInfo(request string, commandJSON CommandJSON) string {
	commandInfo := commandJSON.formatCommandInfo(request)
	return commandInfo
}
