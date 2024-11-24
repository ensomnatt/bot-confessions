FROM golang

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o cf-bot ./cmd 

CMD ["./cf-bot"]
