package command

import (
	"strings"

    telebot "github.com/OvyFlash/telegram-bot-api"
)

const (
    helloSyntax string = "/hello - Gofer greets you!"
    helpSyntax string = "/help [command?] - Describes command functionality and syntax, specific command can be specified"
)

var allHelpSyntaxes = []string{helloSyntax, helpSyntax}

type Command interface {
    GenerateMessage() error
    SendMessage(api *telebot.BotAPI) error
}

func ParseMsgCommand(msg *telebot.Message) Command {
    msgStr := msg.Text
    tokens := strings.Split(msgStr, " ")
    commandName := tokens[0]
    commandParams := tokens[1:]
    switch commandName {
    case "/hello": 
        return &HelloCommand{chatID: msg.Chat.ID, firstName: msg.From.FirstName, lastName: msg.From.LastName, userName: msg.From.UserName}
    case "/help":
        var helpRequest string
        if len(commandParams) < 1 {
            helpRequest = ""
        } else {
            helpRequest = commandParams[0]
        }
        return &HelpCommand{chatID: msg.From.ID, request: helpRequest}
    default:
        return &InvalidCommand{chatID: msg.Chat.ID, request: commandName}
    }
} 
