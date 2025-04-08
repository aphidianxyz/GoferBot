package command

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
)

func TestGeneratingEmptyJSON(t *testing.T) {
	got, err := GenerateCommandJSON("../../test/command/emptyjson.json")
	if err != nil {
		t.Errorf("Did not expect an error: %v", err.Error())
	}
	want := CommandJSON{}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Got: %v\nWanted: %v", got, want)
	}
}

func TestErrorJSONFileNotExisting(t *testing.T) {
	nonExistantFileName := "non-existant-file.json"
	errorStr := fmt.Sprintf("open %v: no such file or directory", nonExistantFileName)
	want := errors.New(errorStr)
	_, got := GenerateCommandJSON(nonExistantFileName)
	if got.Error() != want.Error() { 
		t.Errorf("Got: %v\nWanted: %v", got, want)
	}
}

func TestGenerateCommandJSON(t *testing.T) {
    jsonFilePath := "../../test/command/basic-commands.json"
    commandTest := CommandInfo{CommandName: "test", CommandSyntax: "/test", CommandDesc: "a test command for parsing tests"}
    param1 := Param{ParamName: "param1", ParamDesc: "first param desc"}
    param2 := Param{ParamName: "param2", ParamDesc: "second param desc"}
    commandTest2 := CommandInfo{CommandName: "test2", CommandSyntax: "/test2", CommandDesc: "a second test command for parsing tests", Params: []Param{param1, param2}}
    emptyCmdInfo := CommandInfo{}
    var commands []CommandInfo = []CommandInfo{commandTest, commandTest2, emptyCmdInfo}
    want := CommandJSON{Commands: commands}

    got, err := GenerateCommandJSON(jsonFilePath)
    if err != nil {
        t.Errorf("Did not expect an error")
    }
    if !reflect.DeepEqual(got, want) {
        t.Errorf("Got: %+v\nWanted: %+v", got, want)
    }
}

func TestInvalidJSONFile(t *testing.T) {
	invalidJSONFile := "../../test/command/bad-json.json"
	_, got := GenerateCommandJSON(invalidJSONFile)
	if got == nil {
		t.Error("Expected an error, but got none")
	}
}

func TestFormatCommandInfo(t *testing.T) {
	cmdName := "test"
	jsonFilePath := "../../test/command/basic-commands.json"
	cmdJSON, err := GenerateCommandJSON(jsonFilePath)
	if err != nil {
		t.Error("Did not expect an error")
	}
	want := fmt.Sprintf("test - /test\na test command for parsing tests\n\n")

	got := cmdJSON.formatCommandInfo(cmdName)

	if got != want {
		t.Errorf("Got: %v\nWanted: %v", got, want)
	}
}

func TestFormatAllCommandInfo(t *testing.T) {
	jsonFilePath := "../../test/command/basic-commands.json"
	cmdJSON, err := GenerateCommandJSON(jsonFilePath)
	if err != nil {
		t.Error("Did not expect an error")
	}

	var want string
	for _, command := range cmdJSON.Commands {
		want += format(command)
	}
	got := cmdJSON.formatCommandInfo("")

	if got != want {
		t.Errorf("Got: %v\nWanted: %v", got, want)
	}
}
