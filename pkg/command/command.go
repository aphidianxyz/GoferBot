package command

import (
	"database/sql"
	"strings"

	telebot "github.com/OvyFlash/telegram-bot-api"
)

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
		url, err := getUrl(msg.Text)
		if err != nil {
			return buildReplyCaptionCommand(api, msg)
		}
		return &CaptionCommand{msg: *msg, url: url}
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

func buildReplyCaptionCommand(api *telebot.BotAPI, msg *telebot.Message) Command {
	// a caption command can work with a reply to an image
	if exReply := msg.ExternalReply; exReply != nil  { // can't reference a pointer if nil, so double check is necessary
		if exReply.Photo != nil {
			replyMsg := *msg
			replyMsg.Photo = msg.ExternalReply.Photo
			replyMsg.Caption = msg.Text
			return &CaptionImgCommand{api: api, msg: replyMsg, originalMsg: msg}
		} else if exReply.Sticker != nil {
			// you can't accompany text w/ a sticker, but you can reply to one
			// which lets you caption a sticker if you reply to it
			return &CaptionStickerCommand{api: api, sticker: *exReply.Sticker, originalMsg: msg}
		}
	} else if reply := msg.ReplyToMessage; reply != nil {
		if reply.Photo != nil {
			replyMsg := msg.ReplyToMessage
			replyMsg.Caption = msg.Text
			return &CaptionImgCommand{api: api, msg: *replyMsg, originalMsg: msg}
		} else if reply.Sticker != nil {
			return &CaptionStickerCommand{api: api, sticker: *reply.Sticker, originalMsg: msg}
		}
	}
	// TODO: make a command response dedicated to failed commands
	return &InvalidCommand{chatID: msg.Chat.ID, request: "Unable to caption reply"}
}
