FROM golang:1.23.4-alpine
WORKDIR /app

RUN mkdir logs

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o main .

CMD ["./main"]
EXPOSE 8080
