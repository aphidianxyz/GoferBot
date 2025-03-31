package command

import (
	"fmt"
	"strings"

	telebot "github.com/OvyFlash/telegram-bot-api"
)

type HelpCommand struct {
    chatID int64
    request string
	commandJSON CommandJSON
    sendConfig telebot.Chattable
}

func (hc *HelpCommand) GenerateMessage() {
	var cmdInfo string
	trimmed := strings.TrimLeft(hc.request, "/")
	if hc.request == ""{
		cmdInfo = hc.commandJSON.formatAllCommandInfo()
	} else {
		cmdInfo = hc.commandJSON.formatCommandInfo(trimmed)
	}
	if cmdInfo == "" {
		ErrCmdDoesntExist := fmt.Sprintf("Error: command: \"%v\" does not exist", hc.request)
		hc.sendConfig = telebot.NewMessage(hc.chatID, ErrCmdDoesntExist)
		return
	}
	msgConfig := telebot.NewMessage(hc.chatID, cmdInfo) 
	msgConfig.ParseMode = "MarkDown"
	hc.sendConfig = msgConfig
}

func (hc HelpCommand) SendMessage(api *telebot.BotAPI) error {
    if _, err := api.Send(hc.sendConfig); err != nil {
        return err
    }
    return nil
}

func getCommandInfo(request string, commandJSON CommandJSON) string {
	commandInfo := commandJSON.formatCommandInfo(request)
	return commandInfo
}
