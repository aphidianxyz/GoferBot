package gofer

import (
	"log"
	"os"
	"strings"

	telebot "github.com/OvyFlash/telegram-bot-api"
	cmd "github.com/aphidianxyz/GoferBot/pkg/command"
)

type Gofer struct {
    api *telebot.BotAPI
}

func (g *Gofer) Initialize() {
    token := os.Getenv("TOKEN")
    bot, err := telebot.NewBotAPI(token)
    if err != nil {
        log.Panic("Failed to initialize bot: " + err.Error())
    }
    g.api = bot
    log.Println("Bot initialized! Account: " + g.api.Self.UserName)
}

func (g *Gofer) Update(timeout int) {
    updateConfig := telebot.NewUpdate(0)
    updateConfig.Timeout = timeout 

    updates := g.api.GetUpdatesChan(updateConfig)

    for update := range updates {
        msg := update.Message
        edit := update.EditedMessage // edit is nil when msg isn't and vice-versa
        if msg == nil {
            if edit != nil {
                handleEdits(&update)
            }
            continue
        } else if msg.IsCommand() {
            command := cmd.ParseMsgCommand(msg)
            if err := command.GenerateMessage(); err != nil {
                sendError(msg.Chat.ID, err.Error(), g.api)
                continue
            } // TODO: this impl currently doesn't support multi-step commands
            if err := command.SendMessage(g.api); err != nil {
                sendError(msg.Chat.ID, err.Error(), g.api)
                continue
            }
        } else if msg.Photo != nil { // msg w/ photos have captions, manual parsing required
            if !isCaptionCommand(msg.Caption) {
                continue
            }
            command := cmd.ParseImgCommand(msg)
            if err := command.GenerateMessage(); err != nil {
                sendError(msg.Chat.ID, err.Error(), g.api)
                continue
            }
            if err := command.SendMessage(g.api); err != nil {
                sendError(msg.Chat.ID, err.Error(), g.api)
                continue
            }
        } else {
            // TODO: handle registered responses
        }
    }
}

func isCaptionCommand(caption string) bool {
    tokens := strings.Split(caption, " ")
    if len(tokens) == 0 {
        return false
    }
    if commandName := tokens[0]; commandName[0] != '/' {
        return true 
    }
    return false
}

func handleEdits(update *telebot.Update) {
    if update.EditedMessage == nil {
        return
    }
    // TODO: add operations if we want to handle certain edit events
}

func sendError(chatID int64, errStr string, api *telebot.BotAPI) {
    errSuffix := "Error: "
    errorMessage := telebot.NewMessage(chatID, errSuffix + errStr)
    api.Send(errorMessage)
}
