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
Get your bot API token from (https://telegram.me/BotFather)[BotFather], register a name and set privacy mode to *DISABLED*
Then run it! Make sure you set your TOKEN environment variable to your bot api key provided by telegram's BotFather
```
docker run -e TOKEN=$TOKEN -v chats-db:$GOFERDIR/sql aphidianxyz/goferbot
```
Make sure to give the bot Admin permissions, as some of the commands do not work without it

#### Manual
// TODO
