package main

import (
	"log"
	"os"

	telebot "github.com/go-telegram-bot-api/telegram-bot-api"
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
}

func (g *Gofer) Update(timeout int) {
    updateConfig := telebot.NewUpdate(0)
    updateConfig.Timeout = timeout 

    updates, err := g.api.GetUpdatesChan(updateConfig)
    if err != nil {
        log.Panic("Failed to get updates: " + err.Error())
    }

    for update := range updates {
        msg := update.Message
        edit := update.EditedMessage // edit is nil when msg isn't and vice-versa
        if msg == nil && edit == nil {
            continue
        } else if msg.IsCommand() {
            command := ParseMsgCommand(msg)
            if err := command.GenerateMessage(); err != nil {
                errorMessage := telebot.NewMessage(msg.Chat.ID, "Error: " + err.Error())
                g.api.Send(errorMessage)
                continue
            }
            if err := command.SendMessage(g.api); err != nil {
                errorMessage := telebot.NewMessage(msg.Chat.ID, "Error: " + err.Error())
                g.api.Send(errorMessage)
                continue
            }
        } 
    }
}

// todo: main loop has too many responsibilities rn
func main() {
    gofer := Gofer{}
    gofer.Initialize()
    gofer.Update(60)
}

