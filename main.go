package main

import ( 
	"os"

	g "github.com/aphidianxyz/GoferBot/pkg/gofer"
)

func main() {
	gofer := g.Gofer{DatabasePath: "./sql/chats.db", APIToken: os.Getenv("TOKEN"), CommandJSONFilePath: "./assets/cmd-descriptions/cmd-desc-en.json"}
    gofer.Initialize()
    gofer.Update(60)
}

