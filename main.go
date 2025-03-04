package main

import (
	"fmt"
	"hash/fnv"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	telebot "github.com/go-telegram-bot-api/telegram-bot-api"
	//magick "gopkg.in/gographics/imagick.v3/imagick"
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
        } else if msg.Photo != nil { // for commands that have an attached photo, since any text attached w/ a photo
                                     // is not considered as text (so we can't check if it's a cmd), but a caption
            var msgConfig telebot.MessageConfig
            tokens := strings.Split(msg.Caption, " ")
            commandName := tokens[0]
            if commandName[0] != '/' { // since we can't check if a caption is a command w/ tg API
                break
            }
            commandParams := tokens[1:]
            switch commandName {
            case "/caption": 
                if len(commandParams) < 1 {
                    helpString := `Correct usage: /caption ["Top Text"] ["Bottom Text"] (Brackets not included)`
                    msgConfig = telebot.NewMessage(msg.Chat.ID, helpString)
                    break
                }
                photos := *msg.Photo
                // msg.Photo is a slice of PhotoSizes from the TG API which are offerings of the same
                // photo in various sizes
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
            default:
                msgConfig = telebot.NewMessage(msg.Chat.ID, "Unknown command")
                msgConfig.ReplyToMessageID = msg.MessageID
            }
            bot.Send(msgConfig)
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
    }
    defer file.Close()

    _, err = io.Copy(file, response.Body)
    if err != nil {
        log.Printf("Failed to write image to disk: %e", err)
    }
    
    return filepath, nil
}
