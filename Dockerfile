FROM golang:alpine

WORKDIR /usr/local/goferbot

RUN apk add --no-cache build-base pkgconfig imagemagick-dev

COPY go.mod go.sum ./
RUN go mod download 

COPY . .

RUN go build -o gofer .

CMD ["./gofer"]
