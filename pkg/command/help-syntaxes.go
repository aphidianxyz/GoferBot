package command

const (
    helloSyntax string = "/hello - Gofer greets you!"
    helpSyntax string = "/help [command?] - Describes command functionality and syntax\ncommand (optional) - a specific command to describe"
    pingSyntax string = "/ping - Bot responds with 'pong'"
    captionSyntax string = "/caption [url] [\"top\"] [\"bot\"] - Creates an impact font caption meme\nurl - the url of the image to be captioned\n\"top\" - the top caption, encapsulated by quotes\n\"bot\" - the bottom caption, encapsulated by quotes"
    captionImgSyntax string = "/caption [\"top\"] [\"bot\"] (with an image attached) - Creates an impact font caption meme\n\"top\" - the top caption, encapsulated by quotes\n\"bot\" - the bottom caption, encapsulated by quotes"
    everyoneSyntax string = "/everyone [message?]\nmessage (optional) - a message that accompanies a ping to everyone"
)

var allHelpSyntaxes = []string{helloSyntax, helpSyntax, pingSyntax, captionSyntax, captionImgSyntax, everyoneSyntax}

