package command

import (
	"fmt"
	"reflect"
	"slices"
	"testing"

	telebot "github.com/OvyFlash/telegram-bot-api"
)

func TestSpecificHelpCommandWOSlashes(t *testing.T) {
	specificCommandsTestHelper(t, false)
}

func TestSpecificHelpCommandsWSlashes(t *testing.T) {
	specificCommandsTestHelper(t, true)
}

func TestAllCommands(t *testing.T) {
	var chatID int64 = 0
	cmdJSON := generateCmdJSON(t)
	msgConfig := telebot.NewMessage(chatID, cmdJSON.formatAllCommandInfo())
	msgConfig.ParseMode = "MarkDown"
	want := msgConfig
	
	helpCmd := HelpCommand{chatID: chatID, request: "", commandJSON: cmdJSON}
	helpCmd.GenerateMessage()
	got := helpCmd.sendConfig

	compareChattables(t, got, want)
}

func TestInvalidHelp(t *testing.T) {
	var chatID int64 = 0
	cmdJSON := generateCmdJSON(t)
	invalidCommandName := "/invalidcommand"
	invalidCommandMsg := fmt.Sprintf("%v%v%v", invalidCmdMsgPrefix, invalidCommandName, invalidCmdMsgSuffix)
	msgConfig := telebot.NewMessage(chatID, invalidCommandMsg)
	want := msgConfig
	
	helpCmd := HelpCommand{chatID: chatID, request: invalidCommandName, commandJSON: cmdJSON}
	helpCmd.GenerateMessage()
	got := helpCmd.sendConfig

	compareChattables(t, got, want)
}

func generateCmdJSON(t testing.TB) CommandJSON {
	t.Helper()
	cmdJSON, err := GenerateCommandJSON("../../assets/cmd-descriptions/cmd-desc-en.json")
	if err != nil {
		t.Errorf("Invalid command description json file")
	}
	return cmdJSON
}

func compareChattables(t testing.TB, got, want telebot.Chattable) {
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %+v\nwanted %+v", got, want)
	}
}

func specificCommandsTestHelper(t testing.TB, slashes bool) {
	t.Helper()
	var chatID int64 = 0
	cmdJSON := generateCmdJSON(t)
	var commands []string
	commands = getAllCommandNames(t, cmdJSON, slashes)

	want := createMessages(t, chatID, commands, cmdJSON)
	got := generateMessages(t, chatID, commands, cmdJSON)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %+v\nwanted %+v", got, want)
	}
}

func generateMessages(t testing.TB, chatID int64, commands []string, cmdJSON CommandJSON) (messages []telebot.Chattable) {
	t.Helper()
	for _, command := range commands {
		helpCmd := HelpCommand{chatID: chatID, request: command, commandJSON: cmdJSON}
		helpCmd.GenerateMessage()
		messages = append(messages, helpCmd.sendConfig)
	}
	return messages
}

func createMessages(t testing.TB, chatID int64, commands []string, cmdJSON CommandJSON) (messages []telebot.Chattable) {
	t.Helper()
	for _, command := range commands {
		formattedCommand := cmdJSON.formatCommandInfo(command)
		wantMessage := telebot.NewMessage(chatID, formattedCommand)
		wantMessage.ParseMode = "MarkDown"
		messages = append(messages, wantMessage)
	}
	return messages
}

func getAllCommandNames(t testing.TB, cmdJSON CommandJSON, slashes bool) (commandNames []string) {
	t.Helper()
	for _, command := range cmdJSON.Commands {
		commandName := command.CommandName
		if !slices.Contains(commandNames, commandName) {
			if slashes {
				commandName = "/" + commandName
			}
			commandNames = append(commandNames, command.CommandName)
		}
	}
	return commandNames
}
