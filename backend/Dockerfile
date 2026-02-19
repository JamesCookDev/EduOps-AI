FROM golang:1.24.0-alpine

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

RUN go build -o eduops ./cmd/main.go

EXPOSE 8080

CMD ["./eduops"]
