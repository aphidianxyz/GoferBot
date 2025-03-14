package command

import (
	"errors"
	"fmt"
	"hash/fnv"
	"io"
	"net/http"
	"os"
	"regexp"

    im "gopkg.in/gographics/imagick.v3/imagick"
	telebot "github.com/OvyFlash/telegram-bot-api"
)

type CaptionImgCommand struct {
    chatID int64
    msg telebot.Message
    api *telebot.BotAPI
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
        return errors.New("Failed to send a CaptionImgCommand")
    }
    os.Remove(ci.imgFilePath)
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

// TODO: maybe allow these to be public fns, since CaptionCommand will need them
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
    topCaptionWand, err := drawCaption(mWand.GetImageWidth(), mWand.GetImageHeight()/3, topCap, true)
    defer topCaptionWand.Destroy()
    if err != nil {
        return errors.New("Failed to draw top caption: " + err.Error())
    }
    mWand.CompositeImageGravity(topCaptionWand, im.COMPOSITE_OP_OVER, im.GRAVITY_NORTH)
    botCaptionWand, err := drawCaption(mWand.GetImageWidth(), mWand.GetImageHeight()/3, botCap, false)
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
    wand.SetFont("Impact")
    wand.SetOption("stroke", "black")
    wand.SetOption("strokewidth", "1")
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
    response, err := http.Get(url)
    if err != nil {
        return "", err
    }
    defer response.Body.Close()

    filepath = genUniqueFileName()

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

func genUniqueFileName() string {
    hash := fnv.New32a()
    tempFilenameSuffix := hash.Sum32()
    filename := "temp_caption_" + fmt.Sprint(tempFilenameSuffix) + ".jpg"
    return filename
}

func parseCaptions(prompt string) (topCaption, botCaption string, error error) {
    regex := regexp.MustCompile(`^/caption\s+"([^"]*[a-zA-Z\s\\"]*)"\s+"([^"]*[a-zA-Z\s\\"]*)"$`)
    captions := regex.FindStringSubmatch(prompt)
    if len(captions) != 3 {
        return "", "", errors.New("Expected 2 captions, each encapsulated in quotations")
    }
    return captions[1], captions[2], nil
}
