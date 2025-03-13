package main

import g "github.com/aphidian.xyz/bettergrambot/pkg/gofer"

func main() {
    gofer := g.Gofer{}
    gofer.Initialize()
    gofer.Update(60)
}

