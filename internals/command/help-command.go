package command

import (
    "strings"
    "errors"

	telebot "github.com/go-telegram-bot-api/telegram-bot-api"
)

type HelpCommand struct {
    chatID int64
    request string
    sendConfig telebot.MessageConfig
}

func (hc *HelpCommand) GenerateMessage() error {
    var config telebot.MessageConfig
    request := strings.TrimPrefix(hc.request, "/")
    switch request {
    case "":
        config = telebot.NewMessage(hc.chatID, "all commands: ...")
    case "hello":
        config = telebot.NewMessage(hc.chatID, "/hello - Gofer says hello to you!")
    case "help":
        helpCommandDesc := "/help [command?] - Describes command functionality and syntax, specific command can be specified"
        config = telebot.NewMessage(hc.chatID, helpCommandDesc)
    default:
        invalidCommandMsg := hc.request + " is not a known command"
        config = telebot.NewMessage(hc.chatID, invalidCommandMsg)
    }
    hc.sendConfig = config
    return nil
}

func (hc *HelpCommand) SendMessage(api *telebot.BotAPI) error {
    if _, err := api.Send(hc.sendConfig); err != nil {
        return errors.New("Failed to send a HelpCommand")
    }
    return nil
}
