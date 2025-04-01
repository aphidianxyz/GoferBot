package command

import (
	"reflect"
	"testing"

	telebot "github.com/OvyFlash/telegram-bot-api"
)

func TestGeneratingAboutCommand(t *testing.T) {
	var chatID int64 = 0
	aboutCmd := AboutCommand{chatID: chatID}
	aboutCmd.GenerateMessage()

	got := aboutCmd.sendConfig
	want := telebot.NewMessage(chatID, aboutMsg)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("expected %+v\n, got %+v", got, want)
	}
}
