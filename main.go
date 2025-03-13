package main

import (
	"log"
	"os"

	telebot "github.com/OvyFlash/telegram-bot-api"
	cmd "github.com/aphidian.xyz/bettergrambot/internals/command"
)

type Gofer struct {
    api *telebot.BotAPI
}

func (g *Gofer) initialize() {
    token := os.Getenv("TOKEN")
    bot, err := telebot.NewBotAPI(token)
    if err != nil {
        log.Panic("Failed to initialize bot: " + err.Error())
    }
    g.api = bot
    log.Println("Bot initialized! Account: " + g.api.Self.UserName)
}

func (g *Gofer) update(timeout int) {
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
            }
            if err := command.SendMessage(g.api); err != nil {
                sendError(msg.Chat.ID, err.Error(), g.api)
                continue
            }
        } else if msg.Photo != nil { // text accompanied with a picture is considered a caption
                                     // manual parsing required for commands w/ a picture
        }
    }
}

func handleEdits(update *telebot.Update) {
    if update.EditedMessage == nil {
        return
    }
}

func sendError(chatID int64, errStr string, api *telebot.BotAPI) {
    errSuffix := "Error: "
    errorMessage := telebot.NewMessage(chatID, errSuffix + errStr)
    api.Send(errorMessage)
}

func main() {
    gofer := Gofer{}
    gofer.initialize()
    gofer.update(60)
}

