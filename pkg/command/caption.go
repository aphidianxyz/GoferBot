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

type Caption struct {
    msg telebot.Message
    imgFilePath string
    sendConfig telebot.Chattable
}

func MakeCaption(msg telebot.Message, url string) Command {
	imgFilePath, err := downloadImage(url)
    if err != nil {
		return MakeError(msg, "/caption", "invalid image attachment: " + err.Error())
    }
    // get captions
    topCapStr, botCapStr, err := parseCaptions(msg.Text)
    if err != nil {
		os.Remove(imgFilePath)
		return MakeError(msg, "/caption", "failed to parse captions: " + err.Error())
    }
    // generate image
    if err := captionImage(imgFilePath, topCapStr, botCapStr); err != nil {
		os.Remove(imgFilePath)
		return MakeError(msg, "/caption", "failed to draw captions: " + err.Error())
    }
	return &Caption{msg: msg, imgFilePath: imgFilePath}
}

func MakeCaptionSticker(api *telebot.BotAPI, sticker telebot.Sticker, originalMsg *telebot.Message) Command {
	var imgFilePath string
	imgFileID := sticker.FileID
	imgFileURL, err := api.GetFileDirectURL(imgFileID)
	if err != nil {
		return MakeError(*originalMsg, "caption", "could not retrieve sticker from Telegram: " + err.Error())
	}
	imgFilePath, err = downloadImage(imgFileURL)
	if err != nil {
		os.Remove(imgFilePath)
		return MakeError(*originalMsg, "caption", "failed to download sticker: " + err.Error())
	}
    // get captions
    topCapStr, botCapStr, err := parseCaptions(originalMsg.Text)
    if err != nil {
		os.Remove(imgFilePath)
		return MakeError(*originalMsg, "caption", "could not parse captions: " + err.Error())
    }
    // generate image
    if err := captionImage(imgFilePath, topCapStr, botCapStr); err != nil {
		os.Remove(imgFilePath)
		return MakeError(*originalMsg, "caption", "could not draw captions: " + err.Error())
    }
	return &Caption{imgFilePath: imgFilePath, msg: *originalMsg}
}

func (c *Caption) GenerateMessage() {
    // generate message
    image := telebot.FilePath(c.imgFilePath)
    image.UploadData()
    photoConfig := telebot.NewPhoto(c.msg.Chat.ID, image)
	photoConfig.Caption = "Here's your meme!\n" + fmt.Sprintf("[%v](tg://user?id=%v)", c.msg.From.FirstName, c.msg.From.ID)
	photoConfig.ParseMode = "MarkDown"
    c.sendConfig = photoConfig
}

func (c *Caption) SendMessage(api *telebot.BotAPI) error {
	if _, err := api.Send(c.sendConfig); err != nil {
        if c.imgFilePath != "" {
            os.Remove(c.imgFilePath)
        }
        return err
    }
    os.Remove(c.imgFilePath)
    deleteMessage(api, 3 * time.Second, c.msg.Chat.ID, c.msg.MessageID)
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
