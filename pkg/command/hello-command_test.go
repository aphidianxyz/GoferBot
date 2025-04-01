package command

import (
	"reflect"
	"testing"

	telebot "github.com/OvyFlash/telegram-bot-api"
)

func TestHelloCommand(t *testing.T) {
	var chatID int64 = 0 
	firstName := "Bob"
	lastName := "Dylan"
	userName := "bobdylan.minecraft2003"
	helloStr := "Hello, " + firstName + " " + lastName + "!\nAKA: " + userName
	helloCmd := HelloCommand{chatID: chatID, firstName: firstName, lastName: lastName, userName: userName}
	want := telebot.NewMessage(chatID, helloStr)

	helloCmd.GenerateMessage()
	got := helloCmd.sendConfig
	
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Expected: %v\nGot: %v", got, want)
	}
}
