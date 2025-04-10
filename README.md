# GoferBot
A telegram chat utility bot

### Installation
#### Docker
First pull the docker image:
```
docker pull aphidianxyz/goferbot
```
Before you run the image, create a container so that you can populate a chats database for commands like /everyone
```
docker volume create chats-db
```
Before you can run your own instance of GoferBot, you'll need a Bot API token from Telegram

Get your bot API token by registering a bot instance through the [BotFather](https://telegram.me/BotFather) then set its privacy mode to *DISABLED*

Set your TOKEN environment variable to the API key you've received from the BotFather, then run the bot:
```
docker run -e TOKEN=$TOKEN -v chats-db:$GOFERDIR/sql aphidianxyz/goferbot
```

When you add the bot into your group chat, make sure it is an admin, otherwise several features will not work properly

#### Manual
Install Go 1.21.X, then ImageMagick 7+ on your system.

Compile a build by running `go build` while in the project root

Get your bot API token by registering a bot instance through the [BotFather](https://telegram.me/BotFather) then set its privacy mode to *DISABLED*

Set your TOKEN environment variable to the API key you've received from the BotFather, then run the bot:
```
./GoferBot
```

### Features
1. Captioning images (URL, replies, attachments)
2. Pinging /everyone (only works if people have messaged at least once after the bot was added to a chat)
3. Pinning posts
