FROM go 1.24.0

WORKDIR /usr/local/goferbot

COPY . . 
RUN go mod download 

RUN go build -o gofer .

CMD ["./gofer"]
