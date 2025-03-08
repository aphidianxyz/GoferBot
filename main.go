package main

import (
	"errors"
	"fmt"
	"hash/fnv"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	telebot "github.com/go-telegram-bot-api/telegram-bot-api"
	magick "gopkg.in/gographics/imagick.v3/imagick"
)

// todo: main loop has too many responsibilities rn
func main() {
    token := os.Getenv("TOKEN")
    bot, err := telebot.NewBotAPI(token)

    if err != nil {
        log.Panic("Failed to initialize bot: " + err.Error())
    }

    log.Printf("Authorized on account %s", bot.Self.UserName)

    updateConfig := telebot.NewUpdate(0)
    updateConfig.Timeout = 60

    updates, err := bot.GetUpdatesChan(updateConfig)
    if err != nil {
        log.Panic("Failed to get updates: " + err.Error())
    }

    for update := range updates {
        msg := update.Message 
        editedMsg := update.EditedMessage
        log.Printf("%+v\n", msg)

        if msg == nil {
            if editedMsg == nil {
                log.Println("Invalid message was sent")
            }
            log.Println("Message was edited (handle it)")
            continue
        } else if msg.IsCommand() { // commands without a picture
            tokens := strings.Split(msg.Text, " ")
            commandName := tokens[0]
            //commandParams := tokens[1:]
            var msgConfig telebot.MessageConfig 

            switch commandName {
            case "/hello":
                helloString := "Hello, " + msg.Chat.FirstName + "!"
                msgConfig = telebot.NewMessage(msg.Chat.ID, helloString)
            /*
            // TODO: write /caption for URLs
            case "/caption": 
                if len(commandParams) < 1 {
                    helpString := `Correct usage: /caption [Image URL] ["Top Text"] ["Bottom Text"] (Brackets not included)`
                    msgConfig = telebot.NewMessage(msg.Chat.ID, helpString)
                    break
                }
                if msg.Photo == nil || msg.ReplyToMessage.Photo == nil {
                    noImageError := "No image was given to caption!"
                    msgConfig = telebot.NewMessage(msg.Chat.ID, noImageError)
                    break
                }
                photos := *msg.Photo
                // TODO: handle operations on whole photo albums
                targetFileID := photos[0].FileID
                targetFileURL, err := bot.GetFileDirectURL(targetFileID)
                if err != nil {
                    failedImgDownloadError := "Failed to get URL to download image - " + err.Error()
                    msgConfig = telebot.NewMessage(msg.Chat.ID, failedImgDownloadError)
                    break
                }
                img, err := DownloadImage(targetFileURL)
                if err != nil {
                    failedDownloadError := "Failed to download image - " + err.Error()
                    msgConfig = telebot.NewMessage(msg.Chat.ID, failedDownloadError)
                    break
                }

                ready := fmt.Sprintf("Image downloaded and ready to caption: %v", img)
                msgConfig = telebot.NewMessage(msg.Chat.ID, ready)
            */ 
            default:
                msgConfig = telebot.NewMessage(msg.Chat.ID, "Unknown command")
                msgConfig.ReplyToMessageID = msg.MessageID
            }
            bot.Send(msgConfig)
        } else if msg.Photo != nil { // for commands that have an attached photo; any text attached w/ a photo
                                     // is considered as a caption (so we can't check if it's a cmd)
            var sendConfig telebot.Chattable
            tokens := strings.Split(msg.Caption, " ")
            commandName := tokens[0]
            if len(commandName) == 0 {
                log.Print("no command")
                continue
            }
            if commandName[0] != '/' { // since we can't check if a caption is a command w/ tg API
                log.Print("not a command")
                continue
            }
            commandParams := tokens[1:]
            switch commandName {
            case "/caption": 
                if len(commandParams) < 1 {
                    helpString := `Correct usage: /caption ["Top Text"] ["Bottom Text"] (Brackets not included)`
                    sendConfig = telebot.NewMessage(msg.Chat.ID, helpString)
                    break
                }
                photos := *msg.Photo
                // msg.Photo is a slice of PhotoSizes from the TG API which are offerings of the same
                // photo in various sizes
                // TODO: find biggest photo in PhotoSizes
                targetFileID := photos[1].FileID
                targetFileURL, err := bot.GetFileDirectURL(targetFileID)
                if err != nil {
                    failedImgDownloadError := "Failed to get URL to download image - " + err.Error()
                    sendConfig = telebot.NewMessage(msg.Chat.ID, failedImgDownloadError)
                    break
                }
                imgFilePath, err := DownloadImage(targetFileURL)
                if err != nil {
                    failedDownloadError := "Failed to download image - " + err.Error()
                    sendConfig = telebot.NewMessage(msg.Chat.ID, failedDownloadError)
                    break
                }
                topCaptionText, botCaptionText, err := ParseCaptions(msg.Caption)
                if err != nil {
                    sendConfig = telebot.NewMessage(msg.Chat.ID, err.Error())
                    break
                }
                err = CaptionImage(imgFilePath, topCaptionText, botCaptionText)
                if err != nil {
                    failedToCaptionImg := "Failed to caption image - " + err.Error()
                    sendConfig = telebot.NewMessage(msg.Chat.ID, failedToCaptionImg)
                    break
                }
                photoConfig := telebot.NewPhotoUpload(msg.Chat.ID, imgFilePath)
                sendConfig = photoConfig 
                bot.Send(sendConfig)
                os.Remove(imgFilePath)
                continue
            default:
                msgConfig := telebot.NewMessage(msg.Chat.ID, "Unknown command")
                msgConfig.ReplyToMessageID = msg.MessageID
                sendConfig = msgConfig
            }
            bot.Send(sendConfig)
        } else { // non-commands, but we can generate replies to certain keywords in chat
            if msg.Text == "goat" {
                photoConfig := telebot.NewPhotoUpload(msg.Chat.ID, "./goat.jpg")
                photoConfig.Caption = "Did someone say GOAT???"
                photoConfig.ReplyToMessageID = msg.MessageID

                bot.Send(photoConfig)
            } else {
                log.Println("[", msg.From.UserName, "]", msg.Text)

                msgConfig := telebot.NewMessage(msg.Chat.ID, msg.Text)
                msgConfig.ReplyToMessageID = msg.MessageID

                bot.Send(msgConfig)
            }
        }
    }

}

