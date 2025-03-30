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
	var syntax string
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
		syntax = allSyntaxes
	case "about":
		syntax = aboutSyntax
    case "caption":
		syntax = captionSyntax
    case "everyone":
		syntax = everyoneSyntax
    case "hello":
		syntax = helloSyntax
    case "help":
		syntax = helpSyntax
    default:
        syntax = hc.request + " is not a known command"
    }
	hc.sendConfig = telebot.NewMessage(hc.chatID, syntax)
}

func (hc *HelpCommand) SendMessage(api *telebot.BotAPI) error {
    if _, err := api.Send(hc.sendConfig); err != nil {
        return errors.New("Failed to send a HelpCommand")
    }
    return nil
}
