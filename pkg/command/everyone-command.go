package command

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	telebot "github.com/OvyFlash/telegram-bot-api"
)

type EveryoneCommand struct {
    chatID int64
    db *sql.DB
    sendConfig telebot.Chattable
}

func (ec *EveryoneCommand) GenerateMessage() {
    mentions, err := ec.generateMentions()
    if err != nil {
        ec.sendConfig = telebot.NewMessage(ec.chatID, "Failed to retrieve users from this chat")
    }

    msgConfig := telebot.NewMessage(ec.chatID, mentions)
    msgConfig.ParseMode = "MarkDown"
    ec.sendConfig = msgConfig
}

func (ec *EveryoneCommand) SendMessage(api *telebot.BotAPI) error {
    if _, err := api.Send(ec.sendConfig); err != nil {
        return errors.New("Failed to send an EveryoneCommand")
    }
    return nil
}

func (ec *EveryoneCommand) generateMentions() (string, error) {
    var mentionsMessage string
    queryStmt := "select * from chats where chatID=?"
    rows, err := ec.db.Query(queryStmt, ec.chatID)
    if err != nil {
        return "", err 
    }
    defer rows.Close()
    // parse results
    for rows.Next() {
        var id int64
        var chatID int64
        var userID int64
        var username string 
        var firstName string
        if err = rows.Scan(&id, &chatID, &userID, &username, &firstName); err != nil {
            return "", err
        }
        // users can omit having a username, but a first name is required, which is used as fallback
        var name = firstName
        if username != "" {
            name = username
        }
        mention := fmt.Sprintf("[%v](tg://user?id=%v) ", name, userID)
        mentionsMessage += mention
    }
    return mentionsMessage, nil
}
