FROM golang:1.22

WORKDIR /app

COPY go.mod go.sum ./
COPY . .

RUN go mod download

EXPOSE 8080

RUN go build ./cmd/main.go

CMD ["go", "run", "./cmd/main.go"]