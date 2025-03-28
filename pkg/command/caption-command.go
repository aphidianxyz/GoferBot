package command

import (
	"errors"
	"net/url"
	"os"
	"strings"

	telebot "github.com/OvyFlash/telegram-bot-api"
)

type CaptionCommand struct {
    chatID int64
    msg telebot.Message
    imgFilePath string
    sendConfig telebot.Chattable
}

func (ci *CaptionCommand) GenerateMessage() {
    url, err := getUrl(ci.msg.Text)
    if err != nil {
        ci.sendConfig = telebot.NewMessage(ci.msg.Chat.ID, err.Error())
        return
    }
    ci.imgFilePath, err = downloadImage(url)
    if err != nil {
        ci.sendConfig = telebot.NewMessage(ci.msg.Chat.ID, err.Error())
        return
    }
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
    deleteOriginalMessage(ci.msg, api)
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
        return "./tmp/", errors.New("Invalid or malformed URL was inputted")
    }
    return tokens[1], nil
}
