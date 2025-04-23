package command

import (
	telebot "github.com/OvyFlash/telegram-bot-api"
)

type About struct {
	chatID int64
	sendConfig telebot.Chattable
}

func MakeAbout(chatID int64) Command {
	return &About{chatID: chatID}
}

func (a *About) GenerateMessage() {
	aboutMsg := "GoferBot is a telegram chat utility bot!\n\nFor help, message \"/help\"" +
	"\nTo visit our repo page, go to https://github.com/aphidianxyz/goferbot"
	msgConfig := telebot.NewMessage(a.chatID, aboutMsg)
	a.sendConfig = msgConfig
}

func (a About) SendMessage(api *telebot.BotAPI) error {
    if _, err := api.Send(a.sendConfig); err != nil {
        return err
    }
	return nil
}
