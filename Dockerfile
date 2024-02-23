FROM golang:1.21

RUN mkdir /app
WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build ./cmd/leetgode/main.go

EXPOSE 8080

CMD ["./main"]
