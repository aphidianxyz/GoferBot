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

func (csc *CaptionStickerCommand) GenerateMessage() {
	imgFileID := csc.sticker.FileID
	imgFileURL, err := csc.api.GetFileDirectURL(imgFileID)
	if err != nil {
		csc.sendConfig = telebot.NewMessage(csc.originalMsg.Chat.ID, err.Error())
		return
	}
	csc.imgFilePath, err = downloadImage(imgFileURL)
	if err != nil {
		csc.sendConfig = telebot.NewMessage(csc.originalMsg.From.ID, err.Error())
		return
	}
    // get captions
    topCapStr, botCapStr, err := parseCaptions(csc.originalMsg.Text)
    if err != nil {
        csc.sendConfig = telebot.NewMessage(csc.originalMsg.Chat.ID, err.Error())
        return
    }
    // generate image
    if err := captionImage(csc.imgFilePath, topCapStr, botCapStr); err != nil {
        csc.sendConfig = telebot.NewMessage(csc.originalMsg.Chat.ID, err.Error())
        return
    }
    // generate message
    image := telebot.FilePath(csc.imgFilePath)
    image.UploadData()
    photoConfig := telebot.NewPhoto(csc.originalMsg.Chat.ID, image)
	recipientName := csc.originalMsg.From.FirstName
	recipientID := csc.originalMsg.From.ID
	photoConfig.Caption = "Here's your meme!\n" + fmt.Sprintf("[%v](tg://user?id=%v)", recipientName, recipientID)
	photoConfig.ParseMode = "MarkDown"
	csc.sendConfig = photoConfig
}

func (csc *CaptionStickerCommand) SendMessage(api *telebot.BotAPI) error {
    if _, err := api.Send(csc.sendConfig); err != nil {
        if csc.imgFilePath != "" {
            os.Remove(csc.imgFilePath)
        }
        return err
    }
    os.Remove(csc.imgFilePath)
    // remove the original request if successful, to declutter the chat
	return deleteOriginalMessage(*csc.originalMsg, api)
}
