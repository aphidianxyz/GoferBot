package command

import (
	"errors"
	"fmt"
	"hash/fnv"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"

	telebot "github.com/OvyFlash/telegram-bot-api"
	im "gopkg.in/gographics/imagick.v3/imagick"
)

type CaptionImgCommand struct {
    msg telebot.Message
    api *telebot.BotAPI // required to get a file link from an image uploaded on telegram
    imgFilePath string
    sendConfig telebot.Chattable
}

func (ci *CaptionImgCommand) GenerateMessage() {
    // get image from message 
    imgFileID := getLargestPhotoID(ci.msg.Photo)
    imgFileURL, err := ci.api.GetFileDirectURL(imgFileID)
    if err != nil {
        ci.sendConfig = telebot.NewMessage(ci.msg.Chat.ID, err.Error()) 
        return
    }
    ci.imgFilePath, err = downloadImage(imgFileURL)
    if err != nil {
        ci.sendConfig = telebot.NewMessage(ci.msg.Chat.ID, err.Error())
        return
    }
    // get captions
    topCapStr, botCapStr, err := parseCaptions(ci.msg.Caption)
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

func (ci *CaptionImgCommand) SendMessage(api *telebot.BotAPI) error {
    if _, err := api.Send(ci.sendConfig); err != nil {
        if ci.imgFilePath != "" {
            os.Remove(ci.imgFilePath)
        }
        return err
    }
    os.Remove(ci.imgFilePath)
    // remove the original request if successful, to declutter the chat
    return deleteOriginalMessage(ci.msg, api)
}

func deleteOriginalMessage(original telebot.Message, api *telebot.BotAPI) error {
    deleteConfig := telebot.NewDeleteMessage(original.Chat.ID, original.MessageID)
    if _, err := api.Request(deleteConfig); err != nil {
        return err
    }
    return nil
}

func getLargestPhotoID(photoSizes []telebot.PhotoSize) string {
    largest := photoSizes[0]
    for i, photoSize := range photoSizes {
        if i == 0 {
            continue
        }
        if photoSize.Width > largest.Width || photoSize.Height > largest.Height {
            largest = photoSize
        }
    }
    return largest.FileID
}

func captionImage(filepath, topCap, botCap string) error {
    im.Initialize()
    defer im.Terminate()
    mWand := im.NewMagickWand()
    defer mWand.Destroy()
    // load bg img
    if err := mWand.ReadImage(filepath); err != nil {
        return err
    }

    // draw captions and overlay them on bg img
    // TODO: maybe handle different size configs
    topCaptionWand, err := drawCaption(mWand.GetImageWidth(), mWand.GetImageHeight()/4, topCap, true)
    defer topCaptionWand.Destroy()
    if err != nil {
        return errors.New("Failed to draw top caption: " + err.Error())
    }
    mWand.CompositeImageGravity(topCaptionWand, im.COMPOSITE_OP_OVER, im.GRAVITY_NORTH)
    botCaptionWand, err := drawCaption(mWand.GetImageWidth(), mWand.GetImageHeight()/4, botCap, false)
    defer botCaptionWand.Destroy()
    if err != nil {
        return errors.New("Failed to draw bot caption: " + err.Error())
    }
    mWand.CompositeImageGravity(botCaptionWand, im.COMPOSITE_OP_OVER, im.GRAVITY_SOUTH)

    // write to disk
    if err := mWand.WriteImage(filepath); err != nil {
        return err
    }

    return nil
}

func drawCaption(width, height uint, text string, top bool) (*im.MagickWand, error) {
    wand := im.NewMagickWand()
    wand.SetSize(width, height)
    wand.SetFont("./assets/anton/Anton-Regular.ttf")
    wand.SetOption("stroke", "black")
    wand.SetOption("strokewidth", "2")
    wand.SetOption("fill", "white")
    wand.SetOption("background", "none")
    var gravity im.GravityType = im.GRAVITY_NORTH
    if !top{
        gravity = im.GRAVITY_SOUTH
    }
    wand.SetGravity(gravity)

    if err := wand.ReadImage("caption:" + text); err != nil {
        return nil,  err
    }

    return wand, nil
}

func downloadImage(url string) (filepath string, error error) {
    // get image
    // setting headers and spoofing user-agent to avoid a 403
    request, err := http.NewRequest("GET", url, nil)
    request.Header.Set("Accept", "image/avif,image/webp,image/apng,image/svg+xml,image/*,*/*;q=0.8")
    request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
    client := &http.Client{}
    response, err := client.Do(request)
    if err != nil {
        return "", err
    }
    defer response.Body.Close()

    if response.StatusCode != http.StatusOK {
        return "", errors.New("Failed to download image, got HTTP response: " + strconv.Itoa(response.StatusCode))
    }

    filepath = genUniqueFileName()

    // write to disk
    file, err := os.Create(filepath)
    if err != nil {
        return "", err
    }
    defer file.Close()

    _, err = io.Copy(file, response.Body)
    if err != nil {
        return "", err
    }

    return filepath, nil
}

func parseCaptions(prompt string) (topCaption, botCaption string, error error) {
    regex := regexp.MustCompile(`"([^"]*[a-zA-Z\s\\"]*)"\s+"([^"]*[a-zA-Z\s\\"]*)"$`)
    captions := regex.FindStringSubmatch(prompt)
    if len(captions) != 3 { // the first element is the match w/o groups
        return "", "", errors.New("Expected 2 captions, each encapsulated in quotations")
    }
    return captions[1], captions[2], nil
}

func genUniqueFileName() string {
    hash := fnv.New32a()
    hash.Write([]byte(time.Now().Truncate(time.Nanosecond).String())) // Simplified time handling
    filename := fmt.Sprintf("./tmp/%x.png", hash.Sum32()) // Direct hex format, handles error internally
    return filename
}

