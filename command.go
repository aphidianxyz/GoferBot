package main

import (
	"errors"
	"strings"

	telebot "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Command interface {
    GenerateMessage() error
    SendMessage(api *telebot.BotAPI) error
}

func ParseMsgCommand(msg *telebot.Message) Command {
    msgStr := msg.Text
    tokens := strings.Split(msgStr, " ")
    commandName := tokens[0]
    commandParams := tokens[1:]
    switch commandName {
    case "/hello": // TODO: utilize enums
        return &HelloCommand{chatID: msg.Chat.ID}
    case "/help":
        helpRequest := commandParams[0]
        // TODO: maybe make /help DM the requesting user instead 
        return &HelpCommand{chatID: msg.Chat.ID, request: helpRequest}
    default:
        return &InvalidCommand{chatID: msg.Chat.ID, request: commandName}
    }
} 

type InvalidCommand struct {
    chatID int64
    request string
    sendConfig telebot.MessageConfig
}

// TODO: maybe this would be better described as "UnknownCommand"
func (ic *InvalidCommand) GenerateMessage() error {
    invalidRequest := ic.request + " is not a valid command!"
    ic.sendConfig = telebot.NewMessage(ic.chatID, invalidRequest)
    return nil
}

func (ic *InvalidCommand) SendMessage(api *telebot.BotAPI) error {
    if _, err := api.Send(ic.sendConfig); err != nil {
        return errors.New("Failed to send an InvalidCommand")
    }
    return nil
}

type HelpCommand struct {
    chatID int64
    request string
    sendConfig telebot.MessageConfig
}

func (hc *HelpCommand) GenerateMessage() error {
    var config telebot.MessageConfig
    switch hc.request {
    case "":
    // TODO: general description
    case "hello":
        config = telebot.NewMessage(hc.chatID, "/hello -- Gofer says hello to you!")
    case "help":
        helpCommandDesc := "/help [command?] -- Describes command functionality and syntax, specific command can be specified"
        config = telebot.NewMessage(hc.chatID, helpCommandDesc)
    default:
        invalidCommandMsg := hc.request + " is not a known command"
        config = telebot.NewMessage(hc.chatID, invalidCommandMsg)
    }
    hc.sendConfig = config
    return nil
}

func (hc *HelpCommand) SendMessage(api *telebot.BotAPI) error {
    if _, err := api.Send(hc.sendConfig); err != nil {
        return errors.New("Failed to send a HelpCommand")
    }
    return nil
}

type HelloCommand struct {
    chatID int64
    greetee string
    sendConfig telebot.MessageConfig
}

func (hc *HelloCommand) SetGreetee(greetee string) {
    hc.greetee = greetee
}

func (hc *HelloCommand) GenerateMessage() error {
    helloString := "Hello, " + hc.greetee
    config := telebot.NewMessage(hc.chatID, helloString)
    hc.sendConfig = config 
    return nil
} 

func (hc *HelloCommand) SendMessage(api *telebot.BotAPI) error {
    if _, err := api.Send(hc.sendConfig); err != nil {
        return errors.New("Failed to send a HelloCommand")
    }
    return nil
}

