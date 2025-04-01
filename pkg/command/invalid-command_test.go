package command

import (
	"reflect"
	"testing"

	telebot "github.com/OvyFlash/telegram-bot-api"
)

func TestInvalidCommand(t *testing.T) {
	var chatID int64 = 0
	invalidCmdName := "blah"
	invalidCmdStr := invalidCmdName + invalidCmdSuffix
	want := telebot.NewMessage(chatID, invalidCmdStr) 
	invalidCmd := InvalidCommand{chatID: chatID, request: invalidCmdName}

	invalidCmd.GenerateMessage()
	got := invalidCmd.sendConfig

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Got:%v\nWanted:%v", got, want)
	}
}
