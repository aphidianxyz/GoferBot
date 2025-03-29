package command

import (
	"fmt"
	"os"

	telebot "github.com/OvyFlash/telegram-bot-api"
)

type CaptionStickerCommand struct {
	api *telebot.BotAPI
	sticker telebot.Sticker
	originalMsg *telebot.Message
	imgFilePath string
	sendConfig telebot.Chattable
}

func (cs *CaptionStickerCommand) GenerateMessage() {
	imgFileID := cs.sticker.FileID
	imgFileURL, err := cs.api.GetFileDirectURL(imgFileID)
	if err != nil {
		cs.sendConfig = telebot.NewMessage(cs.originalMsg.Chat.ID, err.Error())
		return
	}
	cs.imgFilePath, err = downloadImage(imgFileURL)
	if err != nil {
		cs.sendConfig = telebot.NewMessage(cs.originalMsg.From.ID, err.Error())
		return
	}
    // get captions
    topCapStr, botCapStr, err := parseCaptions(cs.originalMsg.Text)
    if err != nil {
        cs.sendConfig = telebot.NewMessage(cs.originalMsg.Chat.ID, err.Error())
        return
    }
    // generate image
    if err := captionImage(cs.imgFilePath, topCapStr, botCapStr); err != nil {
        cs.sendConfig = telebot.NewMessage(cs.originalMsg.Chat.ID, err.Error())
        return
    }
    // generate message
    image := telebot.FilePath(cs.imgFilePath)
    image.UploadData()
    photoConfig := telebot.NewPhoto(cs.originalMsg.Chat.ID, image)
	recipientName := cs.originalMsg.From.FirstName
	recipientID := cs.originalMsg.From.ID
	photoConfig.Caption = "Here's your meme!\n" + fmt.Sprintf("[%v](tg://user?id=%v)", recipientName, recipientID)
	photoConfig.ParseMode = "MarkDown"
	cs.sendConfig = photoConfig
}

func (cs *CaptionStickerCommand) SendMessage(api *telebot.BotAPI) error {
    if _, err := api.Send(cs.sendConfig); err != nil {
        if cs.imgFilePath != "" {
            os.Remove(cs.imgFilePath)
        }
        return err
    }
    os.Remove(cs.imgFilePath)
    // remove the original request if successful, to declutter the chat
	return deleteOriginalMessage(*cs.originalMsg, api)
}
