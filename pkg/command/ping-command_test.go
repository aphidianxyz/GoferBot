package command

import (
	"testing"
	"reflect"
    telebot "github.com/OvyFlash/telegram-bot-api"
)

func TestPingCommand(t *testing.T) {
	var chatID int64 = 0
	pingCmd := PingCommand{chatID: chatID}
	want := telebot.NewMessage(chatID, "pong")

	pingCmd.GenerateMessage()
	got := pingCmd.sendConfig
	
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Expected: %v\nGot: %v", got, want)
	}
}
