package command

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
)

type CommandJSON struct {
	Commands []CommandInfo `json:"commands"`
}

type CommandInfo struct {
	CommandName string `json:"commandName"`
	CommandDesc string `json:"commandDesc"`
	Params []struct {
		ParamName string `json:"paramName"`
		ParamDesc string `json:"paramDesc"`
	} `json:"params"`
}

func GenerateCommandJSON(commandDescFilePath string) (CommandJSON, error) {
	var commandsDecoded CommandJSON
	commandDescFile, err := os.ReadFile(commandDescFilePath)
	if err != nil {
		return commandsDecoded, err
	}
	decoder := json.NewDecoder(bytes.NewReader(commandDescFile))
	if err := decoder.Decode(&commandsDecoded); err != nil {
		return commandsDecoded, err
	}
	return commandsDecoded, nil
}

func (c CommandJSON) findCommand(commandName string) (found []CommandInfo) {
	for _, command := range c.Commands {
		if command.CommandName == commandName {
			found = append(found, command) 
		}
	}
	return found
}

func formatCommandHelper(command CommandInfo) string {
	name := command.CommandName
	desc := command.CommandDesc
	params := command.Params
	var paramString string
	for _, param := range params {
		paramString += fmt.Sprintf("+ %v: %v\n", param.ParamName, param.ParamDesc)
	}
	return fmt.Sprintf("%v - %v\n%v\n", name, desc, paramString)
}

func (c CommandJSON) formatCommandInfo(commandName string) (formatted string) {
	var commands []CommandInfo
	commands = c.findCommand(commandName)
	for _, command := range commands {
		formatted += formatCommandHelper(command)
	}
	return formatted
}

func (c CommandJSON) formatAllCommandInfo() (formatted string) {
	for _, command := range c.Commands {
		formatted += formatCommandHelper(command)
	}
	return formatted
}
