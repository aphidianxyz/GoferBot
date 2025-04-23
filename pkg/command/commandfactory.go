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
	// cmds and their param(s) should only appear on the first line
	var msgStr string = strings.Split(msg.Text, "\n")[0]
	if msg.Photo != nil || msg.Video != nil || msg.Animation != nil ||
	msg.Document != nil || msg.Voice != nil {
		msgStr = strings.Split(msg.Caption, "\n")[0]
	}
	cmdTokens := strings.Split(msgStr, " ")
	cmdName := cmdTokens[0]
	cmdParams := cmdTokens[1:]

	switch cmdName {
	case "/about":
		return MakeAbout(msg.Chat.ID)
	case "/caption":
		if msg.Photo != nil {
			return MakeCaptionURL(cf.api, *msg, nil)
		}
		url, err := getUrl(msgStr)
		if err != nil {
			return cf.createReplyCaption(msg)
		}
		return MakeCaption(*msg, url)
	case "/everyone":
		return MakeEveryone(*msg, cf.chatDB)
	case "/hello":
		return MakeHello(*msg)
	case "/help":
		helpRequest := getHelpRequest(cmdParams)
		return MakeHelp(*msg, helpRequest, cf.commandJSON)
	case "/pin":
		return MakePin(*msg)
	case "/ping":
		return MakePing(*msg)
	default:
		return MakeError(*msg, cmdName, "is not a valid command!\nCall /help for a list of commands!")
	}
}

func (cf CommandFactory) createReplyCaption(msg *telebot.Message) Command {
	if exReply := msg.ExternalReply; exReply != nil {
		if exReply.Photo != nil {
			replyMsg := *msg
			replyMsg.Photo = msg.ExternalReply.Photo
			replyMsg.Caption = msg.Text
			return MakeCaptionURL(cf.api, replyMsg, msg)
		} else if exReply.Sticker != nil {
			// you can't accompany text w/ a sticker, but you can reply to one
			// which lets you caption a sticker if you reply to it
			return MakeCaptionSticker(cf.api, *exReply.Sticker, msg)
		}
	} else if reply := msg.ReplyToMessage; reply != nil {
		if reply.Photo != nil {
			replyMsg := msg.ReplyToMessage
			replyMsg.MessageID = msg.MessageID
			replyMsg.Caption = msg.Text
			return MakeCaptionURL(cf.api, *replyMsg, msg)
		} else if reply.Sticker != nil {
			return MakeCaptionSticker(cf.api, *reply.Sticker, msg)
		}
	}
	return MakeError(*msg, "/caption", "Please attach, link or reply to an image to caption it")
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
