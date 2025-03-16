package gofer

import (
	"log"
	"os"
	"strings"
    "database/sql"

    //sqlite "github.com/mattn/go-sqlite3"
	telebot "github.com/OvyFlash/telegram-bot-api"
	cmd "github.com/aphidianxyz/GoferBot/pkg/command"
)


type Gofer struct {
    api *telebot.BotAPI
    db *sql.DB
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
                g.handleEdits(&update)
            }
            continue
        }
        if msg.IsCommand() {
            g.handleCommands(&update)
        } else if msg.Photo != nil { // msg w/ photos have captions, manual parsing required
            g.handlePhotoCommands(&update)
        } else { // TODO: handle messages/command requests with a video or gif attached
            // TODO: handle registered responses
        }
    }
}

func (g *Gofer) handleCommands(update *telebot.Update) {
    msg := update.Message
    command := cmd.ParseMsgCommand(g.api, msg)
    // TODO: this impl currently doesn't support multi-step commands
    command.GenerateMessage()
    if err := command.SendMessage(g.api); err != nil {
        sendError(msg.Chat.ID, err.Error(), g.api)
        return
    }
}

func (g *Gofer) handlePhotoCommands(update *telebot.Update) {
    msg := update.Message
    if !isCaptionCommand(msg.Caption) {
        return
    }
    command := cmd.ParseImgCommand(g.api, msg)
    command.GenerateMessage()
    if err := command.SendMessage(g.api); err != nil {
        sendError(msg.Chat.ID, err.Error(), g.api)
        return
    }
}

func (g *Gofer) handleEdits(update *telebot.Update) {
    if update.EditedMessage == nil {
        return
    }
    // TODO: add operations if we want to handle certain edit events
}

func isCaptionCommand(caption string) bool {
    tokens := strings.Split(caption, " ")
    if len(tokens) == 0 {
        return false
    }
    if commandName := tokens[0]; commandName[0] == '/' {
        return true 
    }
    return false
}


func sendError(chatID int64, errStr string, api *telebot.BotAPI) {
    errSuffix := "Error: "
    errorMessage := telebot.NewMessage(chatID, errSuffix + errStr)
    api.Send(errorMessage)
}
