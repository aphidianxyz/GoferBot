package command

import (
	"database/sql"
	"strings"

	telebot "github.com/OvyFlash/telegram-bot-api"
)

type CommandFactory struct {
	api *telebot.BotAPI 
	chatDB *sql.DB
	commandJSON CommandJSON
}

func ConstructCommandFactory(api *telebot.BotAPI, chatDB *sql.DB, commandJSON CommandJSON) CommandFactory {
	return CommandFactory{api: api, chatDB: chatDB, commandJSON: commandJSON}
}

func (cf CommandFactory) CreateCommand(update *telebot.Update) Command {
	// as of now, all commands only need messages
	msg := update.Message
	// cmds should only appear on the first line
	var msgStr string = strings.Split(msg.Text, "\n")[0]
	// TODO: trying to call a cmd with video or other attachments fail 
	if msg.Photo != nil || msg.Video != nil || msg.Animation != nil ||
	msg.Document != nil || msg.Voice != nil {
		msgStr = strings.Split(msg.Caption, "\n")[0]
	}
	cmdTokens := strings.Split(msgStr, " ")
	cmdName := cmdTokens[0]
	cmdParams := cmdTokens[1:]

	switch cmdName {
	case "/about":
		return MakeAboutCommand(msg.Chat.ID)
	case "/caption":
		if msg.Photo != nil {
			return MakeCaptionImgCommand(cf.api, *msg, nil)
		}
		url, err := getUrl(msgStr)
		if err != nil { // possible reply to an image/sticker
			return cf.createReplyCaptionCommand(msg)
		}
		return MakeCaptionCommand(*msg, url)
	case "/everyone":
		return MakeEveryoneCommand(*msg, cf.chatDB)
	case "/hello":
		return MakeHelloCommand(*msg)
	case "/help":
		helpRequest := getHelpRequest(cmdParams)
		return MakeHelpCommand(*msg, helpRequest, cf.commandJSON)
	case "/pin":
		return &PinCommand{api: cf.api, msg: *msg}
	case "/ping":
		return &PingCommand{chatID: msg.Chat.ID}
	default:
		return &InvalidCommand{msg: *msg, request: cmdName}
	}
}

func (cf CommandFactory) createReplyCaptionCommand(msg *telebot.Message) Command {
	if exReply := msg.ExternalReply; exReply != nil {
		if exReply.Photo != nil {
			replyMsg := *msg
			replyMsg.Photo = msg.ExternalReply.Photo
			replyMsg.Caption = msg.Text
			return MakeCaptionImgCommand(cf.api, replyMsg, msg)
		} else if exReply.Sticker != nil {
			// you can't accompany text w/ a sticker, but you can reply to one
			// which lets you caption a sticker if you reply to it
			return MakeCaptionStickerCommand(cf.api, *exReply.Sticker, msg)
		}
	} else if reply := msg.ReplyToMessage; reply != nil {
		if reply.Photo != nil {
			replyMsg := msg.ReplyToMessage
			replyMsg.MessageID = msg.MessageID
			replyMsg.Caption = msg.Text
			return MakeCaptionImgCommand(cf.api, *replyMsg, msg)
		} else if reply.Sticker != nil {
			return MakeCaptionStickerCommand(cf.api, *reply.Sticker, msg)
		}
	}
	return MakeErrorCommand(*msg, "/caption", "Please attach, link or reply to an image to caption it")
}

// an empty string request will signal a help command to generate
// a complete list of commands and syntaxes
func getHelpRequest(cmdParams []string) (helpRequest string) {
	if len(cmdParams) < 1 {
		helpRequest = ""
	} else {
		helpRequest = cmdParams[0]
	}
	return helpRequest
}
