package gofer

import (
	"database/sql"
	"errors"
	"io/fs"
	"log"
	"os"
	"strings"

	telebot "github.com/OvyFlash/telegram-bot-api"
	cmd "github.com/aphidianxyz/GoferBot/pkg/command"
	_ "github.com/mattn/go-sqlite3"
)

type Gofer struct {
    DatabasePath string
    ApiToken string

    api *telebot.BotAPI
    db *sql.DB
}

func (g *Gofer) Initialize() {
    g.initAPI(g.ApiToken)
    g.initDB(g.DatabasePath)
}

func (g *Gofer) Update(timeout int) {
    defer g.db.Close()
    updateConfig := telebot.NewUpdate(0)
    updateConfig.Timeout = timeout 

    updates := g.api.GetUpdatesChan(updateConfig)

    for update := range updates {
        // check db health
        if errPing := g.db.Ping(); errPing != nil {
            log.Println("Database closed, attempting to reopen...")
            if err := g.initDB(g.DatabasePath); err != nil {
                log.Panicln("Failed to reopen database. " + err.Error())
            }
        }
        msg := update.Message
        edit := update.EditedMessage // edit is nil when msg isn't and vice-versa
        if msg == nil {
            if edit != nil {
                g.handleEdits(&update)
            }
            continue
        }
        // records user in a chat for usage of fns like /everyone or /mention [ping group]
        if err := g.recordUser(msg); err != nil {
            log.Println(err)
        } 
        if msg.IsCommand() {
            g.handleCommands(&update)
        } else if msg.Photo != nil { // msg w/ photos have captions, manual parsing required
            g.handlePhotoCommands(&update)
        } else { // TODO: handle messages/command requests with a video or gif attached
            // TODO: handle registered responses
        }
    }
}

func (g *Gofer) recordUser(msg *telebot.Message) error {
    chatID := msg.Chat.ID
    userID := msg.From.ID
    username := msg.From.UserName
    firstName := msg.From.FirstName
    var stmt string 
    var args []interface{}
    if userExists(g.db, chatID, userID) {
        if username == "" {
            stmt = "update chats set firstname=?, username=NULL where chatID=?;"
            args = []interface{}{firstName, chatID}
        } else {
            stmt = "update chats set firstname=?, username=? where chatID=?;"
            args = []interface{}{firstName, username, chatID}
        }
    } else {
        if username == "" {
            stmt = "insert into chats(chatID, userID, username, firstName) values(?, ?, NULL, ?)"
            args = []interface{}{chatID, userID, firstName}
        } else {
            stmt = "insert into chats(chatID, userID, username, firstname) values(?, ?, ?, ?)"
            args = []interface{}{chatID, userID, username, firstName}
        }
    }
    if _, err := g.db.Exec(stmt, args...); err != nil {
        return err
    }
    return nil
}

func userExists(db *sql.DB, chatID, userID int64) bool {
    queryStmt := "select exists(select chatID, userID from chats where chatID=? and userID=?) as row_exists;"
    row, err := db.Query(queryStmt, chatID, userID)
    if err != nil {
        return false
    }
    defer row.Close()
    for row.Next() {
        var bool int
        err = row.Scan(&bool)
        if err != nil {
            return false
        }
        return bool == 1
    }
    return false
}

func (g *Gofer) initDB(databasePath string) error {
    // create database directory
    var perms fs.FileMode = 0644
    var err error
    if _, err := os.Stat(databasePath); errors.Is(err, fs.ErrNotExist) {
        emptyByte := []byte("")
        err := os.WriteFile(databasePath, emptyByte, perms)
        if err != nil {
            return err
        }
    }
    g.db, err = sql.Open("sqlite3", databasePath)
    if pingErr := g.db.Ping(); pingErr != nil && err != nil {
        return pingErr
    }
    if err := createChatTables(g.db); err != nil {
        return err
    }
    log.Println("Database initialized.")
    return nil
}

func (g *Gofer) initAPI(token string) {
    var err error
    g.api, err = telebot.NewBotAPI(token)
    if err != nil {
        log.Panic("Failed to initialize bot: " + err.Error())
    }
    log.Println("Bot initialized! Account: " + g.api.Self.UserName)
}

func (g *Gofer) handleCommands(update *telebot.Update) {
    msg := update.Message
    command := cmd.ParseMsgCommand(g.api, msg)
    // TODO: this impl currently doesn't support multi-step commands
    command.GenerateMessage()
    if err := command.SendMessage(g.api); err != nil {
        sendError(msg.Chat.ID, err.Error(), g.api)
        return
    }
}

func (g *Gofer) handlePhotoCommands(update *telebot.Update) {
    msg := update.Message
    if !isCaptionCommand(msg.Caption) {
        return
    }
    command := cmd.ParseImgCommand(g.api, msg)
    command.GenerateMessage()
    if err := command.SendMessage(g.api); err != nil {
        sendError(msg.Chat.ID, err.Error(), g.api)
        return
    }
}

func (g *Gofer) handleEdits(update *telebot.Update) {
    if update.EditedMessage == nil {
        return
    }
    // TODO: add operations if we want to handle certain edit events
}

func isCaptionCommand(caption string) bool {
    tokens := strings.Split(caption, " ")
    if len(tokens) == 0 {
        return false
    }
    commandName := tokens[0]
    return commandName[0] == '/'
}


func sendError(chatID int64, errStr string, api *telebot.BotAPI) {
    errSuffix := "Error: "
    errorMessage := telebot.NewMessage(chatID, errSuffix + errStr)
    api.Send(errorMessage)
}

func createChatTables(db *sql.DB) error {
    tableStmt := `
    create table chats(id integer primary key,
    chatID integer not null, 
    userID integer not null, 
    username text, 
    firstName text not null);
    `
    if _, err := db.Exec(tableStmt); err != nil {
        return err
    }
    return nil
}
