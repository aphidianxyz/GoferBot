package main

import ( 
	"os"

	g "github.com/aphidianxyz/GoferBot/pkg/gofer"
)

func main() {
	gofer := g.CreateBot("./sql/chats.db", os.Getenv("TOKEN"), "./assets/cmd-descriptions/cmd-desc-en.json")
    gofer.Initialize()
    gofer.Update(60)
}

