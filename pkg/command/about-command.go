package command

import (
	telebot "github.com/OvyFlash/telegram-bot-api"
)


type AboutCommand struct {
	chatID int64
	sendConfig telebot.Chattable
}

func (ac *AboutCommand) GenerateMessage() {
	aboutMsg := "GoferBot is a telegram chat utility bot!\n\nFor help, message \"/help\"" +
	"\nTo visit our repo page, go to https://github.com/aphidianxyz/goferbot"
	msgConfig := telebot.NewMessage(ac.chatID, aboutMsg)
	ac.sendConfig = msgConfig
}

func (ac AboutCommand) SendMessage(api *telebot.BotAPI) error {
    if _, err := api.Send(ac.sendConfig); err != nil {
        return err
    }
	return nil
}
