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
	var imgCmd bool = false
	if msg.Photo != nil {
		imgCmd = true
		msgStr = strings.Split(msg.Caption, "\n")[0]
	}
	cmdTokens := strings.Split(msgStr, " ")
	cmdName := cmdTokens[0]
	cmdParams := cmdTokens[1:]

	switch cmdName {
	case "/about":
		return &AboutCommand{chatID: msg.Chat.ID}
	case "/caption":
		if imgCmd {
			return &CaptionImgCommand{api: cf.api, msg: *msg}
		}
		url, err := getUrl(msgStr)
		if err != nil { // possible reply to an image/sticker
			return cf.createReplyCaptionCommand(msg)
		}
		return &CaptionCommand{msg: *msg, url: url}
	case "/everyone":
		return &EveryoneCommand{msg: *msg, db: cf.chatDB}
	case "/hello":
		return &HelloCommand{chatID: msg.Chat.ID, firstName: msg.From.FirstName, lastName: msg.From.LastName, userName: msg.From.UserName}
	case "/help":
		helpRequest := getHelpRequest(cmdParams)
		return &HelpCommand{chatID: msg.Chat.ID, request: helpRequest, commandJSON: cf.commandJSON}
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
			return &CaptionImgCommand{api: cf.api, msg: replyMsg, originalMsg: msg}
		} else if exReply.Sticker != nil {
			// you can't accompany text w/ a sticker, but you can reply to one
			// which lets you caption a sticker if you reply to it
			return &CaptionStickerCommand{api: cf.api, sticker: *exReply.Sticker, originalMsg: msg}
		}
	} else if reply := msg.ReplyToMessage; reply != nil {
		if reply.Photo != nil {
			replyMsg := msg.ReplyToMessage
			replyMsg.Caption = msg.Text
			return &CaptionImgCommand{api: cf.api, msg: *replyMsg, originalMsg: msg}
		} else if reply.Sticker != nil {
			return &CaptionStickerCommand{api: cf.api, sticker: *reply.Sticker, originalMsg: msg}
		}
	}
	// TODO: make a command response dedicated to failed commands
	return &ErrorCommand{msg: *msg, originCommand: "/caption", error: "Please attach, link or reply to an image to caption it."}
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
