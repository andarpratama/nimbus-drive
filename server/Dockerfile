FROM golang:1.24.1-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go install github.com/cosmtrek/air@latest

CMD ["air"]
