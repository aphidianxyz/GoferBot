package main

import (
	"os"

	g "github.com/aphidianxyz/GoferBot/pkg/gofer"
)

func main() {
    gofer := g.Gofer{DatabasePath: "./sql/chats.sql", ApiToken: os.Getenv("TOKEN")}
    gofer.Initialize()
    gofer.Update(60)
}

