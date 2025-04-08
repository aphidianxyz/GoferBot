package command

import (
	"database/sql"
	"fmt"
	"io/fs"
	"os"
	"reflect"
	"testing"

	telebot "github.com/OvyFlash/telegram-bot-api"
	_ "github.com/mattn/go-sqlite3"
)

func TestEveryoneCommand(t *testing.T) {
	userOneUsername := "userNameHaver123"
	userOneFirstname := "Joseph"
	userTwoFirstName := "Bob"
	testDatabasePath := "../../test/command/everyone-test.db"
	db := createEveryoneTestDBFile(t, testDatabasePath)
	defer os.Remove(testDatabasePath)
	createEveryoneTestDB(t, db)
	var chatID int64 = 0
	var userOneID int64 = 0
	var userTwoID int64 = 1
	insertUsernameUser(t, db, userOneUsername, userOneFirstname, chatID, userOneID)
	insertFirstNameUser(t, db, userTwoFirstName, chatID, userTwoID)
	origin := telebot.Chat{ID: chatID}
	requestMsg := telebot.Message{Chat: origin}
	everyoneCmd := EveryoneCommand{msg: requestMsg, db: db}
	userOnePing := fmt.Sprintf("[%v](tg://user?id=%v) ", userOneUsername, userOneID)
	userTwoPing := fmt.Sprintf("[%v](tg://user?id=%v) ", userTwoFirstName, userTwoID)
	wantConfig := telebot.NewMessage(chatID, userOnePing + userTwoPing)
	wantConfig.ParseMode = "MarkDown"
	want := wantConfig
	
	everyoneCmd.GenerateMessage()
	got := everyoneCmd.sendConfig

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Got: %v\nWanted: %v", got, want)
	}
}

func createEveryoneTestDBFile(t testing.TB, filePath string) *sql.DB {
	t.Helper()
    var perms fs.FileMode = 0644 
	// touch file
	emptyByte := []byte("") 
	err := os.WriteFile(filePath, emptyByte, perms)
	if err != nil {
		t.Errorf("Did not expect an error while creating the everyone test DB file")
	}
	db, err := sql.Open("sqlite3", filePath)
	if err != nil {
		t.Errorf("Did not expect an error while creating the everyone test DB file")
	}
	return db
}

func insertUsernameUser(t testing.TB, db *sql.DB, username, firstName string, chatID, userID int64) {
	stmt := "insert into chats(chatID, userID, username, firstName) values(?, ?, ?, ?)"
	args := []interface{}{chatID, userID, username, firstName}
    if _, err := db.Exec(stmt, args...); err != nil {
		t.Errorf("Did not expect an error when inserting user into everyone test db")
    }
}

func insertFirstNameUser(t testing.TB, db *sql.DB, firstName string, chatID, userID int64) {
	stmt := "insert into chats(chatID, userID, username, firstName) values(?, ?, NULL, ?)"
	args := []interface{}{chatID, userID, firstName}
    if _, err := db.Exec(stmt, args...); err != nil {
		t.Errorf("Did not expect an error when inserting user into everyone test db")
    }
}

func createEveryoneTestDB(t testing.TB, db *sql.DB) {
	t.Helper()
    tableStmt := `
    create table chats(id integer primary key,
    chatID integer not null, 
    userID integer not null, 
    username text, 
    firstName text not null);
    `
    if _, err := db.Exec(tableStmt); err != nil {
		t.Errorf("Did not expect an error when creating everyone test DB")
    }
}
