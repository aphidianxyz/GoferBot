package main

import g "github.com/aphidianxyz/GoferBot/pkg/gofer"

func main() {
    gofer := g.Gofer{}
    gofer.Initialize()
    gofer.Update(60)
}

