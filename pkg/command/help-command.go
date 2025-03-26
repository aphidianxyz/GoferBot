package command

import (
	"errors"
	"strings"

	telebot "github.com/OvyFlash/telegram-bot-api"
)

type HelpCommand struct {
    chatID int64
    request string
    sendConfig telebot.Chattable
}

func (hc *HelpCommand) GenerateMessage() {
    var config telebot.MessageConfig
    request := strings.TrimPrefix(hc.request, "/")
    switch request {
    case "":
        var allSyntaxes string
        for i, syntaxString := range allHelpSyntaxes {
            if i == len(allHelpSyntaxes) {
                allSyntaxes += syntaxString 
            } else {
                allSyntaxes += syntaxString + "\n"
            }
        }
        config = telebot.NewMessage(hc.chatID, allSyntaxes)
    case "hello":
        config = telebot.NewMessage(hc.chatID, helloSyntax)
    case "help":
        config = telebot.NewMessage(hc.chatID, helpSyntax)
    case "caption":
        config = telebot.NewMessage(hc.chatID, captionSyntax + "\n" + captionImgSyntax)
    case "everyone":
        config = telebot.NewMessage(hc.chatID, everyoneSyntax)
    default:
        invalidCommandMsg := hc.request + " is not a known command"
        config = telebot.NewMessage(hc.chatID, invalidCommandMsg)
    }
    hc.sendConfig = config
}

func (hc *HelpCommand) SendMessage(api *telebot.BotAPI) error {
    if _, err := api.Send(hc.sendConfig); err != nil {
        return errors.New("Failed to send a HelpCommand")
    }
    return nil
}
