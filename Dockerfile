FROM golang:alpine

WORKDIR /usr/local/goferbot

RUN apk add --no-cache build-base pkgconfig imagemagick-dev imagemagick-jpeg imagemagick-raw imagemagick-tiff imagemagick-tiff imagemagick-webp imagemagick-pdf imagemagick-heic

COPY go.mod go.sum ./
RUN go mod download 

COPY . .
RUN touch /usr/local/goferbot/sql/chats.db

RUN go build -o gofer .

CMD ["./gofer"]
