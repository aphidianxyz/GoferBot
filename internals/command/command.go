package command

import (
	"strings"

	telebot "github.com/go-telegram-bot-api/telegram-bot-api"
)

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
    case "/hello": // TODO: utilize enums
        return &HelloCommand{chatID: msg.Chat.ID, firstName: msg.From.FirstName, lastName: msg.From.LastName, userName: msg.From.UserName}
    case "/help":
        var helpRequest string
        if len(commandParams) < 1 {
            helpRequest = ""
        } else {
            helpRequest = commandParams[0]
        }
        // TODO: maybe make /help DM the requesting user instead 
        return &HelpCommand{chatID: msg.Chat.ID, request: helpRequest}
    default:
        return &InvalidCommand{chatID: msg.Chat.ID, request: commandName}
    }
} 