func DownloadImage(url string) (filepath string, err error) {
    response, err := http.Get(url)
    if err != nil {
        log.Println("Failed to retrieve image from URL: " + url + " error: " + err.Error())
        response.Body.Close()
        return "", err
    }
    defer response.Body.Close()

    // create unique name for temp file
    hash := fnv.New32a()
    hash.Write([]byte(url))
    tempFilenameSuffix := hash.Sum32()
    filepath = "temp_" + fmt.Sprint(tempFilenameSuffix) + ".jpg"
    
    file, err := os.Create(filepath)
    if err != nil {
        log.Println("Failed to download image: " + err.Error())
        return "", err 
    }
    defer file.Close()

    _, err = io.Copy(file, response.Body)
    if err != nil {
        log.Printf("Failed to write image to disk: %e", err)
        return "", err
    }
    
    return filepath, nil
}

func CaptionImage(filepath, topText, botText string) error {
    if topText == "" && botText == "" {
        return errors.New("No captions provided")
    }

    magick.Initialize()
    defer magick.Terminate()
    mWand := magick.NewMagickWand()
    defer mWand.Destroy()
    if err := mWand.ReadImage(filepath); err != nil {
        return errors.New("Imagemagick failed to read image @ " + filepath + ". Error: " + err.Error())
    }

    // top caption
    // captions should be bounded by original image's dimensions
    topCaptionWand, err := DrawCaption(mWand.GetImageWidth(), mWand.GetImageHeight()/3, topText, true)
    defer topCaptionWand.Destroy()
    if err != nil {
        return errors.New("Failed to draw top caption")
    }
    mWand.CompositeImageGravity(topCaptionWand, magick.COMPOSITE_OP_OVER, magick.GRAVITY_NORTH)

    // bot caption
    botCaptionWand, err := DrawCaption(mWand.GetImageWidth(), mWand.GetImageHeight()/3, botText, false)
    defer botCaptionWand.Destroy()
    if err != nil {
        return errors.New("Failed to draw top caption")
    }
    mWand.CompositeImageGravity(botCaptionWand, magick.COMPOSITE_OP_OVER, magick.GRAVITY_SOUTH)

    if err := mWand.WriteImage(filepath); err != nil {
        return errors.New("Failed to write captions to original image: " + err.Error())
    }
    
    return nil
}

func DrawCaption(width, height uint, text string, top bool) (*magick.MagickWand, error) {
    wand := magick.NewMagickWand()
    wand.SetSize(width, height)
    wand.SetFont("Impact")
    wand.SetOption("stroke", "black")
    wand.SetOption("strokewidth", "1")
    wand.SetOption("fill", "white")
    wand.SetOption("background", "none")
    var gravity magick.GravityType = magick.GRAVITY_NORTH
    if !top {
        gravity = magick.GRAVITY_SOUTH
    }
    wand.SetGravity(gravity)

    if err := wand.ReadImage("caption:" + text); err != nil {
        return nil, errors.New("Failed to draw caption: " + text)
    }

    return wand, nil
}

func ParseCaptions(prompt string) (topCaption, botCaption string, error error) {
    regex := regexp.MustCompile(`^/caption\s+"([^"]*[a-zA-Z\s\\"]*)"\s+"([^"]*[a-zA-Z\s\\"]*)"$`)
    captions := regex.FindStringSubmatch(prompt)
    log.Println(captions)
    if len(captions) != 3 {
        return "", "", errors.New("Expected 2 captions, each encapsulated in quotations")
    }
    return captions[1], captions[2], nil

}
