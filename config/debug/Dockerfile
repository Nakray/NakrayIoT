FROM golang:1.23.4-alpine
WORKDIR /app

RUN apk add --no-cache git && \
    go install github.com/go-delve/delve/cmd/dlv@latest

COPY . .
RUN go build -o main .

CMD ["dlv", "debug", "--headless", "--listen=:40000", "--api-version=2", "--accept-multiclient", "./main.go"]
