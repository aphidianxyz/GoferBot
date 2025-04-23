package command

import (
	"database/sql"
	"fmt"

	telebot "github.com/OvyFlash/telegram-bot-api"
)

type Everyone struct {
	msg telebot.Message
	mentions string // pings
    sendConfig telebot.Chattable
}

func MakeEveryone(msg telebot.Message, db *sql.DB) Command {
    mentions, err := generateMentions(msg, db)
    if err != nil {
		return MakeError(msg, "/everyone", "failed to retrieve users in this chat: " + err.Error())
    }
	return &Everyone{msg: msg, mentions: mentions}
}

func (ec *Everyone) GenerateMessage() {
    msgConfig := telebot.NewMessage(ec.msg.Chat.ID, ec.mentions)
    msgConfig.ParseMode = "MarkDown"
    // link /everyone to the reply of the invoked command
    if replyTarget := ec.msg.ReplyToMessage; replyTarget != nil {
        msgConfig.ReplyParameters.MessageID = replyTarget.MessageID
    }
    ec.sendConfig = msgConfig
}

func (ec *Everyone) SendMessage(api *telebot.BotAPI) error {
    if _, err := api.Send(ec.sendConfig); err != nil {
        return err
    }
    return nil
}

func generateMentions(msg telebot.Message, db *sql.DB) (string, error) {
    var mentionsMessage string
    queryStmt := "select * from chats where chatID=?"
    rows, err := db.Query(queryStmt, msg.Chat.ID)
    if err != nil {
        return "", err 
    }
    defer rows.Close()
    // parse results
    for rows.Next() {
        var id int64
        var chatID int64
        var userID int64
        var username sql.NullString 
        var firstName string
        if err = rows.Scan(&id, &chatID, &userID, &username, &firstName); err != nil {
            return "", err
        }
        // users can omit having a username, but a first name is required, which is used as fallback
        var name = firstName
        if username.String != "" {
            name = username.String
        }
        mention := fmt.Sprintf("[%v](tg://user?id=%v) ", name, userID)
        mentionsMessage += mention
    }
    return mentionsMessage, nil
}
