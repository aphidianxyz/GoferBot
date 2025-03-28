package command

import (
	"database/sql"
	"strings"

	telebot "github.com/OvyFlash/telegram-bot-api"
)

// maybe this should be a data file, i.e. JSON
const (
    helloSyntax string = "/hello - Gofer greets you!"
    helpSyntax string = "/help [command?] - Describes command functionality and syntax\ncommand (optional) - a specific command to describe"
    captionSyntax string = "/caption [url] [\"top\"] [\"bot\"] - Creates an impact font caption meme\nurl - the url of the image to be captioned\n\"top\" - the top caption, encapsulated by quotes\n\"bot\" - the bottom caption, encapsulated by quotes"
    captionImgSyntax string = "/caption [\"top\"] [\"bot\"] (with an image attached) - Creates an impact font caption meme\n\"top\" - the top caption, encapsulated by quotes\n\"bot\" - the bottom caption, encapsulated by quotes"
    everyoneSyntax string = "/everyone [message?]\nmessage (optional) - a message that accompanies a ping to everyone"
)

var allHelpSyntaxes = []string{helloSyntax, helpSyntax, captionSyntax, captionImgSyntax, everyoneSyntax}

type Command interface {
    GenerateMessage()
    SendMessage(api *telebot.BotAPI) error
}

func ParseMsgCommand(api *telebot.BotAPI, chatDB *sql.DB, msg *telebot.Message) Command {
    msgTxt := strings.Split(msg.Text, "\n")[0] // cmds should only be on the first line
    tokens := strings.Split(msgTxt, " ")
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
        return &HelpCommand{chatID: msg.Chat.ID, request: helpRequest}
    case "/caption":
        return &CaptionCommand{msg: *msg}
    case "/everyone":
        return &EveryoneCommand{msg: *msg, db: chatDB}
    default:
        return &InvalidCommand{chatID: msg.Chat.ID, request: commandName}
    }
} 

func ParseImgCommand(api *telebot.BotAPI, chatDB *sql.DB, msg *telebot.Message) Command {
    msgCap := msg.Caption
    msgCap = strings.Split(msgCap, "\n")[0] // cmds should only be on the first line
    tokens := strings.Split(msgCap, " ")
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
        return &HelpCommand{chatID: msg.Chat.ID, request: helpRequest}
    case "/caption":
        return &CaptionImgCommand{api: api, msg: *msg}
    case "/everyone":
        return &EveryoneCommand{msg: *msg, db: chatDB}
    default:
        return &InvalidCommand{chatID: msg.Chat.ID, request: commandName}
    }
}
