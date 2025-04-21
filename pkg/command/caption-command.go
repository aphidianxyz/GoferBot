package command

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	telebot "github.com/OvyFlash/telegram-bot-api"
)

type CaptionCommand struct {
    msg telebot.Message
    imgFilePath string
    sendConfig telebot.Chattable
}

func MakeCaptionCommand(msg telebot.Message, url string) Command {
	imgFilePath, err := downloadImage(url)
    if err != nil {
		return MakeErrorCommand(msg, "/caption", "invalid image attachment: " + err.Error())
    }
    // get captions
    topCapStr, botCapStr, err := parseCaptions(msg.Text)
    if err != nil {
		os.Remove(imgFilePath)
		return MakeErrorCommand(msg, "/caption", "failed to parse captions: " + err.Error())
    }
    // generate image
    if err := captionImage(imgFilePath, topCapStr, botCapStr); err != nil {
		os.Remove(imgFilePath)
		return MakeErrorCommand(msg, "/caption", "failed to draw captions: " + err.Error())
    }
	return &CaptionCommand{msg: msg, imgFilePath: imgFilePath}
}

func MakeCaptionStickerCommand(api *telebot.BotAPI, sticker telebot.Sticker, originalMsg *telebot.Message) Command {
	var imgFilePath string
	imgFileID := sticker.FileID
	imgFileURL, err := api.GetFileDirectURL(imgFileID)
	if err != nil {
		return MakeErrorCommand(*originalMsg, "caption", "could not retrieve sticker from Telegram: " + err.Error())
	}
	imgFilePath, err = downloadImage(imgFileURL)
	if err != nil {
		os.Remove(imgFilePath)
		return MakeErrorCommand(*originalMsg, "caption", "failed to download sticker: " + err.Error())
	}
    // get captions
    topCapStr, botCapStr, err := parseCaptions(originalMsg.Text)
    if err != nil {
		os.Remove(imgFilePath)
		return MakeErrorCommand(*originalMsg, "caption", "could not parse captions: " + err.Error())
    }
    // generate image
    if err := captionImage(imgFilePath, topCapStr, botCapStr); err != nil {
		os.Remove(imgFilePath)
		return MakeErrorCommand(*originalMsg, "caption", "could not draw captions: " + err.Error())
    }
	return &CaptionCommand{imgFilePath: imgFilePath, msg: *originalMsg}
}

func (ci *CaptionCommand) GenerateMessage() {
    // generate message
    image := telebot.FilePath(ci.imgFilePath)
    image.UploadData()
    photoConfig := telebot.NewPhoto(ci.msg.Chat.ID, image)
	photoConfig.Caption = "Here's your meme!\n" + fmt.Sprintf("[%v](tg://user?id=%v)", ci.msg.From.FirstName, ci.msg.From.ID)
	photoConfig.ParseMode = "MarkDown"
    ci.sendConfig = photoConfig
}

func (ci *CaptionCommand) SendMessage(api *telebot.BotAPI) error {
	if _, err := api.Send(ci.sendConfig); err != nil {
        if ci.imgFilePath != "" {
            os.Remove(ci.imgFilePath)
        }
        return err
    }
    os.Remove(ci.imgFilePath)
    deleteMessage(api, 3 * time.Second, ci.msg.Chat.ID, ci.msg.MessageID)
	return nil

}

// urls are expected as the first param of /caption
// i.e. /caption [url] ["top"] ["bot"]
func getUrl(prompt string) (string, error) {
    tokens := strings.Split(prompt, " ")
    if len(tokens) < 2 {
        return "", errors.New("No URL provided")
    }
    if _, err := url.ParseRequestURI(tokens[1]); err != nil {
        return "", errors.New("Invalid or malformed URL was inputted")
    }
    return tokens[1], nil
}
