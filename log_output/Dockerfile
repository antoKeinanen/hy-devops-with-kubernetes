FROM golang:1.26-alpine

WORKDIR /app

COPY . .

RUN go mod tidy

ENTRYPOINT ["go", "run", "main.go"]

