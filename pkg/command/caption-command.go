package command

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"

	telebot "github.com/OvyFlash/telegram-bot-api"
)

type CaptionCommand struct {
    msg telebot.Message
	url string
    imgFilePath string
    sendConfig telebot.Chattable
}

func (ci *CaptionCommand) GenerateMessage() {
	imgFilePath, err := downloadImage(ci.url)
    if err != nil {
        ci.sendConfig = telebot.NewMessage(ci.msg.Chat.ID, err.Error())
        return
    }
	ci.imgFilePath = imgFilePath
    // get captions
    topCapStr, botCapStr, err := parseCaptions(ci.msg.Text)
    if err != nil {
        ci.sendConfig = telebot.NewMessage(ci.msg.Chat.ID, err.Error())
        return
    }
    // generate image
    if err := captionImage(ci.imgFilePath, topCapStr, botCapStr); err != nil {
        ci.sendConfig = telebot.NewMessage(ci.msg.Chat.ID, err.Error())
        return
    }
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
    return deleteOriginalMessage(ci.msg, api)

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
