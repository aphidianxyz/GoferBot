package main

import (
	"os"

	g "github.com/aphidianxyz/GoferBot/pkg/gofer"
)

func main() {
    gofer := g.Gofer{}
    gofer.Initialize("./sql", os.Getenv("TOKEN"))
    gofer.Update(60)
}

